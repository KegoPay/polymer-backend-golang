package background

type SchedulerType interface {
	StartScheduler()
	Emit(channel string, payload map[string]any) error
}