package entities

type Event struct {
	ID          int    `json:"id" form:"id"`
	Name        string `json:"name" form:"name"`
	Category    string `json:"category"`
	Host        string `json:"host"`
	Description string `json:"description"`
	Datetime    string `json:"datetime"`
	Location    string `json:"location"`
	Photo       string `json:"photo"`
	CreatedAt   string `json:"createdAt"`
	UpdateAt    string `json:"updateAt"`
	UserId      int    `json:"userid" form:"userid"`
	Username 	string 	`json:"username"`
}
