package gameobjects

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"testing"
)

func TestRigidbody_ApplyPhysics(t *testing.T) {
	tests := []struct {
		name    string
		rb      Rigidbody
		delta   float32
		wantPos rl.Vector2
		wantVel rl.Vector2
	}{
		{
			name: "constant velocity no acceleration",
			rb: Rigidbody{
				Transform:    Transform{Position: rl.Vector2{X: 0, Y: 0}},
				Velocity:     rl.Vector2{X: 1, Y: 2},
				Acceleration: rl.Vector2{X: 0, Y: 0},
				MaxVelocity:  0,
			},
			delta:   1.0,
			wantPos: rl.Vector2{X: 1, Y: 2},
			wantVel: rl.Vector2{X: 1, Y: 2},
		},
		{
			name: "constant velocity no acceleration small delta",
			rb: Rigidbody{
				Transform:    Transform{Position: rl.Vector2{X: 0, Y: 0}},
				Velocity:     rl.Vector2{X: 1, Y: 0},
				Acceleration: rl.Vector2{X: 0, Y: 0},
				MaxVelocity:  0,
			},
			delta:   0.3,
			wantPos: rl.Vector2{X: 0.3, Y: 0},
			wantVel: rl.Vector2{X: 1, Y: 0},
		},
		{
			name: "with acceleration",
			rb: Rigidbody{
				Transform:    Transform{Position: rl.Vector2{X: 0, Y: 0}},
				Velocity:     rl.Vector2{X: 1, Y: 0},
				Acceleration: rl.Vector2{X: 2, Y: 0},
				MaxVelocity:  0,
			},
			delta:   1.0,
			wantPos: rl.Vector2{X: 3, Y: 0},
			wantVel: rl.Vector2{X: 3, Y: 0},
		},
		{
			name: "respect max velocity",
			rb: Rigidbody{
				Transform:    Transform{Position: rl.Vector2{X: 0, Y: 0}},
				Velocity:     rl.Vector2{X: 2, Y: 0},
				Acceleration: rl.Vector2{X: 5, Y: 0},
				MaxVelocity:  3,
			},
			delta:   1.0,
			wantPos: rl.Vector2{X: 3, Y: 0},
			wantVel: rl.Vector2{X: 3, Y: 0},
		},
		{
			name: "respect max velocity small delta",
			rb: Rigidbody{
				Transform:    Transform{Position: rl.Vector2{X: 0, Y: 0}},
				Velocity:     rl.Vector2{X: 3, Y: 0},
				Acceleration: rl.Vector2{X: 5, Y: 0},
				MaxVelocity:  1,
			},
			delta:   0.5,
			wantPos: rl.Vector2{X: 0.5, Y: 0}, // Position reflects half of velocity
			wantVel: rl.Vector2{X: 1.0, Y: 0}, // Velocity capped
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := tt.rb
			rb.ApplyPhysics(tt.delta)

			if rb.Position != tt.wantPos {
				t.Errorf("Position = %v, want %v", rb.Position, tt.wantPos)
			}
			if rb.Velocity != tt.wantVel {
				t.Errorf("Velocity = %v, want %v", rb.Velocity, tt.wantVel)
			}
		})
	}
}
