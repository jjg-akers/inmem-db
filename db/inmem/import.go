package inmem

import "time"

// Import ...
type Import struct {
	ID         string
	Status     string
	ImportType string
	ImportTime time.Time
	FileName   string
	User       string
}
