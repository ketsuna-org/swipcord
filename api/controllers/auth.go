package controllers

type LoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterBody struct {
	LoginBody
	Username string `json:"username"`
}
