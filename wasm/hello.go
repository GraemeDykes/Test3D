package main

import (
	"fmt"
	"syscall/js"
	"time"
)

func add(i []js.Value) {
	// js.Global().Set("output", js.ValueOf(i[0].Int()+i[1].Int()))
	// println(js.ValueOf(i[0].Int() + i[1].Int()).String())
	fmt.Println("add")
}

func subtract(i []js.Value) {
	// js.Global().Set("output", js.ValueOf(i[0].Int()-i[1].Int()))
	// println(js.ValueOf(i[0].Int() - i[1].Int()).String())
	fmt.Println("subtract")

	// func (v Value) Invoke(args ...interface{}) Value
	// https://godoc.org/syscall/js

	// https://medium.zenika.com/go-1-11-webassembly-for-the-gophers-ae4bb8b1ee03
	alert := js.Global().Get("alert")
	alert.Invoke("subtract called!")
}

func messageReceived(i []js.Value) {
	alert := js.Global().Get("alert")
	alert.Invoke("Message received.")
}

func executeUserCommand(i []js.Value) {
	alert := js.Global().Get("alert")
	alert.Invoke("executeUserCommand - " + i[0].String())
}

func registerCallbacks() {
	js.Global().Set("add", js.NewCallback(add))
	js.Global().Set("subtract", js.NewCallback(subtract))
	js.Global().Set("messageReceived", js.NewCallback(messageReceived))
	js.Global().Set("executeUserCommand", js.NewCallback(executeUserCommand))
}

func main() {
	fmt.Printf("hello, world\n")

	// Call wasm from JS:
	// https://tutorialedge.net/golang/go-webassembly-tutorial/

	registerCallbacks()

	js.Global().Get("console").Call("log", "Hello world Go/wasm!")

	for {

		js.Global().Get("document").Call("getElementById", "app").Set("innerText", time.Now().String())

		time.Sleep(1 * time.Second)
	}

}
