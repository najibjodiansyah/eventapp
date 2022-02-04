package entities

type Participant struct {
	Name   string `json:"name" form:"name"`
	Avatar string `json:"avatar" form:"avatar"`
}
