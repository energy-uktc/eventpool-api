package models

type SimpleUserModel struct {
	Email    string `json:"email"`
	UserName string `json:"userName"`
}

type UserModel struct {
	Id string `json:"id"`
	SimpleUserModel
}

type UpdateUserModel struct {
	UserName string `json:"userName" binding:"required"`
}
