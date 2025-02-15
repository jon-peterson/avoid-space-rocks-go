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
	if err := game.EventBus.Subscribe("spaceship:fire", SpaceshipFireHandler); err != nil {
		rl.TraceLog(rl.LogError, "error subscribing to spaceship:fire event: %v", err)
		return err
	}
	return nil
}

func DeregisterAudioManager(game *core.Game) error {
	if err := game.EventBus.Unsubscribe("rock:destroyed", RockExplosionHandler); err != nil {
		rl.TraceLog(rl.LogError, "error unsubscribing from rock:destroyed event: %v", err)
		return err
	}
	if err := game.EventBus.Unsubscribe("spaceship:fire", SpaceshipFireHandler); err != nil {
		rl.TraceLog(rl.LogError, "error unsubscribing from spaceship:fire event: %v", err)
		return err
	}
	return nil
}

func RockExplosionHandler(size core.RockSize) {
	switch size {
	case core.RockTiny:
		_ = playSound("explosion_tiny.wav")
	case core.RockSmall:
		_ = playSound("explosion_small.wav")
	case core.RockMedium:
		_ = playSound("explosion_medium.wav")
	case core.RockBig:
		_ = playSound("explosion_large.wav")
	}
}

func SpaceshipFireHandler() {
	_ = playSound("fire.wav")
}

func playSound(filename string) error {
	sound, err := soundFromFile(filename)
	if err == nil {
		rl.PlaySound(*sound)
	}
	return err
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
