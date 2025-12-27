package internal

import (
	"crypto/sha1"
	"encoding/hex"
	"log/slog"
	"os"
	"path/filepath"

	htgotts "github.com/hegedustibor/htgo-tts"
)

type TTSService struct {
	folder   string
	language string
}

func NewTTSService(folder, language string) *TTSService {
	return &TTSService{
		folder:   folder,
		language: language,
	}
}

func (s *TTSService) GenerateAudio(text string, language string) (string, error) {
	hash := sha1.Sum([]byte(text + language)) // Include language in hash
	currentHash := hex.EncodeToString(hash[:])

	baseName := "message"
	fileName := baseName + ".mp3"
	filePath := filepath.Join(s.folder, fileName)
	hashPath := filepath.Join(s.folder, baseName+".sha1")

	storedHash, err := os.ReadFile(hashPath)
	if err == nil && string(storedHash) == currentHash {
		if _, err := os.Stat(filePath); err == nil {
			slog.Info("Audio content unchanged. Using existing file.")
			return fileName, nil
		}
	}

	slog.Info("Generating new audio file...")
	_ = os.Remove(filePath)

	speech := htgotts.Speech{Folder: s.folder, Language: language}
	_, err = speech.CreateSpeechFile(text, baseName)
	if err != nil {
		return "", err
	}

	_ = os.WriteFile(hashPath, []byte(currentHash), 0644)
	return fileName, nil
}
