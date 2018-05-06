package main

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/jakecoffman/cp"
)

var playerSprite *ebiten.Image

type player struct {
	X, Y, Heigth, Width int
	Angle               float64
	LastShoot           time.Time
	BulletSpeed         float64
	Speed               int
}

func (p *player) Move() {
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.X -= p.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.X += p.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		p.Y -= p.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		p.Y += p.Speed
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

func (p *player) Draw(screen *ebiten.Image) {
	// Draws Background Image

	// Draws Player sprite image
	playerOptions := &ebiten.DrawImageOptions{}
	playerOptions.GeoM.Scale(0.5, 0.5)
	// Translate based on the image's size, on the upper left side of screen
	playerOptions.GeoM.Translate(-float64(p.Width)/2, -float64(p.Heigth)/2)
	// Rotate the image. As a result, the anchor point of this rotate is
	// the center of the image.
	playerOptions.GeoM.Rotate(p.Angle)
	// Translate on current position
	playerOptions.GeoM.Translate(float64(p.X), float64(p.Y))
	// Draw on screen, the sprite with options
	screen.DrawImage(playerSprite, playerOptions)

}

func (p *player) Shoot(space *cp.Space, screen *ebiten.Image) {
	// Create bullets if mouse is pressed
	p.LastShoot = time.Now()

	bulletShape := makeBullet(float64(p.X), float64(p.Y))
	// bulletShape.SetCollisionType(BulletCollisionType)

	bulletBody := space.AddBody(bulletShape.Body())
	bulletBody.UserData = "bullet"
	bulletBody.SetAngle(p.Angle)
	bulletBody.SetPositionUpdateFunc(func(body *cp.Body, dt float64) {
		v := cp.Vector{
			X: body.Position().X + (math.Cos(body.Angle()) * p.BulletSpeed * dt),
			Y: body.Position().Y + (math.Sin(body.Angle()) * p.BulletSpeed * dt),
		}
		body.SetPosition(v)
	})
	// space.AddBody(bulletBody)
	// space.AddShape(bulletShape)
}

const BulletCollisionType cp.CollisionType = 0

func makeBullet(x, y float64) *cp.Shape {
	body := cp.NewBody(10.0, cp.BODY_KINEMATIC)
	body.SetPosition(cp.Vector{x, y})

	shape := cp.NewCircle(body, 0.95, cp.Vector{})
	shape.SetElasticity(0)
	shape.SetFriction(0)

	shape.SetCollisionType(BulletCollisionType)

	return shape
}
