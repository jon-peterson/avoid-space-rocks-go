package gameobjects

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
)

type SpriteSheet struct {
	texture     rl.Texture2D // The texture with the packed sprites
	frameWidth  int32        // Width of each frame in pixels
	frameHeight int32        // Height of each frame pixels
	rows        int32        // the number of rows in the spritesheet
	cols        int32        // the number of columns in the spritesheet
	origin      rl.Vector2   // The middle of the sprite (for rotation)
}

func NewSpriteSheet(file string, rows, cols int32) (SpriteSheet, error) {
	// TODO: share textures by name among all instances
	sheetTexture := rl.LoadTexture("assets/sprites/" + file)

	if sheetTexture.Width%cols != 0 || sheetTexture.Height%cols != 0 {
		return SpriteSheet{},
			fmt.Errorf("spritesheet of dimensions (%d,%d) can't be broken into %d rows and %d cols",
				sheetTexture.Width, sheetTexture.Height, rows, cols)
	}

	s := SpriteSheet{
		texture:     sheetTexture,
		frameWidth:  sheetTexture.Width / cols,
		frameHeight: sheetTexture.Height / rows,
		rows:        rows,
		cols:        cols,
	}
	s.origin = rl.NewVector2(float32(s.frameWidth)/2, float32(s.frameHeight)/2)
	return s, nil
}

func (s *SpriteSheet) String() string {
	return fmt.Sprintf("%v (%dx%d)", s.texture, s.frameWidth, s.frameHeight)
}

// Draw the sprite at the given frame at the given location and rotation
func (s *SpriteSheet) Draw(frameRow, frameCol int, loc, rot rl.Vector2) error {
	frame, err := s.frame(frameRow, frameCol)
	if err != nil {
		return err
	}
	destination := rl.Rectangle{
		X:      loc.X,
		Y:      loc.Y,
		Width:  float32(s.frameWidth),
		Height: float32(s.frameHeight),
	}
	rotationDegrees := float32(math.Atan2(float64(rot.Y), float64(rot.X)) * 180 / math.Pi)
	rl.DrawTexturePro(s.texture, frame, destination, s.origin, rotationDegrees, rl.Black)
	return nil
}

// frame returns the rectangle for the given frame in the spritesheet
func (s *SpriteSheet) frame(row, col int) (rl.Rectangle, error) {
	if row < 0 || row >= int(s.rows) || col < 0 || col >= int(s.cols) {
		return rl.Rectangle{}, fmt.Errorf("frame (%d,%d) is out of bounds", row, col)
	}
	return rl.Rectangle{
		X:      float32(col) * float32(s.frameWidth),
		Y:      float32(row) * float32(s.frameHeight),
		Width:  float32(s.frameWidth),
		Height: float32(s.frameHeight),
	}, nil
}
