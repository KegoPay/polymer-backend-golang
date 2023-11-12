package webhook

type CustomerVerificationDTO struct {
	Event string 					  `json:"event"`
	Data  CustomerVerificationDTOData `json:"data"`
}

type CustomerVerificationDTOData struct {
	CustomerID 	 string  `json:"customer_id"`
	CustomerCode string  `json:"customer_code"`
	Email 		 string  `json:"email"`
	Reason 		 string  `json:"reason"`
}