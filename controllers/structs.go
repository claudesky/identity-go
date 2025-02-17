package controllers

type RegisterRequest struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

type LoginRequest = RegisterRequest

type RefreshRequest struct {
	RefreshToken *string `json:"refresh_token"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Message struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}
