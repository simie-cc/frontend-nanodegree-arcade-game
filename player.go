package main

type Positional interface {
	GetXY() (int, int)
}

type Player struct {
	Positional
	cx, cy int
}

func NewPlayer() *Player {
	return &Player{
		cx: 3,
		cy: 3,
	}
}

func (p *Player) GetXY() (int, int) {
	return p.cx * CELL_WIDTH, p.cy * CELL_HEIGHT
}

type Enemy struct {
	Positional
	x, y int
}

func NewEnemy() *Enemy {
	return &Enemy{
		x: 20,
		y: 50,
	}
}

func (en *Enemy) GetXY() (int, int) {
	return en.x, en.y
}
