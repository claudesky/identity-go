package controllers

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/claudesky/identity-go/interfaces/repositories"
	"github.com/claudesky/identity-go/models"
	"github.com/claudesky/identity-go/requests"
	"github.com/claudesky/identity-go/responses"
	"github.com/claudesky/identity-go/services"
	"github.com/claudesky/identity-go/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	th  *services.TokenHandler
	ur  repositories.UserRepository
	tfr repositories.TokenFamilyRepository
	cr  repositories.ClientRepository
	acr repositories.AuthorizationCodeRepository
}

func NewAuthController(
	th *services.TokenHandler,
	ur repositories.UserRepository,
	tfr repositories.TokenFamilyRepository,
	cr repositories.ClientRepository,
	acr repositories.AuthorizationCodeRepository,
) *AuthController {
	return &AuthController{th, ur, tfr, cr, acr}
}

// -- Controller Methods

func (c *AuthController) Login(
	w http.ResponseWriter,
	r *http.Request,
	rq *requests.Login,
) {
	// Get user
	user, err := c.ur.GetUserByEmail(r.Context(), *rq.Email)
	if err != nil {
		// Horrible error handling, should get this handled outside?
		slog.Info("could not find user by email",
			slog.String("error", fmt.Sprintf("%v", err)),
			slog.Group("params",
				slog.String("email", *rq.Email),
			),
		)
		unauthorized(w)
		return
	}

	if user.Password == nil {
		// Better handling later
		slog.Warn("user has no password")
		unauthorized(w)
		return
	}

	// Check Password
	err = bcrypt.CompareHashAndPassword(
		[]byte(*user.Password),
		[]byte(*rq.Password),
	)
	if err != nil {
		unauthorized(w)
		return
	}

	// Token Family ID
	jtf := utils.PseudoUUID()
	now := time.Now()

	// Generate Tokens
	refreshString, tokenString, expRT, err := c.th.GenerateTokens(
		&now,
		user.Id,
		jtf,
		jtf,
		// wonder how we would deal with custom claims in future
		jwt.MapClaims{
			"name":  user.Name,
			"email": user.Email,
		},
	)
	if err != nil {
		internalServerError(w)
		return
	}

	c.tfr.InsertToken(r.Context(), &models.TokenFamily{
		Id:           jtf,
		Sub:          user.Id,
		LastIssued:   jtf,
		CreatedAt:    now,
		LastIssuedAt: now,
		ExpiresAt:    expRT,
	})

	respondWithData(w, "Success", &responses.TokenResponse{
		AccessToken:  tokenString,
		RefreshToken: refreshString,
	}, http.StatusOK)
}

func (c *AuthController) Refresh(
	w http.ResponseWriter,
	r *http.Request,
	rq *requests.Refresh,
) {
	// Verify Token
	ok, claims := c.th.VerifyRefreshToken(*rq.RefreshToken)
	if !ok {
		unauthorized(w)
		return
	}

	// Get the token family
	tf, err := c.tfr.GetTokenById(r.Context(), claims.JTF)
	if err != nil {
		slog.Warn(
			"could not find token family",
			"token",
			rq.RefreshToken,
			"sub",
			claims.SUB,
			"error",
			err,
		)
		unauthorized(w)
		return
	}

	// Make sure sub is correct
	if claims.SUB != tf.Sub {
		slog.Warn(
			"invalid sub for refresh token",
			"token",
			rq.RefreshToken,
			"sub",
			claims.SUB,
			"family_sub",
			tf.Sub,
		)
		unauthorized(w)
		return
	}

	// Check last issued
	if claims.JTI != tf.LastIssued {
		slog.Warn("invalid last issued for refresh token", "token", rq.RefreshToken)
		// TODO: Add revoke here
		unauthorized(w)
		return
	}

	user, err := c.ur.GetUserById(r.Context(), claims.SUB)
	if err != nil {
		slog.Warn(
			"could not user id",
			"user id",
			claims.SUB,
			"error",
			err,
		)
		unauthorized(w)
		return
	}

	// Possible future checks
	// ...

	// New Refresh Token ID
	jtiRT := utils.PseudoUUID()

	now := time.Now()

	// Generate new tokens
	refreshString, tokenString, expRT, err := c.th.
		GenerateTokens(&now, tf.Sub, tf.Id, jtiRT,
			jwt.MapClaims{
				"name":  user.Name,
				"email": user.Email,
			})
	if err != nil {
		internalServerError(w)
		return
	}

	// Update Token Family
	tf.LastIssued = jtiRT
	tf.LastIssuedAt = now
	tf.ExpiresAt = expRT

	c.tfr.UpdateToken(r.Context(), tf)

	// Send new tokens
	respondWithData(w, "Success", &responses.TokenResponse{
		AccessToken:  tokenString,
		RefreshToken: refreshString,
	}, http.StatusOK)
}

