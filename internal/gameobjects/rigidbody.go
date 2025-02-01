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

func (b *Rigidbody) String() string {
	return fmt.Sprintf("vel (%f,%f)", b.Velocity.X, b.Velocity.Y)
}

func (b *Rigidbody) ApplyPhysics() {
	b.Velocity = rl.Vector2Add(b.Velocity, b.Acceleration)
	b.Velocity = rl.Vector2ClampValue(b.Velocity, 0, b.MaxVelocity)
	b.Position = rl.Vector2Add(b.Position, b.Velocity)
}
