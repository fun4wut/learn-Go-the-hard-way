package main

import (
	"reflect"
)

func MakeMap(fpt interface{}) {
	fnV := reflect.ValueOf(fpt).Elem()
	fnI := reflect.MakeFunc(fnV.Type(), implMap)
	fnV.Set(fnI)
}

//TODO:completes implMap function.
// 第一个参数是fn，第二个是arr | map，返回新的arr | map
func implMap(args []reflect.Value) []reflect.Value {
	obj := args[1]
	fn := args[0]
	var newObj reflect.Value
	switch obj.Kind() {
	case reflect.Slice:
		newObj = reflect.MakeSlice(obj.Type(), obj.Len(), obj.Cap())
		for i := 0; i < obj.Len(); i++ {
			newObj.Index(i).Set(
				fn.Call(
					[]reflect.Value{obj.Index(i)},
				)[0],
			)
		}
	case reflect.Map:
		newObj = reflect.MakeMap(obj.Type())
		iter := obj.MapRange()
		for iter.Next() {
			newObj.SetMapIndex(
				iter.Key(),
				fn.Call([]reflect.Value{iter.Value()})[0],
			)
		}
	}
	return []reflect.Value{newObj}
}

func main() {

	println("It is said that Go has no generics.\nHowever we have many other ways to implement a generics like library if less smoothly,one is reflect.MakeFunc.\nUnderscore is a very useful js library,and now let's implement part of it-map,it will help you to understand how reflect works.\nPlease finish the 'implMap' function and pass the test.")
}
