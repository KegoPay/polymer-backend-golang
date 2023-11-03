package interfaces

type ApplicationContext[T interface{}] struct{
	Body *T
	Keys map[string]any
	Ctx any
}

func (ac *ApplicationContext[T]) GetContextData(key string) (value any, exists bool) {
	value, exists = ac.Keys[key]
	return
}

func (ac *ApplicationContext[T]) GetStringContextData(key string) (value string) {
	if val, ok := ac.GetContextData(key); ok && val != nil {
		value = val.(string)
	}
	return
}

func (ac *ApplicationContext[T]) GetBoolContextData(key string) (value bool) {
	if val, ok := ac.GetContextData(key); ok && val != nil {
		value = val.(bool)
	}
	return
}

