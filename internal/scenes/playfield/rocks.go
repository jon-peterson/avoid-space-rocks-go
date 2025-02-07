package playfield

import (
	"avoid_the_space_rocks/internal/gameobjects"
	"avoid_the_space_rocks/internal/utils"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Rock is a game object that has a consistent rotation speed and constant velocity.
type Rock struct {
	gameobjects.Rigidbody
	gameobjects.SpriteSheet
	rotationSpeed float32 // rotations per second
	isAlive       bool
}

// NewRockBig creates a new large rock with a random position and velocity, spinning randomly.
func NewRockBig() Rock {
	game := GetGame()
	sheet, _ := gameobjects.NewSpriteSheet("rock_big.png", 1, 1)
	rock := Rock{
		SpriteSheet: sheet,
		Rigidbody: gameobjects.Rigidbody{
			Transform: gameobjects.Transform{
				Position: game.World.RandomBorderLocation(),
				Rotation: rl.Vector2{X: 1, Y: 0},
			},
		},
		rotationSpeed: random.Float32(rockMaxRotate) / 4,
		isAlive:       true,
	}
	// Half of 'em rotate counterclockwise
	if random.Chance(0.5) {
		rock.rotationSpeed = -rock.rotationSpeed
	}
	// Randomize the speed and direction
	rock.MaxVelocity = rockMaxSpeed / 4
	rock.Velocity = rl.Vector2{
		X: random.Float32InRange(-rock.MaxVelocity, rock.MaxVelocity),
		Y: random.Float32InRange(-rock.MaxVelocity, rock.MaxVelocity),
	}
	return rock
}

// Update applies physics to the rock so it moves per its velocity and rotation speed.
func (r *Rock) Update() error {
	game := GetGame()
	delta := rl.GetFrameTime()
	r.Rotation = rl.Vector2Rotate(r.Rotation, r.rotationSpeed*delta)
	r.Rigidbody.ApplyPhysics()
	r.Position = game.World.Wraparound(r.Position)
	return nil
}

// Draw renders the rock to the screen.
func (r *Rock) Draw() error {
	return r.SpriteSheet.Draw(0, 0, r.Position, r.Rotation)
}

func (r *Rock) IsAlive() bool {
	return r.isAlive
}
