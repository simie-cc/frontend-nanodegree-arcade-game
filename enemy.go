package main

import (
	"math/rand"
	"time"
)

var EnemySpeeds = [...]float64{50, 200, 300}

type Enemy struct {
	Positional
	x, y            int
	secondSpeedByPx float64
}

func NewEnemy() *Enemy {

	en := &Enemy{
		x: 0,
		y: 0,
	}
	en.randomAll()
	return en
}

func (en *Enemy) randomAll() {
	randomStart := rand.Intn(8) + 1
	en.x = -randomStart * CELL_WIDTH

	randomRow := rand.Intn(3) + 1
	en.y = randomRow * CELL_HEIGHT

	randomSpeedIndex := rand.Intn(len(EnemySpeeds))
	en.secondSpeedByPx = EnemySpeeds[randomSpeedIndex]
}

func (en *Enemy) GetXY() (int, int) {
	return en.x, en.y - 20
}

func (en *Enemy) GetCollisionRect() *Rect {
	px, py := en.x, en.y+IMAGE_Y_SHIFT
	return &Rect{
		x1: px + 10, y1: py + 10,
		x2: px + CELL_WIDTH - 10, y2: py + CELL_HEIGHT - 10,
	}
}

func (en *Enemy) tick(dt time.Duration) {
	shift := float64(en.secondSpeedByPx * (float64(dt) / float64(time.Second)))
	en.x += int(shift)

	if en.x > width {
		en.randomAll()
	}
}
