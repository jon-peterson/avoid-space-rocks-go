package gameobjects

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
	"sync"
)

type SpriteSheet struct {
	name        string
	texture     rl.Texture2D // The texture with the packed sprites
	frameWidth  int32        // Width of each frame in pixels
	frameHeight int32        // Height of each frame pixels
	rows        int32        // the number of rows in the spritesheet
	cols        int32        // the number of columns in the spritesheet
	origin      rl.Vector2   // The middle of the sprite (for rotation)
}

type SpriteManager struct {
	spritesMap map[string]*SpriteSheet
	mapLock    sync.RWMutex
}

var spriteManager *SpriteManager = newSpriteManager()

func newSpriteManager() *SpriteManager {
	return &SpriteManager{
		spritesMap: make(map[string]*SpriteSheet),
	}
}

// LoadSpriteSheet loads a spritesheet from the given file, or returns an error if it can't.
// It initializes it with the specified rows and columns. SpriteSheets are cached.
func LoadSpriteSheet(file string, rows, cols int32) (*SpriteSheet, error) {
	// Return the existing sprite sheet if it's already been loaded
	spriteManager.mapLock.RLock()
	if sprite, ok := spriteManager.spritesMap[file]; ok {
		spriteManager.mapLock.RUnlock()
		return sprite, nil
	}
	spriteManager.mapLock.RUnlock()
	// Create, store, and return the spritesheet at that file
	spriteManager.mapLock.Lock()
	defer spriteManager.mapLock.Unlock()
	sheetTexture := rl.LoadTexture("assets/sprites/" + file)
	if sheetTexture.Width%cols != 0 || sheetTexture.Height%cols != 0 {
		return &SpriteSheet{},
			fmt.Errorf("spritesheet of dimensions (%d,%d) can't be broken into %d rows and %d cols",
				sheetTexture.Width, sheetTexture.Height, rows, cols)
	}

	s := SpriteSheet{
		name:        file,
		texture:     sheetTexture,
		frameWidth:  sheetTexture.Width / cols,
		frameHeight: sheetTexture.Height / rows,
		rows:        rows,
		cols:        cols,
	}
	s.origin = rl.NewVector2(float32(s.frameWidth)/2, float32(s.frameHeight)/2)
	spriteManager.spritesMap[file] = &s
	return &s, nil
}

func (s *SpriteSheet) String() string {
	return fmt.Sprintf("%s (%dx%d)", s.name, s.frameWidth, s.frameHeight)
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

// GetSize returns the size of the sprite in pixels as a vector
func (s *SpriteSheet) GetSize() rl.Vector2 {
	return rl.Vector2{
		X: float32(s.frameWidth),
		Y: float32(s.frameHeight),
	}
}

// GetRectangle returns the bounding rectangle where this sprite will be drawn centered at center
func (s *SpriteSheet) GetRectangle(center rl.Vector2) rl.Rectangle {
	return rl.Rectangle{
		X:      center.X - float32(s.frameWidth)/2,
		Y:      center.Y - float32(s.frameHeight)/2,
		Width:  float32(s.frameWidth),
		Height: float32(s.frameHeight),
	}
}
