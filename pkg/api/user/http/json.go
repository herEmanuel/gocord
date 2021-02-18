package http

type RegisterJSON struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type LoginJSON struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
