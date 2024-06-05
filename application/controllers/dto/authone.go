package dto

type AuthOneSendEmail struct {
	Email    string         `json:"email"`
	Template string         `json:"template"`
	Subject  string         `json:"subject"`
	Opts     map[string]any `json:"opts"`
}

type AuthOneCreateUserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
