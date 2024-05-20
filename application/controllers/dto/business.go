package dto

type BusinessDTO struct {
	Name  string `json:"name"`
	Email string `json:"-"`
}

type UpdateBusinessDTO struct {
	Name string `json:"name" validate:"required"`
	ID   string `json:"id" validate:"required"`
}

type SearchCACByName struct {
	Name string `json:"name" validate:"required"`
}

type SetCACInfo struct {
	RCNumber string `json:"rcNumber" validate:"required,min=7,max=7"`
}
