package model

type Signout struct {
	Id  string `json:"id" binding:"required"`
	Jwt string `json:"jwt" binding:"required"`
}
