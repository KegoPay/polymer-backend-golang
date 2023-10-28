package infrastructure

func StartServer(){
	var server serverInterface = &ginServer{}
	server.Start()
}