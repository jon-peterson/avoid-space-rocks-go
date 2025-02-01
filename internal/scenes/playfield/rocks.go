package playfield

import (
	"avoid_the_space_rocks/internal/gameobjects"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math/rand"
)

type Rock struct {
	gameobjects.Rigidbody
	gameobjects.SpriteSheet
	rotationSpeed float32 // rotations per second
}

func NewRockBig() Rock {
	sheet, _ := gameobjects.NewSpriteSheet("rock_big.png", 1, 1)
	rock := Rock{
		SpriteSheet: sheet,
	}
	game := GetGame()
	rock.Position = rl.Vector2{
		X: rand.Float32() * game.World.width,
		Y: rand.Float32() * game.World.height,
	}
	rock.rotationSpeed = rand.Float32() * rockMaxRotate
	if rand.Intn(2) == 0 {
		rock.rotationSpeed = -rock.rotationSpeed
	}
	// Create a random velocity vector
	rock.MaxVelocity = rockMaxSpeed
	rock.Velocity = rl.Vector2{
		X: (rand.Float32() * rockMaxSpeed * 2) - rand.Float32(),
		Y: (rand.Float32() * rockMaxSpeed * 2) - rand.Float32(),
	}
	return rock
}

func (r Rock) Update() error {
	game := GetGame()
	delta := rl.GetFrameTime()
	r.Rotation = rl.Vector2Rotate(r.Rotation, r.rotationSpeed*delta)
	r.Rigidbody.ApplyPhysics()
	game.World.Wraparound(r.Position)
	return nil
}

func (r Rock) Draw() error {
	return r.SpriteSheet.Draw(0, 0, r.Position, r.Rotation)
}
