package prembly

type PremblyFaceMatchResponse struct {
	Status      	bool    		`json:"status"`
	Detail      	string  		`json:"detail"`
	Message		 	string 		`json:"message"`
	Confidence      float32 		`json:"confidence"`
}
