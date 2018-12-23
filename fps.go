package main

import (
	"fmt"
	"time"
)

var fps = 0
var FPS_LIMIT = 40
var FRAME_INTERVAL = time.Second / time.Duration(FPS_LIMIT)
var lastFrameAt time.Time = time.Now()

func StartFpsCounting() {
	fmt.Println("FRAME_INTERVAL", FRAME_INTERVAL)

	secondTicker := time.NewTicker(1 * time.Second)
	fpsDiv := document.Call("getElementById", "fps")

	go func() {
		for {
			<-secondTicker.C
			fpsDiv.Set("innerHTML", fmt.Sprintf("%d", fps))
			fps = 0
		}
	}()
}

func AddFps() {
	fps++
}

func shouldRenderNextFrame() bool {
	duration := time.Since(lastFrameAt)
	if duration <= FRAME_INTERVAL {
		return false
	}

	lastFrameAt = time.Now()
	return true
}
