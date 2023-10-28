package datastore

import "context"

var cancelCtx *context.CancelFunc

func ConnectToDatabase(){
	cancelCtx = connectMongo()
}

func CleanUp(){
	(*cancelCtx)()
}