func (c *AuthController) Validate(w http.ResponseWriter, r *http.Request) {
	tokenString := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

	token, err := c.th.VerifyToken(tokenString)
	if err != nil {
		slog.Info("token verification failed", "error", err)
		unauthorized(w)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		respondWithData(w, "Success", claims, http.StatusOK)
	} else {
		fmt.Println(err)
		internalServerError(w)
		return
	}
}

func (c *AuthController) Authorize(
	w http.ResponseWriter,
	r *http.Request,
	rq *requests.Authorize,
) {
	// Verify client exists and redirect_uri is valid
	client, err := c.cr.GetClientById(r.Context(), *rq.ClientId)
	if err != nil {
		slog.Info("client not found",
			slog.String("client_id", *rq.ClientId),
			slog.String("error", fmt.Sprintf("%v", err)),
		)
		unauthorized(w)
		return
	}

	// Validate redirect URI
	validRedirect := false
	for _, uri := range client.RedirectUris {
		if uri == *rq.RedirectUri {
			validRedirect = true
			break
		}
	}
	if !validRedirect {
		slog.Warn("invalid redirect_uri",
			slog.String("client_id", *rq.ClientId),
			slog.String("redirect_uri", *rq.RedirectUri),
		)
		unauthorized(w)
		return
	}

	// Get authenticated user from context (assumes middleware sets this)
	// For now, we'll need to authenticate the user first
	// This endpoint should be called AFTER user login/authentication
	// You might want to get user from session or require them to login first

	// For this implementation, assume user is authenticated via session/cookie
	// and stored in context by middleware
	userId, ok := r.Context().Value("user_id").(string)
	if !ok || userId == "" {
		slog.Warn("user not authenticated for authorization")
		unauthorized(w)
		return
	}

	// Create authorization code
	scope := ""
	if rq.Scope != nil {
		scope = *rq.Scope
	}

	codeChallenge := ""
	codeChallengeMethod := ""
	if rq.CodeChallenge != nil {
		codeChallenge = *rq.CodeChallenge
	}
	if rq.CodeChallengeMethod != nil {
		codeChallengeMethod = *rq.CodeChallengeMethod
	}

	authCode := models.NewAuthorizationCode(
		client.Id,
		userId,
		*rq.RedirectUri,
		scope,
		codeChallenge,
		codeChallengeMethod,
	)

	err = c.acr.InsertAuthorizationCode(r.Context(), authCode)
	if err != nil {
		slog.Error("failed to insert authorization code",
			slog.String("error", fmt.Sprintf("%v", err)),
		)
		internalServerError(w)
		return
	}

	// Build redirect URL with code and state
	redirectUrl, _ := url.Parse(*rq.RedirectUri)
	query := redirectUrl.Query()
	query.Set("code", authCode.Code)
	if rq.State != nil {
		query.Set("state", *rq.State)
	}
	redirectUrl.RawQuery = query.Encode()

	// Redirect to client's redirect_uri with authorization code
	http.Redirect(w, r, redirectUrl.String(), http.StatusFound)
}

