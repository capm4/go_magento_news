package model

type Config struct {
	Id    int64  `json:"id" db:"id"`
	Path  string `json:"path" db:"path"`
	Value string `json:"value" db:"value"`
}
