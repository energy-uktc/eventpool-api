package models

type UserModel struct {
	Email    string `json:"email"`
	UserName string `json:"userName"`
}

type UpdateUserModel struct {
	UserName string `json:"userName" binding:"required"`
}
