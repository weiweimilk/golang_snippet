package main

// refer to "The Laws of Refelection"
import (
	"reflect"
	"fmt"
)

func Decorator(decoPtr, fn interface{}) {
	var decoratedFunc, targetFunc reflect.Value

	decoratedFunc = reflect.ValueOf(decoPtr).Elem()
	targetFunc = reflect.ValueOf(fn)

	v := reflect.MakeFunc(targetFunc.Type(),
		func(in []reflect.Value) (out []reflect.Value) {
			fmt.Println("before")
			out = targetFunc.Call(in)
			fmt.Println("after")
			return
		})

	decoratedFunc.Set(v)
	return
}

func foo(a, b, c int) int {
	fmt.Printf("%d, %d, %d \n", a, b, c)
	return a + b + c
}

func bar(a, b string) string {
	fmt.Printf("%s, %s \n", a, b)
	return a + b
}

func main() {

	type MyFoo func(int, int, int) int
	var myFoo MyFoo
	Decorator(&myFoo, foo)
	myFoo(1, 2, 3)


	mybar := bar
	Decorator(&mybar, bar)
	mybar("hello", "world!")
}
