package entities

type Event struct {
	Id          int    `json:"id" form:"id"`
	Name        string `json:"name" form:"name"`
	Category    string `json:"category" form:"category"`
	Host        string `json:"host" form:"host"`
	Description string `json:"description" form:"description"`
	Datetime    string `json:"datetime" form:"datetime"`
	Location    string `json:"location" form:"location"`
	Photo       string `json:"photo" form:"photo"`
	UserName    string `json:"username" form:"username"`
	HostId      int    `json:"hostid" form:"hostid"`
}
