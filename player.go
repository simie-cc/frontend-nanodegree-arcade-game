package main

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

type Player struct {
	Positional
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
}

type Enemy struct {
	Positional
	x, y int
}

func NewEnemy() *Enemy {
	//randomRow :=
	return &Enemy{
		x: 20,
		y: 50,
	}
}

func (en *Enemy) GetXY() (int, int) {
	return en.x, en.y
}
