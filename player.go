package main

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/jakecoffman/cp"
)

type player struct {
	X, Y, Heigth, Width int
	Angle               float64
	Sprite              *ebiten.Image
	LastShoot           time.Time
	BulletSpeed         float64
	BulletSprite        *ebiten.Image
	MoveSpeed           int
}

func (p *player) Move() {
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.X -= p.MoveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.X += p.MoveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		p.Y -= p.MoveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		p.Y += p.MoveSpeed
	}
}

func (p *player) CheckScreen(screenWidth, screenHeight int) {
	if p.X > screenWidth-p.Width {
		p.X = screenWidth - p.Width
	}
	if p.X < 0+p.Width {
		p.X = 0 + p.Width
	}
	if p.Y > screenHeight-p.Heigth {
		p.Y = screenHeight - p.Heigth
	}
	if p.Y < 0+p.Heigth {
		p.Y = 0 + p.Heigth
	}
}

func (p player) Draw(screen *ebiten.Image) {
	// Draws Background Image

	// Draws Player sprite image
	playerOptions := &ebiten.DrawImageOptions{}
	playerOptions.GeoM.Scale(0.5, 0.5)
	// Translate based on the image's size, on the upper left side of screen
	playerOptions.GeoM.Translate(-float64(Player.Width)/2, -float64(Player.Heigth)/2)
	// Rotate the image. As a result, the anchor point of this rotate is
	// the center of the image.
	playerOptions.GeoM.Rotate(p.Angle)
	// Translate on current position
	playerOptions.GeoM.Translate(float64(p.X), float64(p.Y))
	// Draw on screen, the sprite with options
	screen.DrawImage(p.Sprite, playerOptions)

}

func (p *player) Shoot() {
	// Create bullets if mouse is pressed
	if time.Now().Sub(p.LastShoot) > 1*time.Second {
		p.LastShoot = time.Now()

		bulletShape := makeBullet(float64(p.X), float64(p.Y))
		bulletBody := space.AddBody(bulletShape.Body())
		bulletBody.SetAngle(p.Angle)
		bulletBody.SetPositionUpdateFunc(func(body *cp.Body, dt float64) {
			v := cp.Vector{
				X: body.Position().X + (math.Cos(body.Angle()) * p.BulletSpeed * dt),
				Y: body.Position().Y + (math.Sin(body.Angle()) * p.BulletSpeed * dt),
			}
			body.SetPosition(v)
		})
	}
}

func makeBullet(x, y float64) *cp.Shape {
	body := cp.NewBody(1.0, cp.INFINITY)
	body.SetPosition(cp.Vector{x, y})

	shape := cp.NewCircle(body, 0.95, cp.Vector{})
	shape.SetElasticity(0)
	shape.SetFriction(0)

	return shape
}
