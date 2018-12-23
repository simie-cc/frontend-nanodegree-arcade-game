package main

import (
	"fmt"
	"math"
	"time"
)

type DIRECTION int

const (
	DIRECTION_UP DIRECTION = iota
	DIRECTION_DOWN
	DIRECTION_LEFT
	DIRECTION_RIGHT
)

type Positional interface {
	GetXY() (int, int)
}

type Rect struct {
	x1, y1 int
	x2, y2 int
}

func (rect *Rect) isIntersect(rect2 *Rect) bool {
	return rect.PointInRect(rect2.x1, rect2.y1) ||
		rect.PointInRect(rect2.x1, rect2.y2) ||
		rect.PointInRect(rect2.x2, rect2.y1) ||
		rect.PointInRect(rect2.x2, rect2.y2)
}

func (rect *Rect) PointInRect(x, y int) bool {
	return rect.x1 <= x && x <= rect.x2 &&
		rect.y1 <= y && y <= rect.y2
}

type CollisionAble interface {
	GetCollisionRect() *Rect
}

type Player struct {
	Positional
	CollisionAble
	cx, cy int
	shake  *Shaker
}

func NewPlayer() *Player {
	p := &Player{
		cx: 2,
		cy: 5,
	}
	p.shake = nil
	return p
}

func (p *Player) GetXY() (int, int) {
	sx, sy := 0, 0
	if p.shake != nil {
		if p.shake.IsEnd() {
			p.shake = nil
		} else {
			sx, sy = p.shake.GetXY()
		}
	}
	return p.cx*CELL_WIDTH + sx, p.cy*CELL_HEIGHT - 20 + sy
}

func (p *Player) GetCollisionRect() *Rect {
	px, py := p.cx*CELL_WIDTH, p.cy*CELL_HEIGHT+IMAGE_Y_SHIFT
	return &Rect{
		x1: px + 10, y1: py + 10,
		x2: px + CELL_WIDTH - 10, y2: py + CELL_HEIGHT - 10,
	}
}

func (p *Player) IsReachEdge() bool {
	return p.cy == 0
}

func (p *Player) move(dir DIRECTION) {
	rx, ry := p.cx, p.cy

	switch dir {
	case DIRECTION_UP:
		p.cy--
	case DIRECTION_DOWN:
		p.cy++
	case DIRECTION_LEFT:
		p.cx--
	case DIRECTION_RIGHT:
		p.cx++
	}

	if p.isOutOfRange() {
		p.cx, p.cy = rx, ry
		p.shake = NewShaker(500 * time.Millisecond)
	}

	fmt.Println("Player pos", p.cx, p.cy)
}

func (p *Player) isOutOfRange() bool {
	return p.cx < 0 || p.cx > COL_COUNT-1 ||
		p.cy < 0 || p.cy > ROW_COUNT-1
}

func (p *Player) isCollision(another CollisionAble) bool {
	anoRect := another.GetCollisionRect()

	return p.GetCollisionRect().isIntersect(anoRect)
}

type Shaker struct {
	start    time.Time
	duration time.Duration
}

func NewShaker(duration time.Duration) *Shaker {
	return &Shaker{
		start:    time.Now(),
		duration: duration,
	}
}

func (sh *Shaker) GetXY() (int, int) {
	const shakeDiff = float64(5)
	since := time.Since(sh.start)
	shakeTick := (float64(since) / float64(time.Second)) * 72
	xdiff := shakeDiff * math.Sin((math.Pi/float64(12))*shakeTick)
	return int(xdiff), 0
}

func (sh *Shaker) IsEnd() bool {
	since := time.Since(sh.start)
	return since > sh.duration
}
