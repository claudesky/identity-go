package main

import (
	"crypto"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

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

func init() {
	// PORT
	v, x := os.LookupEnv("IDG_PORT")
	if x && v != "" {
		idg_port = ":" + v
	}

	// PRIVATE KEY
	v, x = os.LookupEnv("IDG_PKEY_64")
	if x && v != "" {
		pkeyPem, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			panic(err)
		}

		idg_pkey, err = jwt.ParseEdPrivateKeyFromPEM([]byte(pkeyPem))
		if err != nil {
			panic(err)
		}
	} else {
		v, x := os.LookupEnv("IDG_PKEY_PATH")
		if x && v != "" {
			var err error
			idg_pkey, err = fetchPkey(v)
			if err != nil {
				panic(err)
			}
		} else {
			var err error
			idg_pkey, err = fetchPkey("./keys/private.pem") // Default IDG_PKEY_PATH
			if err != nil {
				panic(err)
			}
		}
	}

	// PUBLIC KEY
	v, x = os.LookupEnv("IDG_PUBKEY_64")
	if x && v != "" {
		pubkeyPem, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			panic(err)
		}

		idg_pubkey, err = jwt.ParseEdPrivateKeyFromPEM([]byte(pubkeyPem))
		if err != nil {
			panic(err)
		}
	} else {
		v, x := os.LookupEnv("IDG_PUBKEY_PATH")
		if x && v != "" {
			var err error
			idg_pubkey, err = fetchPubkey(v)
			if err != nil {
				panic(err)
			}
		} else {
			var err error
			idg_pubkey, err = fetchPubkey("./keys/public.pem") // Default IDG_PUBKEY_PATH
			if err != nil {
				panic(err)
			}
		}
	}

	// DB CONNECTION
	v, x = os.LookupEnv("IDG_DB_CONN")
	if x && v != "" {
		idg_db_conn = v
	}

	// DB PASSWORD
	v, x = os.LookupEnv("IDG_DB_PASS")
	if x && v != "" {
		idg_db_pass = v
	}

	// MAIL

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
