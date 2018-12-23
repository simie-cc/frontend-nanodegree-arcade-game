package main

type Renderer interface {
	Init(width, height int)
	Release()
	UpdateFPS(val string)

	PrepareImage(url string)
	WaitImage()

	ClearRect()
	DrawRect(rect *Rect)
	DrawRoundedRect(rect *Rect)
	DrawImage(url string, x, y int)
	DrawText(text string, x, y int)
	SetFont(fontDesc string)
	SetTextAlign(align string)

	ListenKeyboardEvent(fn func(eventType, key string))
	ListenClickEvent(fn func(eventType string, clientX, clientY int))
	RegisterRenderFunction(fn func())

	StartRender()
	StopRender()
}

type EmptyRenderer struct {
	Renderer
}

var renderer Renderer = &EmptyRenderer{}
