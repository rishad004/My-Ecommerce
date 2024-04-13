package models

type Login struct {
	Email    string `json:"email"`
	Password string `json:"pass"`
}
type AddCat struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
type Coup struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Code        string `json:"code"`
	Condition   int    `json:"condition"`
	Value       int    `json:"off"`
	Duration    int    `json:"day"`
}