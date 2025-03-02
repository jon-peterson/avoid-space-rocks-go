package playfield

import (
	"os"
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func TestMain(m *testing.M) {
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
