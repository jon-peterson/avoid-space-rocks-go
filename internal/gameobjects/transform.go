package gameobjects

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Transform struct {
	position rl.Vector2
	rotation rl.Vector2
}

func (t *Transform) String() string {
	return fmt.Sprintf("pos (%f,%f) rot (%f,%f)", t.position.X, t.position.Y, t.rotation.X, t.rotation.Y)
}
