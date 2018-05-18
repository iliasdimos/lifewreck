package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/Tarliton/collision2d"

	"github.com/hajimehoshi/ebiten"
)

var enemies map[string]*enemy

var enemySpeed float64 = 1

type enemy struct {
	x, y, height, width float64
	angle               float64
	Speed               float64
}

func initEnemies() {
	enemies = make(map[string]*enemy)
}

func NewEnemy(x, y float64) {
	e := enemy{}

	e.x, e.y = position()
	e.Speed = enemySpeed // Add random for speed +- 5
	e.angle = math.Atan2(float64(y-e.y), float64(x-e.x))
	e.width = float64(enemySprite.Bounds().Dx()) / 2
	e.height = float64(enemySprite.Bounds().Dy()) / 2
	enemies[time.Now().String()] = &e
}

func (e *enemy) Update(x, y float64) {
	e.angle = math.Atan2(float64(y-e.y), float64(x-e.x))
	if x < e.x {
		e.x -= e.Speed // * screenWidth / 500
	}
	if x > e.x {
		e.x += e.Speed // * screenWidth / 500
	}
	if y < e.y {
		e.y -= e.Speed // * screenHeight / 500
	}
	if y > e.y {
		e.y += e.Speed // * screenHeight / 500
	}
}

func (e *enemy) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Reset()
	op.GeoM.Scale(0.5, 0.5)

	op.GeoM.Translate(-float64(e.width)/2, -float64(e.height)/2)

	op.GeoM.Rotate(e.angle)

	op.GeoM.Translate(float64(e.X()), float64(e.Y()))
	screen.DrawImage(enemySprite, op)
}

func position() (x, y float64) {
	rand.Seed(time.Now().UTC().UnixNano())
	side := rand.Intn(3)

	switch side {
	case 0:
		x = 0
		y = float64(rand.Intn(screenHeight))
	case 1:
		x = float64(rand.Intn(screenWidth))
		y = 0
	case 2:
		x = screenWidth
		y = float64(rand.Intn(screenHeight))
	case 3:
		x = float64(rand.Intn(screenWidth))
		y = screenHeight
	}
	return x, y
}

// func (e *enemy) CheckCollision(s Something) bool {
// 	if image.Rect(e.x, e.y, e.x+e.isDown()*e.width, e.y+e.isUp()*e.height).Overlaps(image.Rect(s.X(), s.Y(), s.SX(), s.SY())) {
// 		fmt.Println(e.angle)
// 		return true
// 	}
// 	return false
// }

func (e *enemy) CheckCollision(s Something) bool {
	v1 := collision2d.Vector{X: float64(s.X()), Y: float64(s.Y())}
	c2 := collision2d.Circle{Pos: collision2d.Vector{X: float64(e.X()), Y: float64(e.Y())}, R: float64(e.SX() - e.X())}
	return collision2d.PointInCircle(v1, c2)

	// p2 := collision2d.NewPolygon(collision2d.Vector{X: float64(e.X()), Y: float64(e.Y())}, collision2d.Vector{X: float64(e.SX()), Y: float64(e.SY())}, e.Angle(), []float64{float64(e.X()), float64(e.Y()), float64(e.SX()), float64(e.SY())})

	// p2 := collision2d.NewPolygon(collision2d.Vector{X: float64(e.X()), Y: float64(e.Y())}, collision2d.Vector{X: float64(e.X()), Y: float64(e.Y())}, e.Angle(), []float64{float64(e.X()), float64(e.Y()), float64(e.SX()), float64(e.SY())})
	// return collision2d.PointInPolygon(v1, p2)

}

func (e *enemy) isUp() int {
	if e.angle > 0 {
		return -1
	}
	return 1
}

func (e *enemy) isDown() int {
	if e.angle < math.Pi/2 || e.angle > -math.Pi/2 {
		return -1
	}
	return 1
}

func (e enemy) X() float64     { return e.x }
func (e enemy) Y() float64     { return e.y }
func (e enemy) SY() float64    { return e.y + e.height }
func (e enemy) SX() float64    { return e.x + e.width }
func (e enemy) Angle() float64 { return e.angle }
