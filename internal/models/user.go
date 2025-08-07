package models

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleId   int    `json:"roleId"`
}

type UserInfo struct {
	ID     int    `json:"id"`
	Email  string `json:"email"`
	RoleId int    `json:"roleId"`
}
