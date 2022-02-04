package entities

type User struct {
	Id           int `json:"id" form:"id"`
	Name         string
	Email        string
	Password     string
	Organization string
	PhoneNumber  string
	Avatar       string
}
