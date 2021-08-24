package models

type UserModel struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	UserName string `json:"userName"`
}

type UpdateUserModel struct {
	UserName string `json:"userName" binding:"required"`
}
