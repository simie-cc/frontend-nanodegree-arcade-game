//Wasming
// compile: GOOS=js GOARCH=wasm go build -o main.wasm ./main.go
package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type GAMEMODE int8

const (
	GAMEMODE_INTRO GAMEMODE = iota
	GAMEMODE_MENU
	GAMEMODE_PLAY
	GAMEMODE_PLAY_READONLY
)

const (
	width  int = 505
	height int = 606

	CELL_WIDTH  = 101
	CELL_HEIGHT = 83

	ROW_COUNT = 6
	COL_COUNT = 5

	IMAGE_Y_SHIFT = 50

	IMAGE_WATER     = "images/water-block.png"
	IMAGE_STONE     = "images/stone-block.png"
	IMAGE_GRASS     = "images/grass-block.png"
	IMAGE_BOY       = "images/char-boy.png"
	IMAGE_CAT_GIRL  = "images/char-cat-girl.png"
	IMAGE_PINK_GIRL = "images/char-pink-girl.png"
	IMAGE_HORN_GIRL = "images/char-horn-girl.png"
	IMAGE_PRIN_GIRL = "images/char-princess-girl.png"

	IMAGE_BUG   = "images/enemy-bug.png"
	IMAGE_ARROW = "images/arrow.png"

	DEBUG = false
)

var (
	CharImages = [...]string{
		IMAGE_BOY, IMAGE_CAT_GIRL, IMAGE_PINK_GIRL, IMAGE_HORN_GIRL, IMAGE_PRIN_GIRL,
	}
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

	player     *Player = nil
	enemies    []*Enemy
	enemyCount = 1

	gameMode          GAMEMODE
	charSelected      = IMAGE_PINK_GIRL
	charSelectedIndex = 0
	menuImageRects    []Rect

	verboseLevel = 1
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	fmt.Println("Loading...")

	prepareResources()
	StartFpsCounting()
	renderer.ListenClickEvent(handleClickEvent)
	renderer.ListenKeyboardEvent(handleKeyboardEvent)
	renderer.ListenMouseMoveEvent(handleMouseMoveEvent)

	systemService.RegisterGetVerboseLevel(func() int {
		return verboseLevel
	})
	systemService.RegisterSetVerboseLevel(func(newLevel int) {
		verboseLevel = newLevel
	})

	fmt.Println("Start game...")
	startIntro()
	//startGame(3)

	<-done

	fmt.Println("Main exit")
}

func prepareResources() {
	renderer.Init(width, height)
	renderer.PrepareImage(IMAGE_WATER)
	renderer.PrepareImage(IMAGE_STONE)
	renderer.PrepareImage(IMAGE_GRASS)
	renderer.PrepareImage(IMAGE_BOY)
	renderer.PrepareImage(IMAGE_CAT_GIRL)
	renderer.PrepareImage(IMAGE_PINK_GIRL)
	renderer.PrepareImage(IMAGE_HORN_GIRL)
	renderer.PrepareImage(IMAGE_PRIN_GIRL)
	renderer.PrepareImage(IMAGE_BUG)
	renderer.PrepareImage(IMAGE_ARROW)
	renderer.WaitImage()
}

func startIntro() {
	gameMode = GAMEMODE_INTRO
	fmt.Println("Gamemode to", gameMode)

	renderer.RegisterRenderFunction(func() {
		renderer.ClearRect()
		renderGameTitle("這是一個小遊戲", "點擊滑鼠開始...")

		AddFps()
	})
	renderer.StartRender()
}

func renderGameTitle(title, subtitle string) {
	renderer.SetFont("50px '微軟正黑體'")
	renderer.SetTextAlign("center")

	centerX := width / 2
	renderer.DrawText(title, centerX, 200)
	renderer.DrawRoundedRect(&Rect{x1: 30, y1: 100, x2: width - 30, y2: 250})

	if len(subtitle) > 0 {
		renderer.SetFont("30px '微軟正黑體'")
		renderer.DrawText(subtitle, centerX, 400)
	}
}

func startMenu() {

	gameMode = GAMEMODE_MENU
	fmt.Println("Gamemode to", gameMode)

	menuImageRects = make([]Rect, len(CharImages))

	for idx, _ := range CharImages {
		cx, cy := idx*CELL_WIDTH, 300
		menuImageRects[idx] = Rect{x1: cx, y1: cy,
			x2: cx + CELL_WIDTH - 1, y2: cy + 171}
	}

	renderer.RegisterRenderFunction(func() {
		renderer.ClearRect()
		renderGameTitle("這是一個小遊戲", "")
		renderOptions()

		renderer.SetFont("30px '微軟正黑體'")
		centerX := width / 2
		renderer.DrawText("點按Enter選擇角色", centerX, 550)

		AddFps()
	})
	renderer.StartRender()
}

