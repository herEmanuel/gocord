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

type ChangePasswordJSON struct {
	OldPassword          string `json:"oldPassword"`
	NewPassword          string `json:"newPassword"`
	ConfirmationPassword string `json:"confirmationPassword"`
}

type EnterServerJSON struct {
	InviteCode string `json:"inviteCode"`
}
