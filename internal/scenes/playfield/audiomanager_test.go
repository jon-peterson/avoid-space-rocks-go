package playfield

import (
	"os"
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func TestMain(m *testing.M) {
	// Change to the project root directory
	if err := os.Chdir("../../.."); err != nil {
		panic("failed to change to project root directory: " + err.Error())
	}

	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()
	os.Exit(m.Run())
}

func TestSoundFromFile(t *testing.T) {
	audioManager := NewAudioManager()

	// Test loading a sound file
	filename := "fire.wav"
	sound, err := audioManager.soundFromFile(filename)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if sound.Stream.Buffer == nil {
		t.Fatalf("expected valid sound, got nil buffer")
	}

	// Test loading the same sound file from cache
	cachedSound, err := audioManager.soundFromFile(filename)
	if err != nil {
		t.Fatalf("expected no error loading cache, got %v", err)
	}
	if cachedSound.Stream.Buffer == nil {
		t.Fatalf("expected valid sound loading cache, got nil buffer")
	}
	if sound != cachedSound {
		t.Fatalf("expected cached sound, got different instance")
	}
}

func TestSoundFromFile_NotExist(t *testing.T) {
	audioManager := NewAudioManager()

	// Test loading a non-existent sound file
	filename := "non_existent_file.wav"
	_, err := audioManager.soundFromFile(filename)
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
}

func TestMusicFromFile(t *testing.T) {
	audioManager := NewAudioManager()

	// Test loading a music file
	filename := "fuel_burn.wav"
	music, err := audioManager.musicFromFile(filename)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if music.Stream.Buffer == nil {
		t.Fatalf("expected valid music, got nil buffer")
	}

	// Test loading the same music file from cache
	cachedMusic, err := audioManager.musicFromFile(filename)
	if err != nil {
		t.Fatalf("expected no error loading cache, got %v", err)
	}
	if cachedMusic.Stream.Buffer == nil {
		t.Fatalf("expected valid music loading cache, got nil buffer")
	}
	if music != cachedMusic {
		t.Fatalf("expected cached music, got different instance")
	}
}
