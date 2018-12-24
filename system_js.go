package main

import (
	"syscall/js"
)

type JSSystemService struct {
	JSReleasable
}

var systemService = &JSSystemService{}

func getWindow() js.Value {
	return js.Global().Get("window")
}

func (s *JSSystemService) RegisterToggleDebug(fn func()) {
	callbackFn := js.NewCallback(func(args []js.Value) {
		fn()
	})
	s.deferRelease(callbackFn)

	getWindow().Set("ToggleDebug", callbackFn)
}
