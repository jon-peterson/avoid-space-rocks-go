package playfield

import (
	"avoid_the_space_rocks/internal/gameobjects"
	random "avoid_the_space_rocks/internal/utils"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// The World object represents the state of the game within the playfield
type World struct {
	width     float32 // Width of the playfield in worldspace
	height    float32 // Height of the playfield in worldspace
	Spaceship Spaceship
	Objects   gameobjects.GameObjectCollection
}

func NewWorld(width, height int32) World {
	w := World{
		width:  float32(width),
		height: float32(height),
	}
	return w
}

func (w *World) InitializeLevel(level int) {
	w.Objects = gameobjects.NewGameObjectCollection()
	// Spaceship starts in the middle pointing up
	w.Spaceship = NewSpaceship()
	w.Spaceship.Position = rl.Vector2{
		X: w.width / 2,
		Y: w.height / 2,
	}
	// Random rocks based on the level number
	for i := 0; i < 4; i++ {
		rock := NewRockBig()
		w.Objects.Add(&rock)
	}
}

// Wraparound returns the position of the given position, wrapping around the edges of the playfield
func (w *World) Wraparound(p rl.Vector2) rl.Vector2 {
	if p.X < 0 {
		p.X = w.width
	} else if p.X > w.width {
		p.X = 0
	}
	if p.Y < 0 {
		p.Y = w.height
	} else if p.Y > w.height {
		p.Y = 0
	}
	return p
}

func (w *World) RandomBorderLocation() rl.Vector2 {
	if random.Chance(0.5) {
		return rl.Vector2{
			X: random.RndFloat32(w.width),
			Y: random.Choice([]float32{0, w.height}),
		}
	}
	return rl.Vector2{
		X: random.Choice([]float32{0, w.width}),
		Y: random.RndFloat32(w.height),
	}

}
