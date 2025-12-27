package internal

import (
	"os"
	"strconv"

	"github.com/hegedustibor/htgo-tts/voices"
)

type Config struct {
	Language    string
	MessageText string
	Volume      float64
	AudioFolder string
	Port        string
}

func LoadConfig() *Config {
	volumePercent := getEnvInt("VOLUME", 100)
	volumeFloat := float64(volumePercent) / 100.0
	if volumeFloat > 1.0 {
		volumeFloat = 1.0
	} else if volumeFloat < 0 {
		volumeFloat = 0
	}

	return &Config{
		Language:    getEnv("LANGUAGE", voices.English),
		MessageText: getEnv("MESSAGE_TEXT", "Time to sleep"),
		Volume:      volumeFloat,
		AudioFolder: getEnv("AUDIO_FOLDER", "audio"),
		Port:        getEnv("PORT", ":8080"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	valueStr, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return fallback
	}
	return value
}
