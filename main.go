//Wasming
// compile: GOOS=js GOARCH=wasm go build -o main.wasm ./main.go
package main

import (
	"fmt"
	"math/rand"
	"syscall/js"
	"time"
)

const (
	width  int = 505
	height int = 606

	CELL_WIDTH  = 101
	CELL_HEIGHT = 83

	ROW_COUNT = 6
	COL_COUNT = 5

	IMAGE_Y_SHIFT = 50

	IMAGE_WATER = "images/water-block.png"
	IMAGE_STONE = "images/stone-block.png"
	IMAGE_GRASS = "images/grass-block.png"
	IMAGE_GIRL  = "images/char-cat-girl.png"
	IMAGE_BUG   = "images/enemy-bug.png"
)

var (
	RowImages = [...]string{
		IMAGE_WATER,
		IMAGE_STONE,
		IMAGE_STONE,
		IMAGE_STONE,
		IMAGE_GRASS,
		IMAGE_GRASS,
	}
)

var (
	done = make(chan bool, 0)

	document js.Value
	window   js.Value
	ctx      js.Value
)

var (
	player     *Player
	enemies    []*Enemy
	enemyCount = 5
	rendering  = true
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	fmt.Println("Start")

	document = js.Global().Get("document")
	window = js.Global().Get("window")

	canvasEl := document.Call("getElementById", "mycanvas")
	canvasEl.Set("width", width)
	canvasEl.Set("height", height)
	ctx = canvasEl.Call("getContext", "2d")

	prepareImage(IMAGE_WATER)
	prepareImage(IMAGE_STONE)
	prepareImage(IMAGE_GRASS)
	prepareImage(IMAGE_GIRL)
	prepareImage(IMAGE_BUG)

	waitImage()

	player = NewPlayer()
	enemies = make([]*Enemy, enemyCount)
	for idx, _ := range enemies {
		enemies[idx] = NewEnemy()
	}

	StartFpsCounting()

	fmt.Println("is continue1")

	var keyboardEventHandler = js.NewCallback(handleKeyboardEvent)
	defer keyboardEventHandler.Release()

	document.Call("addEventListener", "keydown", keyboardEventHandler)

	var timeTick = time.Now()

	var renderFrame js.Callback
	renderFrame = js.NewCallback(func(args []js.Value) {
		if rendering {
			defer js.Global().Call("requestAnimationFrame", renderFrame)
		} else {
			return
		}

		if !shouldRenderNextFrame() {
			return
		}

		dt := time.Since(timeTick)
		timeTick = time.Now()

		clearRect()
		drawBaseGround()
		drawEntity(player, IMAGE_GIRL)
		drawRect(player.GetCollisionRect())

		for _, enemy := range enemies {
			enemy.tick(dt)
			drawEntity(enemy, IMAGE_BUG)
			drawRect(enemy.GetCollisionRect())

			if player.isCollision(enemy) {
				fmt.Println("Collision!!")
				rendering = false
				return
			}
		}

		AddFps()
	})
	defer renderFrame.Release()

	js.Global().Call("requestAnimationFrame", renderFrame)

	<-done

	fmt.Println("Main exit")
}

func stopAnimation() {
	rendering = false
}

func clearRect() {
	ctx.Call("clearRect", 0, 0, width, height)
}

func drawRect(rect *Rect) {
	ctx.Call("strokeRect", rect.x1, rect.y1, (rect.x2 - rect.x1), (rect.y2 - rect.y1))
}

func drawBaseGround() {
	for r := 0; r < ROW_COUNT; r++ {
		for c := 0; c < COL_COUNT; c++ {
			x, y := c*CELL_WIDTH, r*CELL_HEIGHT
			ctx.Call("drawImage", imageCache[RowImages[r]], x, y)
			drawRect(&Rect{x1: x, y1: y, x2: x + 3, y2: y + 3})
		}
	}

}

func drawEntity(p Positional, charurl string) {
	x, y := p.GetXY()
	ctx.Call("drawImage", imageCache[charurl],
		x, y)
}

func handleKeyboardEvent(args []js.Value) {
	event := args[0]
	eventType := event.Get("type")
	key := event.Get("key").String()

	switch key {
	case "ArrowUp":
		fmt.Println("Event", eventType, key)
		player.move(DIRECTION_UP)
	case "ArrowDown":
		fmt.Println("Event", eventType, key)
		player.move(DIRECTION_DOWN)
	case "ArrowLeft":
		fmt.Println("Event", eventType, key)
		player.move(DIRECTION_LEFT)
	case "ArrowRight":
		fmt.Println("Event", eventType, key)
		player.move(DIRECTION_RIGHT)
	}
}

