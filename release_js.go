package main

import "syscall/js"

type JSReleasable struct {
	Releasable
	callbacks []js.Callback
}

func (r *JSReleasable) deferRelease(callback js.Callback) {
	r.callbacks = append(r.callbacks, callback)
}

func (r *JSReleasable) Release() {
	for _, r := range r.callbacks {
		r.Release()
	}
	r.callbacks = nil
}
