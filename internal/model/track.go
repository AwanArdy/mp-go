package model

import "time"

type Track struct {
	ID        string
	Title     string
	Artist    string
	Album     string
	Path      string
	Duration  time.Duration
	CoverData []byte
}
