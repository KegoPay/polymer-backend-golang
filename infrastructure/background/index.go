package background

import (
	"os"

	"usepolymer.co/infrastructure/background/gocraft"
)

var Scheduler SchedulerType

func StartScheduler() {
	Scheduler = (&gocraft.GoCraftScheduler{
		CacheAddress:  os.Getenv("REDIS_ADDR"),
		CachePassword: os.Getenv("REDIS_PASSWORD"),
	})
	Scheduler.StartScheduler()
}
