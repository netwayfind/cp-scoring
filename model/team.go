package model

type TeamSummary struct {
	ID   int64
	Name string
}

type Team struct {
	ID      int64
	Name    string
	POC     string
	Email   string
	Enabled bool
}
