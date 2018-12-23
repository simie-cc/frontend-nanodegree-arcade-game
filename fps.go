package main

import (
	"fmt"
	"time"
)

var fps = 0
var FPS_LIMIT = 20
var FRAME_INTERVAL = time.Second / time.Duration(FPS_LIMIT)
var lastFrameAt time.Time = time.Now()

func StartFpsCounting() {
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

func waitToNextFrame() {
	duration := time.Since(lastFrameAt)
	if duration < FRAME_INTERVAL {
		time.Sleep(FRAME_INTERVAL - duration)
	}

	lastFrameAt = time.Now()
}