func (c *AuthController) TokenExchange(
	w http.ResponseWriter,
	r *http.Request,
	rq *requests.TokenExchange,
) {
	// Verify client credentials
	client, err := c.cr.GetClientById(r.Context(), *rq.ClientId)
	if err != nil {
		slog.Info("client not found",
			slog.String("client_id", *rq.ClientId),
		)
		unauthorized(w)
		return
	}

	// Verify client secret
	err = bcrypt.CompareHashAndPassword(
		[]byte(client.ClientSecret),
		[]byte(*rq.ClientSecret),
	)
	if err != nil {
		slog.Warn("invalid client secret",
			slog.String("client_id", *rq.ClientId),
		)
		unauthorized(w)
		return
	}

	// Get authorization code
	authCode, err := c.acr.GetAuthorizationCodeByCode(r.Context(), *rq.Code)
	if err != nil {
		slog.Warn("authorization code not found",
			slog.String("code", *rq.Code),
		)
		unauthorized(w)
		return
	}

	// Verify code hasn't been used
	if authCode.Used {
		slog.Warn("authorization code already used",
			slog.String("code", *rq.Code),
		)
		unauthorized(w)
		return
	}

	// Verify code hasn't expired
	if time.Now().After(authCode.ExpiresAt) {
		slog.Warn("authorization code expired",
			slog.String("code", *rq.Code),
		)
		unauthorized(w)
		return
	}

	// Verify client_id matches
	if authCode.ClientId != client.Id {
		slog.Warn("client_id mismatch",
			slog.String("code_client_id", authCode.ClientId),
			slog.String("request_client_id", client.Id),
		)
		unauthorized(w)
		return
	}

	// Verify redirect_uri matches
	if authCode.RedirectUri != *rq.RedirectUri {
		slog.Warn("redirect_uri mismatch",
			slog.String("code_redirect_uri", authCode.RedirectUri),
			slog.String("request_redirect_uri", *rq.RedirectUri),
		)
		unauthorized(w)
		return
	}

	// Verify PKCE if code_challenge was provided
	if authCode.CodeChallenge != "" {
		if rq.CodeVerifier == nil {
			slog.Warn("code_verifier required but not provided")
			unauthorized(w)
			return
		}

		// Verify code_verifier against code_challenge
		valid := verifyPKCE(authCode.CodeChallenge, authCode.CodeChallengeMethod, *rq.CodeVerifier)
		if !valid {
			slog.Warn("invalid code_verifier")
			unauthorized(w)
			return
		}
	}

	// Get user
	user, err := c.ur.GetUserById(r.Context(), authCode.UserId)
	if err != nil {
		slog.Warn("user not found",
			slog.String("user_id", authCode.UserId),
		)
		unauthorized(w)
		return
	}

	// Mark code as used
	err = c.acr.MarkCodeAsUsed(r.Context(), authCode.Code)
	if err != nil {
		slog.Error("failed to mark code as used",
			slog.String("code", authCode.Code),
		)
		internalServerError(w)
		return
	}

	// Generate tokens
	jtf := utils.PseudoUUID()
	now := time.Now()

	refreshString, tokenString, expRT, err := c.th.GenerateTokens(
		&now,
		user.Id,
		jtf,
		jtf,
		jwt.MapClaims{
			"name":  user.Name,
			"email": user.Email,
			"scope": authCode.Scope,
		},
	)
	if err != nil {
		internalServerError(w)
		return
	}

	// Store token family
	c.tfr.InsertToken(r.Context(), &models.TokenFamily{
		Id:           jtf,
		Sub:          user.Id,
		LastIssued:   jtf,
		CreatedAt:    now,
		LastIssuedAt: now,
		ExpiresAt:    expRT,
	})

	respondWithData(w, "Success", &responses.TokenResponse{
		AccessToken:  tokenString,
		RefreshToken: refreshString,
	}, http.StatusOK)
}

// Verify PKCE code_verifier against code_challenge
func verifyPKCE(codeChallenge, method, codeVerifier string) bool {
	switch method {
	case "S256":
		hash := sha256.Sum256([]byte(codeVerifier))
		computed := base64.RawURLEncoding.EncodeToString(hash[:])
		return computed == codeChallenge
	case "plain":
		return codeVerifier == codeChallenge
	}
	return false
}
