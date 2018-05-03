package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"time"

	"github.com/jakecoffman/cp"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	// rplatformer "github.com/hajimehoshi/ebiten/examples/resources/images/platformer"
)

const (
	// Settings
	screenWidth  = 1024
	screenHeight = 512
)

var (
	playerSprite *ebiten.Image
	space        *cp.Space
	dot          *ebiten.Image
)

func init() {

	// Preload images
	image.RegisterFormat("png", "PNG", png.Decode, png.DecodeConfig)
	playerFile, err := os.Open("player.png")
	if err != nil {
		log.Fatal(err)
	}
	defer playerFile.Close()

	img, _, err := image.Decode(playerFile)
	if err != nil {
		log.Fatal("Could not open file: ", err)
	}

	playerSprite, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	width, heigth := playerSprite.Size()
	charWidth = width / 2
	charHeigth = heigth / 2

	space = cp.NewSpace()
	space.Iterations = 1

	dot, _ = ebiten.NewImage(5, 5, ebiten.FilterNearest)
	dot.Fill(color.White)
	speed = 300

	Player = player{
		X:      screenWidth / 2,
		Y:      screenHeight / 2,
		Heigth: heigth / 2,
		Width:  width / 2,
		Sprite: playerSprite}

}

type bullet struct {
	x, y  int
	angle float64
}

var (
	charX      = screenWidth / 2
	charY      = screenHeight / 2
	charHeigth int
	charWidth  int
	m          mouse
	count      int
	bullets    []bullet
	speed      float64
	lastBullet time.Time
	Player     player
)

type mouse struct{ X, Y int }

type player struct {
	X, Y, Heigth, Width int
	Sprite              *ebiten.Image
}

func (p *player) Move() {
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.X -= 3
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.X += 3
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		p.Y -= 3
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		p.Y += 3
	}
}

func update(screen *ebiten.Image) error {

	// Controls
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		charX -= 3
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		charX += 3
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		charY -= 3
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		charY += 3
	}

	// Find mouse position
	m.X, m.Y = ebiten.CursorPosition()
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("Mouse x: %d Mouse y: %d", m.X, m.Y))

	// Find angle to mouse cursor
	angle := math.Atan2(float64(m.Y-charY), float64(m.X-charX))

	// Create bullets if mouse is pressed
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if time.Now().Sub(lastBullet) > 1*time.Second {
			lastBullet = time.Now()
			b := bullet{x: charX, y: charY, angle: angle}
			bulletShape := makeBullet(float64(b.x), float64(b.y))

			bulletBody := space.AddBody(bulletShape.Body())
			bulletBody.SetAngle(b.angle)
			bulletBody.SetPositionUpdateFunc(func(body *cp.Body, dt float64) {
				v := cp.Vector{
					X: body.Position().X + (math.Cos(body.Angle()) * speed * dt),
					Y: body.Position().Y + (math.Sin(body.Angle()) * speed * dt),
				}
				body.SetPosition(v)
			})
		}

	}
	// Sides
	if charX > screenWidth-charHeigth {
		charX = screenWidth - charHeigth
	}
	if charX < 0+charHeigth {
		charX = 0 + charHeigth
	}
	if charY > screenHeight-charHeigth {
		charY = screenHeight - charHeigth
	}
	if charY < 0+charHeigth {
		charY = 0 + charHeigth
	}

	if ebiten.IsRunningSlowly() {
		return nil
	}

	// Draws Background Image
	playerOptions := &ebiten.DrawImageOptions{}

	// Draws Player sprite image
	playerOptions = &ebiten.DrawImageOptions{}
	playerOptions.GeoM.Scale(0.5, 0.5)
	// Translate based on the image's size, on the upper left side of screen
	playerOptions.GeoM.Translate(-float64(charWidth)/2, -float64(charHeigth)/2)
	// Rotate the image. As a result, the anchor point of this rotate is
	// the center of the image.
	playerOptions.GeoM.Rotate(angle)
	// Translate on current position
	playerOptions.GeoM.Translate(float64(charX), float64(charY))
	// Draw on screen, the sprite with options
	screen.DrawImage(playerSprite, playerOptions)

	// Draw bullets
	space.Step(1.0 / ebiten.FPS)

	// drawBullets(screen)

	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(200.0/255.0, 200.0/255.0, 200.0/255.0, 1)

	space.EachBody(func(body *cp.Body) {
		if body.Position().X > screenWidth || body.Position().X < 0 {
			return
		}
		if body.Position().Y > screenHeight || body.Position().Y < 0 {
			return
		}
		op.GeoM.Reset()
		op.GeoM.Translate(body.Position().X, body.Position().Y)
		screen.DrawImage(dot, op)
	})

	// FPS counter
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %f", ebiten.CurrentFPS()))

	return nil

}

// func createBullet(x, y int, angle float64) {
// 	bullets = append(bullets, bullet{x: x, y: y, angle: angle})
// }

// toRadian Converts a float coordinate to radian
// func toRadian(c float64) float64 {
// 	return c * (math.Pi / 180)
// }

func main() {

	if err := ebiten.Run(update, screenWidth, screenHeight, 1, "Lifewreck"); err != nil {
		panic(err)
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
