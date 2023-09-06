package convertutil

import (
	"reflect"
	"time"
)

type UserPo struct {
	Id        int64
	Username  string
	Age       int64
	CreatedAt time.Time
}

type UserVo struct {
	Id        int64
	Username  string
	Age       int64
	CreatedAt string
}

func ConvertStruct[K, V any](source K) (target V) {
	valueOf := reflect.ValueOf(source)
	typeOf := reflect.TypeOf(target)
	targetValue := reflect.New(typeOf).Elem()

	for i := 0; i < valueOf.NumField(); i++ {
		field := valueOf.Type().Field(i)
		if structField, b := typeOf.FieldByName(field.Name); b && structField.Type == field.Type {
			targetValue.FieldByName(structField.Name).Set(valueOf.Field(i))
		}
	}
	return targetValue.Interface().(V)
}
