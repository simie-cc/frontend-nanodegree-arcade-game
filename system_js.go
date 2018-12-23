package main

import (
	"fmt"
	"syscall/js"
)

type JSSystemService struct {
	JSReleasable
}

var systemService = &JSSystemService{}

func init() {
	systemService = &JSSystemService{}
}

func getWindow() js.Value {
	return js.Global().Get("window")
}

func (s *JSSystemService) RegisterGetVerboseLevel(fn func() int) {
	callbackFn := js.NewCallback(func(args []js.Value) {
		callId := args[0].String()
		val := fn()
		fmt.Println("Getreturn vale", val)
		getWindow().Set(fmt.Sprintf("rv_%v", callId), val)
	})
	s.deferRelease(callbackFn)

	getWindow().Set("GetVerboseLevel", callbackFn)
}

func (s *JSSystemService) RegisterSetVerboseLevel(fn func(newLevel int)) {
	callbackFn := js.NewCallback(func(args []js.Value) {
		intValue := args[0].Int()
		fn(intValue)
	})
	s.deferRelease(callbackFn)

	getWindow().Set("SetVerboseLevel", callbackFn)
}
