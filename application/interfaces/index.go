package interfaces

type ApplicationContext[T interface{}] struct{
	Body *T
	Keys map[string]any
	Header map[string][]string
	Ctx any
}

func (ac *ApplicationContext[T]) GetContextData(key string) (value any, exists bool) {
	value, exists = ac.Keys[key]
	return
}

func (ac *ApplicationContext[T]) SetContextData(key string, data any) {
	if ac.Keys == nil {
		ac.Keys = map[string]any{}
	}
	ac.Keys[key] = data
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

func (ac *ApplicationContext[T]) GetHeader(key string) (value any) {
	header := ac.Header
	if header == nil {
		return nil
	}
	if header[key] == nil {
		return nil
	}
	return header[key][0]
}
