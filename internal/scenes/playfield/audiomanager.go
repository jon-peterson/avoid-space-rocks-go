package playfield

import (
	"avoid_the_space_rocks/internal/core"
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"sync"
)

var (
	soundMap = make(map[string]rl.Sound)
	mapLock  = sync.RWMutex{}
)

func RegisterAudioManager(game *core.Game) error {
	if err := game.EventBus.Subscribe("rock:destroyed", RockExplosionHandler); err != nil {
		rl.TraceLog(rl.LogError, "error subscribing to rock:destroyed event: %v", err)
		return err
	}
	return nil
}

func DeregisterAudioManager(game *core.Game) error {
	if err := game.EventBus.Unsubscribe("rock:destroyed", RockExplosionHandler); err != nil {
		rl.TraceLog(rl.LogError, "error unsubscribing from rock:destroyed event: %v", err)
		return err
	}
	return nil
}

func RockExplosionHandler(size core.RockSize) {
	switch size {
	case core.RockTiny:
		core.GetGame().Score += 100
	case core.RockSmall:
		core.GetGame().Score += 75
	case core.RockMedium:
		core.GetGame().Score += 50
	case core.RockBig:
		core.GetGame().Score += 20
	}
}

// soundFromFile loads a sound from a file, or returns an error if it can't.
// The sound files are cached forever, but they take up very little memory.
func soundFromFile(filename string) (*rl.Sound, error) {
	// Almost all the time we have the sound already
	mapLock.RLock()
	if sound, ok := soundMap[filename]; ok {
		mapLock.RUnlock()
		return &sound, nil
	}
	mapLock.RUnlock()
	// Load and save the sound file so we need a writers lock
	mapLock.Lock()
	defer mapLock.Unlock()
	sound := rl.LoadSound(fmt.Sprintf("assets/audio/%s", filename))
	if sound.Stream.Buffer == nil {
		return &sound, fmt.Errorf("could not load sound file %s", filename)
	}
	soundMap[filename] = sound
	return &sound, nil
}
