package main

import (
	"time"

	"github.com/hajimehoshi/ebiten"
)

var playerSprite *ebiten.Image

type player struct {
	x, y, height, width float64
	Angle               float64
	LastShoot           time.Time
	BulletSpeed         float64
	Speed               float64
}

func (p *player) Move() {
	if len(ebiten.TouchIDs()) > 0 {
		p.checkTouch()
		return
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.x -= p.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.x += p.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		p.y -= p.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		p.y += p.Speed
	}
}

func (p *player) checkTouch() {
	for _, t := range ebiten.TouchIDs() {
		x, y := ebiten.TouchPosition(t)
		if x > int(p.x) {
			p.x += p.Speed
		}
		if x < int(p.x) {
			p.x -= p.Speed
		}
		if y > int(p.y) {
			p.y += p.Speed
		}
		if y < int(p.y) {
			p.y -= p.Speed
		}
	}
	return
}

func (p *player) CheckScreen(screenWidth, screenHeight float64) {
	if p.x > screenWidth-p.width {
		p.x = screenWidth - p.width
	}
	if p.x < 0+p.width {
		p.x = 0 + p.width
	}
	if p.y > screenHeight-p.height {
		p.y = screenHeight - p.height
	}
	if p.y < 0+p.height {
		p.y = 0 + p.height
	}
}

func (p *player) Draw(screen *ebiten.Image) {
	// Draws Background Image

	// Draws Player sprite image
	playerOptions := &ebiten.DrawImageOptions{}
	playerOptions.GeoM.Scale(0.5, 0.5)
	// Translate based on the image's size, on the upper left side of screen
	playerOptions.GeoM.Translate(-float64(p.width)/2, -float64(p.height)/2)
	// Rotate the image. As a result, the anchor point of this rotate is
	// the center of the image.
	playerOptions.GeoM.Rotate(p.Angle)
	// Translate on current position
	playerOptions.GeoM.Translate(float64(p.x), float64(p.y))
	// Draw on screen, the sprite with options
	screen.DrawImage(playerSprite, playerOptions)

}

func (p *player) Shoot() {
	// Create bullets if mouse is pressed
	p.LastShoot = time.Now()
	NewBullet(p.x, p.y, p.Angle)
}

func (p player) X() float64  { return p.x }
func (p player) Y() float64  { return p.y }
func (p player) SY() float64 { return p.y + p.height }
func (p player) SX() float64 { return p.x + p.width }
