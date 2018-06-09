package calendar

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	domain "github.com/naokirin/slan-go/app/domain/calendar"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

var _ domain.Calendar = (*Calendar)(nil)

// Calendar is implemented calendar repository
type Calendar struct{}

// Retrieve a token, saves the token, then returns the generated client.
func (c *Calendar) getClient(config *oauth2.Config, tokenPath string) (*http.Client, error) {
	tok, err := c.tokenFromFile(tokenPath)
	if err != nil {
		log.Fatalf("Token file:%s is not found", tokenPath)
		return nil, err
	}
	return config.Client(context.Background(), tok), nil
}

// Retrieves a token from a local file.
func (c *Calendar) tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// GetCalendarItems returns CalendarItems
func (c *Calendar) GetCalendarItems(min time.Time, max time.Time, secretPath string, tokenPath string) []domain.Item {

	result := []domain.Item{}

	b, err := ioutil.ReadFile(secretPath)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
		return result
	}

	// If modifying these scopes, delete your previously saved client_secret.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
		return result
	}

	client, err := c.getClient(config, tokenPath)
	if err != nil {
		log.Fatalf("Unable to retrieve Google client: %v", err)
		return result
	}
	srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
		return result
	}

	minf := min.Format(time.RFC3339)
	maxf := max.Format(time.RFC3339)
	events, err := srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(minf).TimeMax(maxf).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
		return result
	}
	if len(events.Items) == 0 {
		return result
	}
	for _, item := range events.Items {
		var (
			start time.Time
			end   time.Time
		)
		if item.Start.DateTime != "" {
			start, _ = time.Parse("2006-01-02T15:04:05-07:00", item.Start.DateTime)
		} else {
			start, _ = time.Parse("2006-01-02", item.Start.Date)
		}
		if item.End.DateTime != "" {
			end, _ = time.Parse("2006-01-02T15:04:05-07:00", item.End.DateTime)
		} else {
			end, _ = time.Parse("2006-01-02", item.End.Date)
		}

		data := domain.Item{
			Summary:  item.Summary,
			Location: item.Location,
			Start:    start,
			End:      end,
		}
		result = append(result, data)
	}
	return result
}
