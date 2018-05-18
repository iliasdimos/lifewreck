package main

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten"
)

type bullet struct {
	x, y    float64
	angle   float64
	deleted bool
}

var bullets map[string]*bullet

var bulletSpeed float64 = 10.0

func initBullets() {
	bullets = make(map[string]*bullet)
}
func NewBullet(x, y float64, angle float64) {
	b := bullet{x: x, y: y, angle: angle}
	bullets[time.Now().String()] = &b
	// return &b
}

func (b *bullet) Update() {
	b.x = b.x + math.Cos(b.angle)*float64(bulletSpeed)
	b.y = b.y + math.Sin(b.angle)*float64(bulletSpeed)

	if b.x < 0 || b.x > screenWidth {
		b.deleted = true
	}
	if b.y < 0 || b.y > screenHeight {
		b.deleted = true
	}
}

func (b *bullet) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(200.0/255.0, 200.0/255.0, 200.0/255.0, 1)

	op.GeoM.Reset()
	op.GeoM.Translate(float64(b.x), float64(b.y))
	screen.DrawImage(bulletSprite, op)

}

func CleanBullets() {
	for i, b := range bullets {
		if b.deleted == true {
			delete(bullets, i)
		}
	}
}

func (b *bullet) X() float64  { return b.x }
func (b *bullet) Y() float64  { return b.y }
func (b *bullet) SY() float64 { return b.x + 2 }
func (b *bullet) SX() float64 { return b.y + 2 }
