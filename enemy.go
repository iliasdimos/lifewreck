package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/jakecoffman/cp"
)

type enemy struct {
	X, Y, Heigth, Width int
	Angle               float64
	Speed               float64
}

func NewEnemy(p *player, space *cp.Space) {
	e := enemy{}

	e.X, e.Y = position()
	e.Speed = 30
	e.Angle = math.Atan2(float64(p.Y-e.Y), float64(p.X-e.X))

	enemyShape := makeEnemy(float64(e.X), float64(e.Y))
	// enemyShape.SetCollisionType(EnemyCollisionType)

	enemyBody := space.AddBody(enemyShape.Body())
	enemyBody.UserData = "enemy"
	enemyBody.SetAngle(e.Angle)
	enemyBody.SetPositionUpdateFunc(func(body *cp.Body, dt float64) {
		v := cp.Vector{
			X: body.Position().X + (math.Cos(body.Angle()) * e.Speed * dt),
			Y: body.Position().Y + (math.Sin(body.Angle()) * e.Speed * dt),
		}
		body.SetPosition(v)
	})
	// space.AddBody(enemyBody)
	// space.AddShape(enemyShape)
}

const EnemyCollisionType cp.CollisionType = 1

func makeEnemy(x, y float64) *cp.Shape {
	body := cp.NewBody(10.0, cp.BODY_KINEMATIC)
	body.SetPosition(cp.Vector{x, y})

	shape := cp.NewCircle(body, 0.95, cp.Vector{})
	shape.SetElasticity(0)
	shape.SetFriction(0)

	shape.SetCollisionType(EnemyCollisionType)

	return shape
}

func (e *enemy) Update(p *player) {
	e.Angle = math.Atan2(float64(p.Y-e.Y), float64(p.X-e.X))

}

func position() (x, y int) {
	rand.Seed(time.Now().UTC().UnixNano())
	side := rand.Intn(3)

	switch side {
	case 0:
		x = 0
		y = rand.Intn(screenHeight)
	case 1:
		x = rand.Intn(screenWidth)
		y = 0
	case 2:
		x = screenWidth
		y = rand.Intn(screenHeight)
	case 3:
		x = rand.Intn(screenWidth)
		y = screenHeight
	}
	return x, y
}
