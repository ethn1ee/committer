package models

type Diff struct {
	Path       string
	StatusCode string
	Before     string
	After      string
}
