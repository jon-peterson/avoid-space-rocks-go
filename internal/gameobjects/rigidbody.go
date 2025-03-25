package gameobjects

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Rigidbody struct {
	Transform
	Acceleration rl.Vector2
	Velocity     rl.Vector2
	MaxVelocity  float32 // The maximum magnitude of the velocity vector
}

func (rb *Rigidbody) String() string {
	return fmt.Sprintf("vel (%f,%f)", rb.Velocity.X, rb.Velocity.Y)
}

// ApplyPhysics applies acceleration to the velocity and then moves the object
func (rb *Rigidbody) ApplyPhysics() {
	delta := rl.GetFrameTime()
	rb.Velocity = rl.Vector2Add(rb.Velocity, rb.Acceleration)
	if rb.MaxVelocity > 0 {
		// TODO: fix this; it's not applying max velocity correctly
		rb.Velocity = rl.Vector2ClampValue(rb.Velocity, 0, rb.MaxVelocity)
	}
	move := rl.Vector2Scale(rb.Velocity, delta)
	rb.Position = rl.Vector2Add(rb.Position, move)
}