// 	// Init Canvas stuff
// 	doc := js.Global().Get("document")
// 	canvasEl := doc.Call("getElementById", "mycanvas")
// 	width = doc.Get("body").Get("clientWidth").Float()
// 	height = doc.Get("body").Get("clientHeight").Float()
// 	canvasEl.Set("width", width)
// 	canvasEl.Set("height", height)
// 	ctx = canvasEl.Call("getContext", "2d")

// 	done := make(chan struct{}, 0)

// 	dt := DotThing{speed: 160, size: 6}

// 	mouseMoveEvt := js.NewCallback(func(args []js.Value) {
// 		e := args[0]
// 		mousePos[0] = e.Get("clientX").Float()
// 		mousePos[1] = e.Get("clientY").Float()
// 	})
// 	defer mouseMoveEvt.Release()

// 	// Event handler for count range
// 	// countChangeEvt := js.NewCallback(func(args []js.Value) {
// 	// 	evt := args[0]
// 	// 	intVal, err := strconv.Atoi(evt.Get("target").Get("value").String())
// 	// 	if err != nil {
// 	// 		log.Println("Invalid value", err)
// 	// 		return
// 	// 	}
// 	// 	dt.SetNDots(intVal)
// 	// })
// 	// defer countChangeEvt.Release()

// 	// Event handler for speed range
// 	// speedInputEvt := js.NewCallback(func(args []js.Value) {
// 	// 	evt := args[0]
// 	// 	fval, err := strconv.ParseFloat(evt.Get("target").Get("value").String(), 64)
// 	// 	if err != nil {
// 	// 		log.Println("Invalid value", err)
// 	// 		return
// 	// 	}
// 	// 	dt.speed = fval
// 	// })
// 	// defer speedInputEvt.Release()

// 	// Event handler for size
// 	// sizeChangeEvt := js.NewCallback(func(args []js.Value) {
// 	// 	evt := args[0]
// 	// 	intVal, err := strconv.Atoi(evt.Get("target").Get("value").String())
// 	// 	if err != nil {
// 	// 		log.Println("Invalid value")
// 	// 	}
// 	// 	dt.size = intVal
// 	// })
// 	// defer sizeChangeEvt.Release()

// 	// Event handler for lines toggle
// 	// lineChangeEvt := js.NewCallback(func(args []js.Value) {
// 	// 	evt := args[0]
// 	// 	dt.lines = evt.Get("target").Get("checked").Bool()
// 	// })
// 	// defer lineChangeEvt.Release()

// 	// Event handler for dashed toggle
// 	// dashedChangeEvt := js.NewCallback(func(args []js.Value) {
// 	// 	evt := args[0]
// 	// 	dt.dashed = evt.Get("target").Get("checked").Bool()
// 	// })
// 	// defer dashedChangeEvt.Release()

// 	doc.Call("addEventListener", "mousemove", mouseMoveEvt)
// 	// doc.Call("getElementById", "count").Call("addEventListener", "change", countChangeEvt)
// 	// doc.Call("getElementById", "speed").Call("addEventListener", "input", speedInputEvt)
// 	// doc.Call("getElementById", "size").Call("addEventListener", "input", sizeChangeEvt)
// 	// doc.Call("getElementById", "dashed").Call("addEventListener", "change", dashedChangeEvt)
// 	// doc.Call("getElementById", "lines").Call("addEventListener", "change", lineChangeEvt)

// 	dt.SetNDots(100)
// 	dt.lines = false
// 	var renderFrame js.Callback
// 	var tmark float64
// 	var markCount = 0
// 	var tdiffSum float64

// 	renderFrame = js.NewCallback(func(args []js.Value) {
// 		now := args[0].Float()
// 		tdiff := now - tmark
// 		tdiffSum += now - tmark
// 		markCount++
// 		if markCount > 10 {
// 			doc.Call("getElementById", "fps").Set("innerHTML", fmt.Sprintf("FPS: %.01f", 1000/(tdiffSum/float64(markCount))))
// 			tdiffSum, markCount = 0, 0
// 		}
// 		tmark = now

// 		// Pool window size to handle resize
// 		curBodyW := doc.Get("body").Get("clientWidth").Float()
// 		curBodyH := doc.Get("body").Get("clientHeight").Float()
// 		if curBodyW != width || curBodyH != height {
// 			width, height = curBodyW, curBodyH
// 			canvasEl.Set("width", width)
// 			canvasEl.Set("height", height)
// 		}
// 		dt.Update(tdiff / 1000)

