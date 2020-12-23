package models

type RequestResult struct {
	Status  int    `default:0`
	Error   string `default:"something wrong"`
	Message string `default:""`
	Data    string
}
