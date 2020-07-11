package main

import (
	"reflect"
)

//Reverse reverses a slice.
var Reverse func(slice interface{}) = func(slice interface{}) {
	var t = reflect.ValueOf(slice).Elem()  // 使用Elem获取指针所指的元素
	swap := reflect.Swapper(t.Interface()) // 使用Interface获取对应的值
	for i, j := 0, t.Len()-1; i < j; {
		swap(i, j)
		i++
		j--
	}
}

func main() {
	println("Please edit main.go,and complete the 'Reverse' function to pass the test.\nYou should use reflect package to reflect the slice type and make it applly to any type.\nTo run test,please run 'go test'\nIf you pass the test,please run 'git checkout l2' ")
}
