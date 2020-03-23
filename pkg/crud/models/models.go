package models

type User  struct {
		Id              int
		Name            string
		Login           string
		Password        string
		Role            string
		Removed         bool
}

