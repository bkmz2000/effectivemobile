package db

import (
	"encoding/json"
	"fmt"
	"time"
)

// Song represents a song in the database
// @Description Song model
// @Property id int true "ID" example(1)
// @Property title string true "Title" example("Song Title")
// @Property artist string true "Artist" example("Artist Name")
// @Property release_date string true "Release Date" example("2023-01-01")
// @Property text string true "Text" example("Song lyrics here...")
type Song struct {
	Id          uint      `json:"id"`
	Name        string    `json:"name"`
	Group       string    `json:"group"`
	ReleaseDate time.Time `json:"release_date"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
}

// UnmarshalJSON for custom date parsing
func (s *Song) UnmarshalJSON(b []byte) error {
	type Alias Song // Create an alias to avoid recursion
	aux := &struct {
		ReleaseDate string `json:"release_date"` // Use string for custom parsing
		*Alias
	}{
		Alias: (*Alias)(s),
	}

	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}

	// Parse the release date
	parsedDate, err := time.Parse("2006-01-02", aux.ReleaseDate)
	if err != nil {
		return fmt.Errorf("invalid date format: %s", aux.ReleaseDate)
	}
	s.ReleaseDate = parsedDate

	return nil
}
