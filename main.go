//go:generate file2byteslice -package=images -input=images/player.png -output=./images/player.go -var=Player_png
//go:generate file2byteslice -package=images -input=images/zombie-animation.png -output=./images/zombie.go -var=Zombie_png
//go:generate file2byteslice -package=images -input=images/survivor1_gun.png -output=./images/survivor.go -var=Survivor_png
//go:generate file2byteslice -package=images -input=images/zombie1_hold.png -output=./images/zombie1.go -var=Zombie1_png

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
var bulletSprite *ebiten.Image
var enemySprite *ebiten.Image

func init() {

	// Preload images
	img, _, err := image.Decode(bytes.NewReader(images.Survivor_png))
	if err != nil {
		log.Fatal("Could not open file: ", err)
	}

	playerSprite, err = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	if err != nil {
		log.Fatal("Could not create player sprite")
	}

	width, heigth := playerSprite.Size()

	space = cp.NewSpace()
	space.Iterations = 1

	bulletSprite, err = ebiten.NewImage(5, 5, ebiten.FilterNearest)
	if err != nil {
		log.Fatal("Could not create dot image: ", err)
	}
	bulletSprite.Fill(color.White)

	Player = player{
		X:           screenWidth / 2,
		Y:           screenHeight / 2,
		Heigth:      heigth / 2,
		Width:       width / 2,
		LastShoot:   time.Now().Add(-1 * time.Second),
		BulletSpeed: 300,
		Speed:       3,
	}

	img, _, err = image.Decode(bytes.NewReader(images.Zombie1_png))
	if err != nil {
		log.Fatal("Could not open file: ", err)
	}

	enemySprite, err = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	if err != nil {
		log.Fatal("Could not create enemy sprite")
	}

	space.SetCollisionSlop(0.5)
	ch := space.NewCollisionHandler(BulletCollisionType, EnemyCollisionType)
	ch.BeginFunc = bulletEnemy

}

var lastEnemy = time.Now()

func update(screen *ebiten.Image) error {

	Player.Move()

	// Find mouse position
	x, y := ebiten.CursorPosition()

	// Find angle to mouse cursor
	Player.Angle = math.Atan2(float64(y-Player.Y), float64(x-Player.X))

	// Check if mouse is pressed and shoot
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if time.Now().Sub(Player.LastShoot) > 1*time.Second {
			Player.Shoot(space, screen)
		}
	}

	// Check if we are near screen sides
	Player.CheckScreen(screenWidth, screenHeight)

	if ebiten.IsRunningSlowly() {
		return nil
	}

	Player.Draw(screen)
	// Draw bullets
	space.Step(1.0 / ebiten.FPS)

	if time.Now().Sub(lastEnemy) > 1*time.Second {
		NewEnemy(&Player, space)
		lastEnemy = time.Now()
	}

	space.EachBody(func(body *cp.Body) {
		if body.Position().X > screenWidth || body.Position().X < 0 {
			space.RemoveBody(body)
			return
		}
		if body.Position().Y > screenHeight || body.Position().Y < 0 {
			space.RemoveBody(body)
			return
		}
		if body.UserData == "bullet" {
			op := &ebiten.DrawImageOptions{}
			// op.ColorM.Scale(200.0/255.0, 200.0/255.0, 200.0/255.0, 1)

			op.GeoM.Reset()
			op.GeoM.Translate(body.Position().X, body.Position().Y)
			screen.DrawImage(bulletSprite, op)
		}
		if body.UserData == "enemy" {
			enemyOptions := &ebiten.DrawImageOptions{}
			enemyOptions.GeoM.Reset()

			enemyOptions.GeoM.Scale(0.5, 0.5)

			body.SetAngle(math.Atan2(float64(Player.Y)-body.Position().Y, float64(Player.X)-body.Position().X))
			enemyOptions.GeoM.Rotate(body.Angle())

			enemyOptions.GeoM.Translate(body.Position().X, body.Position().Y)
			screen.DrawImage(enemySprite, enemyOptions)
		}

	})

	// FPS counter
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %f", ebiten.CurrentFPS()))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\nMouse x: %d Mouse y: %d", x, y))

	return nil
}

func main() {
	if err := ebiten.Run(update, screenWidth, screenHeight, 1, "Lifewreck"); err != nil {
		panic(err)
	}
}

func bulletEnemy(arb *cp.Arbiter, space *cp.Space, data interface{}) bool {

	// Get pointers to the two bodies in the collision pair and define local variables for them.
	// Their order matches the order of the collision types passed
	// to the collision handler this function was defined for

	bullet, enemy := arb.Bodies()
	// bullet.RemoveShape()
	// additions and removals can't be done in a normal callback.
	// Schedule a post step callback to do it.
	// Use the hook as the key and pass along the arbiter.

	space.RemoveBody(bullet)

	space.RemoveBody(enemy)

	// space.AddPostStepCallback(attachHook, bullet, enemy)

	return true
}
