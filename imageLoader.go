package main

import (
	"fmt"
	"syscall/js"
)

var imageCount int = 0
var imageCache = make(map[string]js.Value)
var imageDone = make(chan bool, 30)

func prepareImage(url string) {
	imageType := window.Get("Image")
	image := imageType.New()

	var imageLoad js.Callback
	imageLoad = js.NewCallback(func(args []js.Value) {

		fmt.Println("image loaded!!", url)
		image.Set("onload", nil)
		imageLoad.Release()
		imageDone <- true
	})

	image.Set("onload", imageLoad)
	image.Set("src", url)
	imageCache[url] = image
	imageCount++
}

func waitImage() {
	received := 0
	for {
		<-imageDone
		received++
		if received >= imageCount {
			break
		}
	}
}