// 		js.Global().Call("requestAnimationFrame", renderFrame)
// 	})
// 	defer renderFrame.Release()

// 	// Start running
// 	js.Global().Call("requestAnimationFrame", renderFrame)

// 	<-done

// }

// // DotThing manager
// type DotThing struct {
// 	dots   []*Dot
// 	dashed bool
// 	lines  bool
// 	speed  float64
// 	size   int
// }

// // Update updates the dot positions and draws
// func (dt *DotThing) Update(dtTime float64) {
// 	if dt.dots == nil {
// 		return
// 	}
// 	ctx.Call("clearRect", 0, 0, width, height)

// 	// Update
// 	for i, dot := range dt.dots {
// 		dir := [2]float64{}
// 		// Bounce
// 		if dot.pos[0] < 0 {
// 			dot.pos[0] = 0
// 			dot.dir[0] *= -1
// 		}
// 		if dot.pos[0] > width {
// 			dot.pos[0] = width
// 			dot.dir[0] *= -1
// 		}

// 		if dot.pos[1] < 0 {
// 			dot.pos[1] = 0
// 			dot.dir[1] *= -1
// 		}

// 		if dot.pos[1] > height {
// 			dot.pos[1] = height
// 			dot.dir[1] *= -1
// 		}
// 		dir = dot.dir

// 		ctx.Set("globalAlpha", 0.5)
// 		ctx.Call("beginPath")
// 		ctx.Set("fillStyle", fmt.Sprintf("#%06x", dot.color))
// 		ctx.Set("strokeStyle", fmt.Sprintf("#%06x", dot.color))
// 		// Dashed array ref: https://github.com/golang/go/blob/release-branch.go1.11/src/syscall/js/js.go#L98
// 		ctx.Call("setLineDash", []interface{}{})
// 		if dt.dashed {
// 			ctx.Call("setLineDash", []interface{}{5, 10})
// 		}
// 		ctx.Set("lineWidth", dt.size)
// 		ctx.Call("arc", dot.pos[0], dot.pos[1], dt.size, 0, 2*math.Pi)
// 		ctx.Call("fill")

// 		mdx := mousePos[0] - dot.pos[0]
// 		mdy := mousePos[1] - dot.pos[1]
// 		d := math.Sqrt(mdx*mdx + mdy*mdy)
// 		if d < 200 {
// 			ctx.Set("globalAlpha", 1-d/200)
// 			ctx.Call("beginPath")
// 			ctx.Call("moveTo", dot.pos[0], dot.pos[1])
// 			ctx.Call("lineTo", mousePos[0], mousePos[1])
// 			ctx.Call("stroke")
// 			if d > 100 { // move towards mouse
// 				dir[0] = (mdx / d) * 2
// 				dir[1] = (mdy / d) * 2
// 			} else { // do not move
// 				dir[0] = 0
// 				dir[1] = 0
// 			}
// 		}

// 		if dt.lines {
// 			for _, dot2 := range dt.dots[i+1:] {
// 				mx := dot2.pos[0] - dot.pos[0]
// 				my := dot2.pos[1] - dot.pos[1]
// 				d := mx*mx + my*my
// 				if d < lineDistSq {
// 					ctx.Set("globalAlpha", 1-d/lineDistSq)
// 					ctx.Call("beginPath")
// 					ctx.Call("moveTo", dot.pos[0], dot.pos[1])
// 					ctx.Call("lineTo", dot2.pos[0], dot2.pos[1])
// 					ctx.Call("stroke")
// 				}
// 			}
// 		}

// 		dot.pos[0] += dir[0] * dt.speed * dtTime
// 		dot.pos[1] += dir[1] * dt.speed * dtTime
// 	}
// }

// // SetNDots reinitializes dots with n size
// func (dt *DotThing) SetNDots(n int) {
// 	dt.dots = make([]*Dot, n)
// 	for i := 0; i < n; i++ {
// 		dt.dots[i] = &Dot{
// 			pos: [2]float64{
// 				rand.Float64() * width,
// 				rand.Float64() * height,
// 			},
// 			dir: [2]float64{
// 				rand.NormFloat64(),
// 				rand.NormFloat64(),
// 			},
// 			color: uint32(rand.Intn(0xFFFFFF)),
// 			size:  10,
// 		}
// 	}
// }

// // Dot represents a dot ...
// type Dot struct {
// 	pos   [2]float64
// 	dir   [2]float64
// 	color uint32
// 	size  float64
// }
