package entity

type User struct {
	Id                string `json:"id,omitempty"`
	Name              string `json:"name" binding:"required"`
	Email             string `json:"email" binding:"required"`
	Age               int    `json:"age" binding:"required"`
	Gender            string `json:"gender" binding:"required"`
	Lastattendacetime string `json:"last attendance time" binding:"required"`
	Password          string `json:"password" binding:"required"`
}
type Token struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
