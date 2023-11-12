package identityverification

// type IdentityPayload struct {
// 	Country string `json:"country"`
// 	Type string `json:"type"`
// 	AccountNumber string `json:"account_number"`
// 	BVN string `json:"bvn"`
// 	BankCode string `json:"bank_code"`
// 	FirstName string `json:"first_name"`
// 	LastName string `json:"last_name"`
// }


type CustomerPayload struct {
	Email         string `json:"email"`
	Phone     	  string `json:"phone"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
}

type AccountPayload struct {
	BVN         	string `json:"bvn"`
	AccountNumber   string `json:"account_number"`
	BankCode    	string `json:"bank_code"`
}


type CustomerVerificationPayload struct {
	Country         string `json:"country"`
	Type     	  	string `json:"type"`
	AccountNumber   string `json:"account_number"`
	BVN     	  	string `json:"bvn"`
	BankCode     	string `json:"bank_code"`
	FirstName     	string `json:"first_name"`
	LastName      	string `json:"last_name"`
}

type CustomerCreationDTO struct {
	Status 	bool						 `json:"status"`
	Message string						 `json:"message"`
	Data	CustomerCreationDTODataField `json:"data"`
}

type CustomerCreationDTODataField struct {
	Email 		 string `json:"email"`
	Integration  int	`json:"integration"`
	CustomerCode string	`json:"customer_code"`
	ID 			 int	`json:"id"`
}