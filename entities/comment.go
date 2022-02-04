package entities

type Comment struct {
	UserName  string `json:"username" form:"username"`
	Avatar    string `json:"avatar" form:"avatar"`
	Content   string `json:"content" form:"content"`
	CreatedAt string `json:"created_at" form:"created_at"`
}
