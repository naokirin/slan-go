package spreadsheets

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"

	"github.com/naokirin/slan-go/app/domain/lunch"
)

var _ lunch.Repository = (*Spreadsheets)(nil)

// Spreadsheets for reading google spreadsheet
type Spreadsheets struct{}

// Retrieve a token, saves the token, then returns the generated client.
func (s *Spreadsheets) getClient(config *oauth2.Config, tokenPath string) (*http.Client, error) {
	tok, err := s.tokenFromFile(tokenPath)
	if err != nil {
		log.Fatalf("Token file:%s is not found", tokenPath)
		return nil, err
	}
	return config.Client(context.Background(), tok), nil
}

// Retrieves a token from a local file.
func (s *Spreadsheets) tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// GetRows returns spreadsheet rows
func (s *Spreadsheets) GetRows(sheetID string, readRange string, secretPath string, tokenPath string) [][]string {
	result := make([][]string, 0)

	b, err := ioutil.ReadFile(secretPath)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
		return result
	}

	// If modifying these scopes, delete your previously saved client_secret.json.
	config, err := google.ConfigFromJSON(b, sheets.SpreadsheetsReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
		return result
	}

	client, err := s.getClient(config, tokenPath)
	if err != nil {
		log.Fatalf("Unable to retrieve Google client: %v", err)
		return result
	}
	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
		return result
	}

	// Prints the names and majors of students in a sample spreadsheet:
	// https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit
	resp, err := srv.Spreadsheets.Values.Get(sheetID, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
		return result
	}

	if len(resp.Values) > 0 {
		for _, row := range resp.Values {
			// Print columns A and E, which correspond to indices 0 and 4.
			r := make([]string, 0)
			for _, x := range row {
				r = append(r, x.(string))
			}
			result = append(result, r)
		}
	}

	return result
}
