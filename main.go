package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var addr = flag.String("addr", ":9102", "Http Service Address")
var privKey = flag.String("privKey", "./keys/private.pem", "Private Key")
var pubKey = flag.String("pubKey", "./keys/public.pem", "Public Key")

type Message struct {
	Message string `json:"message"`
}

func main() {
	flag.Parse()

	pkey, err := fetchPkey()
	if err != nil {
		log.Fatal(err)
		return
	}

	pubkey, err := fetchPubkey()
	if err != nil {
		log.Fatal(err)
		return
	}

	tokenHandler := newTokenHandler(pkey, pubkey)

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			httpHandler(tokenHandler, w, r)
		}))

	log.Fatal((http.ListenAndServe(*addr, mux)))
}

func httpHandler(
	th *TokenHandler,
	w http.ResponseWriter,
	r *http.Request,
) {
	switch r.URL.Path {
	case "/hello-world":
		helloWorld(w, r)
	case "/healthcheck":
		health(w, r)
	case "/auth/login":
		authLogin(th, w, r)
	case "/auth/validate":
		authValidate(th, w, r)
	default:
		http.Error(w, "Not Found.", http.StatusNotFound)
	}
}

func health(w http.ResponseWriter, _ *http.Request) {
	json.NewEncoder(w).Encode(&Message{Message: "ok"})
}

func helloWorld(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

func authLogin(th *TokenHandler, w http.ResponseWriter, _ *http.Request) {
	// TODO: Login Stuff here

	// Assume OK
	tokenString, err := signToken(th, jwt.MapClaims{
		"jti": "123",
		"jtf": "familyId",
		"jtp": "parendId",
		"sub": "subject",
	})
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(&Message{Message: tokenString})
}

func authValidate(th *TokenHandler, w http.ResponseWriter, r *http.Request) {
	tokenString := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

	token, err := verifyToken(th, tokenString)

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		json.NewEncoder(w).Encode(&Message{Message: fmt.Sprint(claims)})
	} else {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
