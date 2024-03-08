package background

import (
	"os"

	"kego.com/infrastructure/background/gocraft"
)

var Scheduler SchedulerType

func StartScheduler() {
	Scheduler =  (&gocraft.GoCraftScheduler{
		CacheAddress: os.Getenv("REDIS_ADDR"),
		CachePassword: os.Getenv("REDIS_PASSWORD"),
	})
	Scheduler.StartScheduler()
}
