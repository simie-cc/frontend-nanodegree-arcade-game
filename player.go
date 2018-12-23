package main

import (
	"fmt"
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
	return rect.pointInRect(rect2.x1, rect2.y1) ||
		rect.pointInRect(rect2.x1, rect2.y2) ||
		rect.pointInRect(rect2.x2, rect2.y1) ||
		rect.pointInRect(rect2.x2, rect2.y2)
}

func (rect *Rect) pointInRect(x, y int) bool {
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
}

func NewPlayer() *Player {
	return &Player{
		cx: 2,
		cy: 5,
	}
}

func (p *Player) GetXY() (int, int) {
	return p.cx * CELL_WIDTH, p.cy*CELL_HEIGHT - 20
}

func (p *Player) GetCollisionRect() *Rect {
	px, py := p.cx*CELL_WIDTH, p.cy*CELL_HEIGHT+IMAGE_Y_SHIFT
	return &Rect{
		x1: px + 10, y1: py + 10,
		x2: px + CELL_WIDTH - 10, y2: py + CELL_HEIGHT - 10,
	}
}

func (p *Player) move(dir DIRECTION) {
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
	fmt.Println("Player pos", p.cx, p.cy)
}

func (p *Player) isCollision(another CollisionAble) bool {
	anoRect := another.GetCollisionRect()

	return p.GetCollisionRect().isIntersect(anoRect)
}
