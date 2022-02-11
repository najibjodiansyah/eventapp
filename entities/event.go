package entities

type Event struct {
	Id          int    `json:"id" form:"id"`
	Name        string `json:"name" form:"name"`
	UserName    string `json:"username" form:"username"`
	Host        string `json:"host" form:"host"`
	Description string `json:"description" form:"description"`
	Datetime    string `json:"datetime" form:"datetime"`
	Location    string `json:"location" form:"location"`
	Category    string `json:"category" form:"category"`
	Photo       string `json:"photo" form:"photo"`
	HostId      int    `json:"hostid" form:"hostid"`
}
