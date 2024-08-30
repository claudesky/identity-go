package main

import (
	"crypto"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var addr string = ":9102"
var pkey crypto.PrivateKey
var pubkey crypto.PublicKey
var dbConn = "postgresql://postgres@localhost/identity_go"
var dbPass string

func init() {
	// PORT
	v, x := os.LookupEnv("IDG_PORT")
	if x && v != "" {
		addr = ":" + v
	}

	// PRIVATE KEY
	v, x = os.LookupEnv("IDG_PKEY_64")
	if x && v != "" {
		pkeyPem, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			panic(err)
		}

		pkey, err = jwt.ParseEdPrivateKeyFromPEM([]byte(pkeyPem))
		if err != nil {
			panic(err)
		}
	} else {
		v, x := os.LookupEnv("IDG_PKEY_PATH")
		if x && v != "" {
			var err error
			pkey, err = fetchPkey(v)
			if err != nil {
				panic(err)
			}
		} else {
			var err error
			pkey, err = fetchPkey("./keys/private.pem") // Default IDG_PKEY_PATH
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

		pubkey, err = jwt.ParseEdPrivateKeyFromPEM([]byte(pubkeyPem))
		if err != nil {
			panic(err)
		}
	} else {
		v, x := os.LookupEnv("IDG_PUBKEY_PATH")
		if x && v != "" {
			var err error
			pubkey, err = fetchPubkey(v)
			if err != nil {
				panic(err)
			}
		} else {
			var err error
			pubkey, err = fetchPubkey("./keys/public.pem") // Default IDG_PUBKEY_PATH
			if err != nil {
				panic(err)
			}
		}
	}

	// DB CONNECTION
	v, x = os.LookupEnv("IDG_DB_CONN")
	if x && v != "" {
		dbConn = v
	}

	// DB PASSWORD
	v, x = os.LookupEnv("IDG_DB_PASS")
	if x && v != "" {
		dbPass = v
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
