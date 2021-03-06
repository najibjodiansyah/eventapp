package entities

type Comment struct {
	Id        int    `json:"id" form:"id"`
	UserId    int    `json:"userId" form:"userId"`
	UserName  string `json:"username" form:"username"`
	Avatar    string `json:"avatar" form:"avatar"`
	Content   string `json:"content" form:"content"`
	CreatedAt string `json:"createdAt" form:"createdAt"`
}
