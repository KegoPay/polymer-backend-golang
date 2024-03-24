package interfaces

import "kego.com/application/utils"

type ApplicationContext[T interface{}] struct{
	Body *T
	Keys map[string]any
	Query map[string]any
	Param map[string]any
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

func (ac *ApplicationContext[T]) GetFloat64ContextData(key string) (value float64) {
	if val, ok := ac.GetContextData(key); ok && val != nil {
		value = val.(float64)
	}
	return
}

func (ac *ApplicationContext[T]) GetBoolContextData(key string) (value bool) {
	if val, ok := ac.GetContextData(key); ok && val != nil {
		value = val.(bool)
	}
	return
}

func (ac *ApplicationContext[T]) GetParameter(key string) any {
	param := ac.Param[key]
	return param
}

func (ac *ApplicationContext[T]) GetStringParameter(key string) string {
	param := ac.Param[key]
	return param.(string)
}

func (ac *ApplicationContext[T]) GetBoolParameter(key string) bool {
	param := ac.Param[key]
	return param.(bool)
}

func (ac *ApplicationContext[T]) GetHeader(key string) (value *string) {
	header := ac.Header
	if header == nil {
		return nil
	}
	if header[key] == nil {
		return nil
	}
	return utils.GetStringPointer(header[key][0])
}
