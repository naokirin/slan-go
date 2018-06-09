package main

import (
	"flag"
	"io/ioutil"
	"log"

	g "github.com/naokirin/slan-go/cmd/google/internal"
	"golang.org/x/oauth2/google"
	sheets "google.golang.org/api/sheets/v4"
)

var (
	clientSecretFlag = flag.String("s", "secrets/google_client_secrets.json", "client secret file")
	tokenFileFlag    = flag.String("t", "secrets/google_token.json", "token file")
)

func main() {
	flag.Parse()

	b, err := ioutil.ReadFile(*clientSecretFlag)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved client_secret.json.
	config, err := google.ConfigFromJSON(b, sheets.SpreadsheetsReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	g.GetClient(config, *tokenFileFlag)
}
