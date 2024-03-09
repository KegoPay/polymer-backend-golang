package dto

type ErrorSupportRequestDTO struct {
	Message 		string	`json:"msg" validate:"required,max=300"`
}
