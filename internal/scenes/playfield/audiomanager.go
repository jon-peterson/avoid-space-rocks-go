package playfield

import (
	"avoid_the_space_rocks/internal/core"
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"sync"
)

type AudioManager struct {
	soundMap map[string]*rl.Sound
	mapLock  sync.RWMutex
}

func NewAudioManager() *AudioManager {
	return &AudioManager{
		soundMap: make(map[string]*rl.Sound),
	}
}

func (mgr *AudioManager) eventMappings() []eventMapping {
	return []eventMapping{
		{"alien:destroyed", mgr.alienDestroyedHandler},
		{"rock:destroyed", mgr.rockExplosionHandler},
		{"spaceship:fire", mgr.spaceshipFireHandler},
		{"spaceship:enter_hyperspace", mgr.spaceshipEnterHyperspaceHandler},
		{"spaceship:destroyed", mgr.spaceshipExplosionHandler},
	}
}

func (mgr *AudioManager) Register(game *core.Game) error {
	for _, sub := range mgr.eventMappings() {
		if err := game.EventBus.SubscribeAsync(sub.event, sub.handler, false); err != nil {
			rl.TraceLog(rl.LogError, "error subscribing to %s event: %v", sub.event, err)
			return err
		}
	}
	return nil
}

func (mgr *AudioManager) Deregister(game *core.Game) error {
	for _, sub := range mgr.eventMappings() {
		if err := game.EventBus.Unsubscribe(sub.event, sub.handler); err != nil {
			rl.TraceLog(rl.LogError, "error unsubscribing from %s event: %v", sub.event, err)
			return err
		}
	}
	return nil
}

func (mgr *AudioManager) rockExplosionHandler(size core.RockSize) {
	switch size {
	case core.RockTiny:
		_ = mgr.playSound("explosion_tiny.wav")
	case core.RockSmall:
		_ = mgr.playSound("explosion_small.wav")
	case core.RockMedium:
		_ = mgr.playSound("explosion_medium.wav")
	case core.RockBig:
		_ = mgr.playSound("explosion_large.wav")
	}
}

func (mgr *AudioManager) alienDestroyedHandler(_ core.AlienSize) {
	_ = mgr.playSound("explosion_alien.wav")
}

func (mgr *AudioManager) spaceshipFireHandler() {
	_ = mgr.playSound("fire.wav")
}

func (mgr *AudioManager) spaceshipExplosionHandler() {
	_ = mgr.playSound("explosion_ship.wav")
}

func (mgr *AudioManager) spaceshipEnterHyperspaceHandler() {
	_ = mgr.playSound("hyperspace.wav")
}

func (mgr *AudioManager) playSound(filename string) error {
	sound, err := mgr.soundFromFile(filename)
	if err == nil {
		rl.PlaySound(*sound)
	}
	return err
}

// soundFromFile loads a sound from a file, or returns an error if it can't.
// The sound files are cached forever, but they take up very little memory.
func (mgr *AudioManager) soundFromFile(filename string) (*rl.Sound, error) {
	// Almost all the time we have the sound already
	mgr.mapLock.RLock()
	if sound, ok := mgr.soundMap[filename]; ok {
		mgr.mapLock.RUnlock()
		return sound, nil
	}
	mgr.mapLock.RUnlock()
	// Load and save the sound file so we need a writers lock
	mgr.mapLock.Lock()
	defer mgr.mapLock.Unlock()
	sound := rl.LoadSound(fmt.Sprintf("assets/audio/%s", filename))
	if sound.Stream.Buffer == nil {
		return &sound, fmt.Errorf("could not load sound file %s", filename)
	}
	mgr.soundMap[filename] = &sound
	return &sound, nil
}
