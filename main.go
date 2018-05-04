//go:generate file2byteslice -package=images -input=images/player.png -output=./images/player.go -var=Player_png
package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"time"

	"github.com/dosko64/lifewreck/images"
	"github.com/jakecoffman/cp"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	// Settings
	screenWidth  = 1024
	screenHeight = 512
)

var (
	space  *cp.Space
	Player player
)

func init() {

	// Preload images
	img, _, err := image.Decode(bytes.NewReader(images.Player_png))
	if err != nil {
		log.Fatal("Could not open file: ", err)
	}

	playerSprite, err := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	if err != nil {
		log.Fatal("Could not create player sprite")
	}

	width, heigth := playerSprite.Size()

	space = cp.NewSpace()
	space.Iterations = 1

	dot, err := ebiten.NewImage(5, 5, ebiten.FilterNearest)
	if err != nil {
		log.Fatal("Could not create dot image")
	}
	dot.Fill(color.White)

	Player = player{
		X:            screenWidth / 2,
		Y:            screenHeight / 2,
		Heigth:       heigth / 2,
		Width:        width / 2,
		Sprite:       playerSprite,
		LastShoot:    time.Now().Add(-1 * time.Second),
		BulletSpeed:  300,
		BulletSprite: dot,
		MoveSpeed:    3,
	}

}

func update(screen *ebiten.Image) error {

	Player.Move()

	// Find mouse position
	x, y := ebiten.CursorPosition()

	// Find angle to mouse cursor
	angle := math.Atan2(float64(y-Player.Y), float64(x-Player.X))
	Player.Angle = angle
	// Check if mouse is pressed and shoot
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		Player.Shoot()
	}

	// Check if we are near screen sides
	Player.CheckScreen(screenWidth, screenHeight)

	if ebiten.IsRunningSlowly() {
		return nil
	}

	Player.Draw(screen)
	// Draw bullets
	space.Step(1.0 / ebiten.FPS)

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
		screen.DrawImage(Player.BulletSprite, op)
	})

	// FPS counter
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %f", ebiten.CurrentFPS()))
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("Mouse x: %d Mouse y: %d", m.X, m.Y))

	return nil

}

func main() {

	if err := ebiten.Run(update, screenWidth, screenHeight, 1, "Lifewreck"); err != nil {
		panic(err)
	}
}
