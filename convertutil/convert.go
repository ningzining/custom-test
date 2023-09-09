package convertutil

import (
	"reflect"
)

func ConvertStruct[T any](source interface{}) (target T) {
	sourceValue := reflect.Indirect(reflect.ValueOf(source))

	targetType := reflect.TypeOf(target)
	targetValue := reflect.New(targetType).Elem()

	for i := 0; i < sourceValue.NumField(); i++ {
		sourceField := sourceValue.Type().Field(i)
		if targetField, exist := targetType.FieldByName(sourceField.Name); exist {
			if sourceField.Type == targetField.Type && sourceField.Type.Kind() != reflect.Struct {
				targetValue.FieldByName(targetField.Name).Set(sourceValue.Field(i))
			}
		}
	}
	return targetValue.Interface().(T)
}
