package request

// Request User Create
type UserCreateRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Request Login User
type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
