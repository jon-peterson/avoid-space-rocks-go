package playfield

import (
	"avoid_the_space_rocks/internal/gameobjects"
	"avoid_the_space_rocks/internal/utils"
	rl "github.com/gen2brain/raylib-go/raylib"
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
		X: random.Float32(game.World.width),
		Y: random.Float32(game.World.height),
	}
	rock.rotationSpeed = random.Float32(rockMaxRotate) / 4
	if random.Chance(0.5) {
		rock.rotationSpeed = -rock.rotationSpeed
	}
	// Create a random velocity vector; big rocks are slow
	rock.MaxVelocity = rockMaxSpeed / 4
	rock.Velocity = rl.Vector2{
		X: random.Float32InRange(-rock.MaxVelocity, rock.MaxVelocity),
		Y: random.Float32InRange(-rock.MaxVelocity, rock.MaxVelocity),
	}
	rock.Rotation = rl.Vector2{X: 1, Y: 0}
	return rock
}

func (r *Rock) Update() error {
	game := GetGame()
	delta := rl.GetFrameTime()
	r.Rotation = rl.Vector2Rotate(r.Rotation, r.rotationSpeed*delta)
	r.Rigidbody.ApplyPhysics()
	r.Position = game.World.Wraparound(r.Position)
	return nil
}

func (r *Rock) Draw() error {
	return r.SpriteSheet.Draw(0, 0, r.Position, r.Rotation)
}
