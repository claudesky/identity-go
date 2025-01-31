package main

import (
	"crypto"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

const defaultPkeyPath = "./keys/private.pem"
const defaultPubkeyPath = "./keys/public.pem"

var idg_port string = ":9102"
var idg_pkey crypto.PrivateKey
var idg_pubkey crypto.PublicKey
var idg_db_conn = "postgresql://postgres@localhost/identity_go"
var idg_db_pass string
var idg_mail_user string
var idg_mail_pass string
var idg_mail_host string
var idg_mail_port string
var idg_mail_addr string
var idg_send_init_email bool

func getEnv(key string, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func init() {
	var v string
	var x bool

	// PORT
	idg_port = getEnv("IDG_PORT", idg_port)

	// PRIVATE KEY
	v, x = os.LookupEnv("IDG_PKEY_64")
	if x && v != "" {
		// Parse private key directly from base64 encoded env
		pkeyPem, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			panic(err)
		}

		idg_pkey, err = jwt.ParseEdPrivateKeyFromPEM([]byte(pkeyPem))
		if err != nil {
			panic(err)
		}
	} else {
		// Get file contents of path as private key
		v, x := os.LookupEnv("IDG_PKEY_PATH")
		if x && v != "" {
			var err error
			idg_pkey, err = fetchPkey(v)
			if err != nil {
				panic(err)
			}
		} else {
			// Use the default private key path
			var err error
			idg_pkey, err = fetchPkey(defaultPkeyPath)
			if err != nil {
				panic(err)
			}
		}
	}

	// PUBLIC KEY
	v, x = os.LookupEnv("IDG_PUBKEY_64")
	if x && v != "" {
		// Parse public key directly from base64 encoded env
		pubkeyPem, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			panic(err)
		}

		idg_pubkey, err = jwt.ParseEdPrivateKeyFromPEM([]byte(pubkeyPem))
		if err != nil {
			panic(err)
		}
	} else {
		// Get file contents of path as public key
		v, x := os.LookupEnv("IDG_PUBKEY_PATH")
		if x && v != "" {
			var err error
			idg_pubkey, err = fetchPubkey(v)
			if err != nil {
				panic(err)
			}
		} else {
			// Use the default public key path
			var err error
			idg_pubkey, err = fetchPubkey(defaultPubkeyPath)
			if err != nil {
				panic(err)
			}
		}
	}

	// DB CONNECTION
	idg_db_conn = getEnv("IDG_DB_CONN", idg_db_conn)

	// DB PASSWORD
	idg_db_pass = getEnv("IDG_DB_PASS", "")

	// MAIL
	idg_mail_user = getEnv("IDG_MAIL_USER", "")
	idg_mail_pass = getEnv("IDG_MAIL_PASS", "")
	idg_mail_host = getEnv("IDG_MAIL_HOST", "")
	idg_mail_port = getEnv("IDG_MAIL_PORT", "")
	idg_mail_addr = getEnv("IDG_MAIL_ADDR", "")

	// SEND INIT EMAIL
	v, x = os.LookupEnv("IDG_SEND_INIT_EMAIL")
	if !x && v == "true" {
		idg_send_init_email = true
	}
}

func fetchPkey(path string) (crypto.PrivateKey, error) {
	// Get the Private Key
	pkeyBytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jwt.ParseEdPrivateKeyFromPEM(pkeyBytes)
}

func fetchPubkey(path string) (crypto.PublicKey, error) {
	// Get the Public Key
	pubkeyBytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jwt.ParseEdPublicKeyFromPEM(pubkeyBytes)
}
