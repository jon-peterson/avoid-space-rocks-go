package playfield

import (
	"avoid_the_space_rocks/internal/core"
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/hashicorp/go-set"
	"sync"
)

type AudioManager struct {
	soundMap  map[string]*rl.Sound
	soundLock sync.RWMutex

	musicMap     map[string]*rl.Music
	musicLock    sync.RWMutex
	playingMusic set.Set[string]
}

var _ core.EventObserver = (*AudioManager)(nil)

func NewAudioManager() *AudioManager {
	return &AudioManager{
		soundMap:     make(map[string]*rl.Sound),
		musicMap:     make(map[string]*rl.Music),
		playingMusic: *set.New[string](10),
	}
}

func (mgr *AudioManager) eventMappings() []eventMapping {
	return []eventMapping{
		{"alien:destroyed", mgr.alienDestroyedHandler},
		{"alien:fire", mgr.alienFireHandler},
		{"rock:destroyed", mgr.rockExplosionHandler},
		{"spaceship:fire", mgr.spaceshipFireHandler},
		{"spaceship:thrust", mgr.spaceshipThrustHandler},
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
	mgr.playingMusic.ForEach(func(filename string) bool {
		rl.StopMusicStream(*mgr.musicMap[filename])
		return true
	})
	return nil
}

// Update is called every frame to update the audio manager. If there are any playing music
// streams, it updates them so they continue to play.
func (mgr *AudioManager) Update(game *core.Game) error {
	if !mgr.playingMusic.Empty() {
		mgr.musicLock.RLock()
		defer mgr.musicLock.RUnlock()
		mgr.playingMusic.ForEach(func(filename string) bool {
			rl.UpdateMusicStream(*mgr.musicMap[filename])
			return true
		})
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

func (mgr *AudioManager) alienFireHandler() {
	_ = mgr.playSound("fire_alien.wav")
}

func (mgr *AudioManager) spaceshipFireHandler() {
	_ = mgr.playSound("fire.wav")
}

func (mgr *AudioManager) spaceshipThrustHandler(start bool) {
	if start {
		mgr.startMusic("fuel_burn.wav")
	} else {
		mgr.stopMusic("fuel_burn.wav")
	}
}

func (mgr *AudioManager) spaceshipExplosionHandler() {
	if mgr.playingMusic.Contains("fuel_burn.wav") {
		mgr.stopMusic("fuel_burn.wav")
	}
	_ = mgr.playSound("explosion_ship.wav")
}

func (mgr *AudioManager) spaceshipEnterHyperspaceHandler() {
	if mgr.playingMusic.Contains("fuel_burn.wav") {
		mgr.stopMusic("fuel_burn.wav")
	}
	_ = mgr.playSound("hyperspace.wav")
}

// playSound plays a sound from a filename, or returns an error if it can't.
func (mgr *AudioManager) playSound(filename string) error {
	sound, err := mgr.soundFromFile(filename)
	if err == nil {
		rl.PlaySound(*sound)
	}
	return err
}

// startMusic starts playing a music file from a filename, or returns an error if it can't.
func (mgr *AudioManager) startMusic(filename string) {
	if !mgr.playingMusic.Contains(filename) {
		rl.TraceLog(rl.LogDebug, "Starting music for %s", filename)
		mgr.withMusic(filename, func(music *rl.Music) {
			mgr.playingMusic.Insert(filename)
			rl.PlayMusicStream(*music)
		})
	}
}

// stopMusic stops playing a music file from a filename, or returns an error if it can't.
func (mgr *AudioManager) stopMusic(filename string) {
	if mgr.playingMusic.Contains(filename) {
		rl.TraceLog(rl.LogDebug, "Stopping music for %s", filename)
		mgr.withMusic(filename, func(music *rl.Music) {
			mgr.playingMusic.Remove(filename)
			rl.StopMusicStream(*music)
		})
	}
}

type musicHandler func(*rl.Music)

func (mgr *AudioManager) withMusic(filename string, callback musicHandler) {
	music, err := mgr.musicFromFile(filename)
	if err == nil {
		callback(music)
	}
}

// soundFromFile loads a sound from a file, or returns an error if it can't.
// The sound files are cached forever, but they take up very little memory.
func (mgr *AudioManager) soundFromFile(filename string) (*rl.Sound, error) {
	// Almost all the time we have the sound already
	mgr.soundLock.RLock()
	if sound, ok := mgr.soundMap[filename]; ok {
		mgr.soundLock.RUnlock()
		return sound, nil
	}
	mgr.soundLock.RUnlock()
	// Load and save the sound file so we need a writers lock
	mgr.soundLock.Lock()
	defer mgr.soundLock.Unlock()
	sound := rl.LoadSound(fmt.Sprintf("assets/audio/%s", filename))
	if sound.Stream.Buffer == nil {
		return &sound, fmt.Errorf("could not load sound file %s", filename)
	}
	mgr.soundMap[filename] = &sound
	return &sound, nil
}

// musicFromFile loads music from a file, or returns an error if it can't.
// The music files are cached forever, but they take up very little memory.
func (mgr *AudioManager) musicFromFile(filename string) (*rl.Music, error) {
	// Almost all the time we have the sound already
	mgr.musicLock.RLock()
	if music, ok := mgr.musicMap[filename]; ok {
		mgr.musicLock.RUnlock()
		return music, nil
	}
	mgr.musicLock.RUnlock()
	// Load and save the music file so we need a writers lock
	mgr.musicLock.Lock()
	defer mgr.musicLock.Unlock()
	music := rl.LoadMusicStream(fmt.Sprintf("assets/audio/%s", filename))
	if music.Stream.Buffer == nil {
		return &music, fmt.Errorf("could not load music file %s", filename)
	}
	music.Looping = true
	mgr.musicMap[filename] = &music
	return &music, nil
}
