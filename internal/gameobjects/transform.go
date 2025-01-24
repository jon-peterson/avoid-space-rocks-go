package gameobjects

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Transform struct {
	Position rl.Vector2
	Rotation rl.Vector2
}

func (t *Transform) String() string {
	return fmt.Sprintf("pos (%f,%f) rot (%f,%f)", t.Position.X, t.Position.Y, t.Rotation.X, t.Rotation.Y)
}
