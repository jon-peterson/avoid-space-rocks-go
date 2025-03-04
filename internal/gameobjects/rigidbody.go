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
	delta := rl.GetFrameTime()
	b.Velocity = rl.Vector2Add(b.Velocity, b.Acceleration)
	move := rl.Vector2Scale(b.Velocity, delta)
	if b.MaxVelocity > 0 {
		move = rl.Vector2ClampValue(move, 0, b.MaxVelocity)
	}
	b.Position = rl.Vector2Add(b.Position, move)
}
