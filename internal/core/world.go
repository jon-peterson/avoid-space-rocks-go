package core

import (
	"avoid_the_space_rocks/internal/gameobjects"
	random "avoid_the_space_rocks/internal/utils"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// The World object represents the state of the game within the playfield
type World struct {
	Width     float32 // Width of the playfield in worldspace
	Height    float32 // Height of the playfield in worldspace
	Spaceship Spaceship
	Objects   gameobjects.GameObjectCollection
}

func NewWorld(width, height float32) *World {
	w := World{
		Width:  width,
		Height: height,
	}
	return &w
}

func (w *World) Initialize() {
	// Spaceship starts in the middle pointing up
	w.Objects = gameobjects.NewGameObjectCollection()
	w.Spaceship = NewSpaceship()
	go w.Spaceship.Spawn()
}

// Wraparound returns the position of the given position, wrapping around the edges of the playfield
func (w *World) Wraparound(p rl.Vector2) rl.Vector2 {
	if p.X < 0 {
		p.X = w.Width
	} else if p.X > w.Width {
		p.X = 0
	}
	if p.Y < 0 {
		p.Y = w.Height
	} else if p.Y > w.Height {
		p.Y = 0
	}
	return p
}

// IsOutsideEdges returns true if the given position is outside the edges of the playfield
func (w *World) IsOutsideEdges(p rl.Vector2) bool {
	return p.X < 0 || p.X > w.Width || p.Y < 0 || p.Y > w.Height
}

// RandomBorderPosition returns a random position on the border of the playfield, each
// equally likely.
func (w *World) RandomBorderPosition() rl.Vector2 {
	if random.Chance(0.5) {
		return rl.Vector2{
			X: random.RndFloat32(w.Width),
			Y: random.Choice([]float32{0, w.Height}),
		}
	}
	return rl.Vector2{
		X: random.Choice([]float32{0, w.Width}),
		Y: random.RndFloat32(w.Height),
	}
}

// RandomPosition returns a random position within the playfield.
func (w *World) RandomPosition() rl.Vector2 {
	return rl.Vector2{
		X: random.RndFloat32(w.Width),
		Y: random.RndFloat32(w.Height),
	}
}
