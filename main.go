//go:generate file2byteslice -package=images -input=images/survivor1_gun.png -output=./images/survivor.go -var=Survivor_png
//go:generate file2byteslice -package=images -input=images/zombie1_hold.png -output=./images/zombie1.go -var=Zombie1_png

//go:generate file2byteslice -package=sound -input=sound/all.mp3 -output=./sound/all.go -var=All_mp3

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

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	// Settings
	screenWidth  = 800
	screenHeight = 800
	sampleRate   = 44100
	Running      = iota
	Paused
)

var Player player

var bulletSprite *ebiten.Image
var enemySprite *ebiten.Image

var Score int
var Health = 15
var HighScore int

var state = Running

// var audioContext *audio.Context

// var audioPlayer *audio.Player

func init() {
	// // Create audio player
	// var err error
	// audioContext, err = audio.NewContext(sampleRate)
	// if err != nil {
	// 	log.Fatal(err)
	// }

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

	bulletSprite, err = ebiten.NewImage(5, 5, ebiten.FilterNearest)
	if err != nil {
		log.Fatal("Could not create dot image: ", err)
	}
	bulletSprite.Fill(color.White)

	// Create player
	Player = player{
		x:           screenWidth / 2,
		y:           screenHeight / 2,
		height:      float64(heigth) / 2,
		width:       float64(width) / 2,
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

	initBullets()
	initEnemies()

}

var lastEnemy = time.Now()

func update(screen *ebiten.Image) error {
	// if audioPlayer == nil {
	// 	// Decode the wav file.
	// 	// wavS is a decoded io.ReadCloser and io.Seeker.
	// 	wavS, err := mp3.Decode(audioContext, audio.BytesReadSeekCloser(sound.All_mp3))
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// Create an infinite loop stream from the decoded bytes.
	// 	// s is still an io.ReadCloser and io.Seeker.
	// 	s := audio.NewInfiniteLoop(wavS, wavS.Length())

	// 	audioPlayer, err = audio.NewPlayer(audioContext, s)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// Play the infinite-length stream. This never ends.
	// 	audioPlayer.Play()
	// }

	if Health < 0 {
		if Score > HighScore {
			HighScore = Score
		}
		Health = 15
		Score = 0
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) && state == Running {
		state = Paused
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) && state == Paused {
		state = Running
	}

	if state == Paused {
		return nil
	}

	Player.Move()

	// Find mouse position
	x, y := ebiten.CursorPosition()

	// Find angle to mouse cursor
	Player.Angle = math.Atan2(float64(y)-Player.Y(), float64(x)-Player.X())

	// Check if mouse is pressed and shoot
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if time.Now().Sub(Player.LastShoot) > 1*time.Second {
			Player.Shoot()
		}
	}

	for t, b := range bullets {
		if b.deleted == true {
			delete(bullets, t)
		}
		b.Update()
		b.Draw(screen)
	}
	// Check if we are near screen sides
	Player.CheckScreen(screenWidth, screenHeight)

	if ebiten.IsRunningSlowly() {
		return nil
	}

	Player.Draw(screen)
	// Draw bullets
	// space.Step(1.0 / ebiten.FPS)

	if time.Now().Sub(lastEnemy) > 1*time.Second {
		NewEnemy((Player.X()+Player.SX())/2, (Player.Y()+Player.SY())/2)
		lastEnemy = time.Now()
	}

	for t, e := range enemies {
		for tt, ee := range enemies {
			if e.CheckCollision(ee) {
				if t != tt {
					e.Speed = 0.5
				}
			}
		}
		if e.CheckCollision(Player) {
			delete(enemies, t)
			Health--
		}
		for _, b := range bullets {
			if e.CheckCollision(b) {
				delete(enemies, t)
				Score++
			}
		}

		e.Update((Player.X()+Player.SX())/2, (Player.Y()+Player.SY())/2)

		e.Draw(screen)
	}

	// FPS counter
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %f", ebiten.CurrentFPS()))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\nScore: %d", Score))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\nHealth: %d", Health))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\n\nHighscore: %d", HighScore))

	return nil
}

func main() {
	if err := ebiten.Run(update, screenWidth, screenHeight, 1, "Lifewreck"); err != nil {
		panic(err)
	}
}

type Something interface {
	X() float64
	Y() float64
	SY() float64
	SX() float64
}
