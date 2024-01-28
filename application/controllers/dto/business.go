package dto

type BusinessDTO struct {
	Name  string `json:"name"`
	Email *string `json:"email"`
}

type UpdateBusinessDTO struct {
	Name string `json:"name" validate:"required"`
	ID 	 string `json:"id" validate:"required"`
}