package main

import (
	"fmt"
	"syscall/js"
)

type JSRenderer struct {
	JSReleasable

	document js.Value
	window   js.Value
	ctx      js.Value

	imageCount int
	imageCache map[string]js.Value
	imageDone  chan bool

	renderFn  func()
	rendering bool
}

var renderer = &JSRenderer{}

func init() {
	renderer = &JSRenderer{
		imageCount: 0,
		imageCache: make(map[string]js.Value),
		imageDone:  make(chan bool, 30),
	}
}

func (r *JSRenderer) Init(width, height int) {
	r.document = js.Global().Get("document")
	r.window = js.Global().Get("window")

	canvasEl := r.document.Call("getElementById", "mycanvas")
	canvasEl.Set("width", width)
	canvasEl.Set("height", height)
	r.ctx = canvasEl.Call("getContext", "2d")
}

func (r *JSRenderer) UpdateFPS(val string) {
	fpsDiv := r.document.Call("getElementById", "fps")
	fpsDiv.Set("innerHTML", val)
}

func (r *JSRenderer) PrepareImage(url string) {
	imageType := r.window.Get("Image")
	image := imageType.New()

	var imageLoad js.Callback
	imageLoad = js.NewCallback(func(args []js.Value) {

		fmt.Println("image loaded!!", url)
		image.Set("onload", nil)
		imageLoad.Release()
		r.imageDone <- true
	})

	image.Set("onload", imageLoad)
	image.Set("src", url)
	r.imageCache[url] = image
	r.imageCount++
}

func (r *JSRenderer) WaitImage() {
	received := 0
	for {
		<-r.imageDone
		received++
		if received >= r.imageCount {
			break
		}
	}
}

func (r *JSRenderer) ClearRect() {
	r.ctx.Call("clearRect", 0, 0, width, height)
}

func (r *JSRenderer) DrawRect(rect *Rect) {
	r.ctx.Call("strokeRect", rect.x1, rect.y1, (rect.x2 - rect.x1), (rect.y2 - rect.y1))
}

func (r *JSRenderer) DrawImage(url string, x, y int) {
	r.ctx.Call("drawImage", r.imageCache[url], x, y)
	r.DrawRect(&Rect{x1: x, y1: y, x2: x + 3, y2: y + 3})
}

func (r *JSRenderer) ListenKeyboardEvent(fn func(eventType, key string)) {
	var keyboardEventHandler = js.NewCallback(func(args []js.Value) {
		event := args[0]
		eventType := event.Get("type").String()
		key := event.Get("key").String()
		fn(eventType, key)
	})
	r.deferRelease(keyboardEventHandler)

	r.document.Call("addEventListener", "keydown", keyboardEventHandler)
}

func (r *JSRenderer) ListenClickEvent(fn func(eventType string, clientX, clientY int)) {
	canvas := r.document.Call("getElementById", "mycanvas")

	var eventHandler = js.NewCallback(func(args []js.Value) {
		event := args[0]
		eventType := event.Get("type").String()
		x := event.Get("clientX").Int()
		y := event.Get("clientY").Int()
		rect := canvas.Call("getBoundingClientRect")
		rx, ry := rect.Get("left").Int(), rect.Get("top").Int()

		fn(eventType, x-rx, y-ry)
	})
	r.deferRelease(eventHandler)

	canvas.Call("addEventListener", "click", eventHandler)
}

func (r *JSRenderer) ListenMouseMoveEvent(fn func(eventType string, clientX, clientY int)) {
	canvas := r.document.Call("getElementById", "mycanvas")

	var eventHandler = js.NewCallback(func(args []js.Value) {
		event := args[0]
		eventType := event.Get("type").String()
		// Refer: https://www.html5canvastutorials.com/advanced/html5-canvas-mouse-coordinates/
		x := event.Get("clientX").Int()
		y := event.Get("clientY").Int()
		rect := canvas.Call("getBoundingClientRect")
		rx, ry := rect.Get("left").Int(), rect.Get("top").Int()

		fn(eventType, x-rx, y-ry)
	})
	r.deferRelease(eventHandler)

	canvas.Call("addEventListener", "mousemove", eventHandler)
}

func (r *JSRenderer) RegisterRenderFunction(fn func()) {
	r.renderFn = fn
}

func (r *JSRenderer) StartRender() {
	if r.rendering {
		return
	}

	r.rendering = true

	var renderFrame js.Callback
	renderFrame = js.NewCallback(func(args []js.Value) {
		if r.rendering {
			defer js.Global().Call("requestAnimationFrame", renderFrame)
			r.renderFn()
		} else {
			renderFrame.Release()
		}
	})

	js.Global().Call("requestAnimationFrame", renderFrame)
}

func (r *JSRenderer) StopRender() {
	r.rendering = false
}

func (r *JSRenderer) SetFont(fontDesc string) {
	r.ctx.Set("font", fontDesc)
}

func (r *JSRenderer) DrawText(text string, x, y int) {
	r.ctx.Call("fillText", text, x, y)
}

func (r *JSRenderer) SetTextAlign(align string) {
	r.ctx.Set("textAlign", align)
}

func (r *JSRenderer) DrawRoundedRect(rect *Rect) {
	const red = 10
	r.ctx.Call("moveTo", rect.x1+red, rect.y1)
	r.ctx.Call("lineTo", rect.x2-red, rect.y1)
	r.ctx.Call("arcTo", rect.x2, rect.y1, rect.x2, rect.y1+red, red)

	r.ctx.Call("lineTo", rect.x2, rect.y2-red)
	r.ctx.Call("arcTo", rect.x2, rect.y2, rect.x2-red, rect.y2, red)

	r.ctx.Call("lineTo", rect.x1+red, rect.y2)
	r.ctx.Call("arcTo", rect.x1, rect.y2, rect.x1, rect.y2-red, red)

	r.ctx.Call("lineTo", rect.x1, rect.y1+red)
	r.ctx.Call("arcTo", rect.x1, rect.y1, rect.x1+red, rect.y1, red)

	r.ctx.Call("stroke")
}