func renderOptions() {
	for idx, img := range CharImages {
		renderer.DrawImage(img, idx*CELL_WIDTH, 300)
	}
	renderer.DrawImage(IMAGE_ARROW,
		charSelectedIndex*CELL_WIDTH+(CELL_WIDTH/2-16),
		460)
}

func roundToInt(f float64) int {
	return int(math.Round(f))
}

func startGame(withEnemyCount int) {

	gameMode = GAMEMODE_PLAY
	enemyCount = withEnemyCount

	charSelected = CharImages[charSelectedIndex]
	player = NewPlayer()
	enemies = make([]*Enemy, enemyCount)
	for idx, _ := range enemies {
		enemies[idx] = NewEnemy()
	}

	var timeTick = time.Now()
	var failMessage = ""
	var messageTimeout time.Time

	renderer.RegisterRenderFunction(func() {

		if !shouldRenderNextFrame() {
			return
		}

		dt := time.Since(timeTick)
		timeTick = time.Now()

		renderer.ClearRect()

		drawBaseGround()
		renderer.SetFont("16px '微軟正黑體'")
		renderer.DrawText(fmt.Sprintf("難度: %d", enemyCount), 50, 20)

		drawEntity(player, charSelected)
		renderer.DrawRect(player.GetCollisionRect())
		if len(failMessage) == 0 && player.IsReachEdge() {
			gameMode = GAMEMODE_PLAY_READONLY
			failMessage = "過關！"
			messageTimeout = time.Now()
			enemyCount += 2
			return
		}

		for _, enemy := range enemies {
			enemy.tick(dt)
			drawEntity(enemy, IMAGE_BUG)
			renderer.DrawRect(enemy.GetCollisionRect())

			if len(failMessage) == 0 && player.isCollision(enemy) {
				fmt.Println("Collision!!")
				gameMode = GAMEMODE_PLAY_READONLY
				failMessage = "已撞到啦！"
				messageTimeout = time.Now()

				return
			}
		}

		if len(failMessage) > 0 {
			renderGameTitle(failMessage, "")
			duration := time.Since(messageTimeout)
			if duration > 1*time.Second {
				renderer.StopRender()
				startGame(enemyCount)
				return
			}
		}

		AddFps()
	})

	renderer.StartRender()
}

func drawBaseGround() {
	for r := 0; r < ROW_COUNT; r++ {
		for c := 0; c < COL_COUNT; c++ {
			x, y := c*CELL_WIDTH, r*CELL_HEIGHT
			renderer.DrawImage(RowImages[r], x, y)
		}
	}
}

func drawEntity(p Positional, charurl string) {
	x, y := p.GetXY()
	renderer.DrawImage(charurl, x, y)
}

func handleClickEvent(eventType string, x, y int) {
	if gameMode == GAMEMODE_INTRO {
		startMenu()
	} else if gameMode == GAMEMODE_MENU {
		for _, imageRect := range menuImageRects {
			if imageRect.PointInRect(x, y) {
				startGame(3)
				return
			}
		}
	}
}

func handleMouseMoveEvent(eventType string, x, y int) {
	if gameMode != GAMEMODE_MENU {
		return
	}

	for idx, imageRect := range menuImageRects {
		if imageRect.PointInRect(x, y) {
			charSelectedIndex = idx
			return
		}
	}
}

func handleKeyboardEvent(eventType, key string) {
	switch gameMode {
	case GAMEMODE_MENU:
		handleKeyboardEventMenu(eventType, key)
	case GAMEMODE_PLAY:
		handleKeyboardEventPlaying(eventType, key)
	}
}

func handleKeyboardEventMenu(eventType, key string) {

	switch key {

	case "ArrowLeft":
		charSelectedIndex--
		if charSelectedIndex < 0 {
			charSelectedIndex = len(CharImages) - 1
		}

	case "ArrowRight":
		charSelectedIndex++
		if charSelectedIndex >= len(CharImages) {
			charSelectedIndex = 0
		}

	case "Enter":
		startGame(3)
	}
}

func handleKeyboardEventPlaying(eventType, key string) {
	if player == nil {
		return
	}

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
