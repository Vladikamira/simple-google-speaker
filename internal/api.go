package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type SpeakRequest struct {
	Message  string `json:"message"`
	Language string `json:"language"`
}

type APIHandler struct {
	cfg     *Config
	tts     *TTSService
	castSvc *CastService
}

func NewAPIHandler(cfg *Config, tts *TTSService, castSvc *CastService) *APIHandler {
	return &APIHandler{
		cfg:     cfg,
		tts:     tts,
		castSvc: castSvc,
	}
}

func (h *APIHandler) SpeakHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SpeakRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("Invalid API request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	slog.Info("Received speak request", "message", req.Message, "language", req.Language)

	// Use defaults from config if not provided
	message := req.Message
	if message == "" {
		message = h.cfg.MessageText
	}

	language := req.Language
	if language == "" {
		language = h.cfg.Language
	}

	volume := h.cfg.Volume

	// 1. Find device
	client := h.castSvc.FindDevice()
	if client == nil {
		http.Error(w, "No Google Cast devices found", http.StatusServiceUnavailable)
		return
	}

	// 2. Generate audio (we need a way to pass language to TTS if it's different)
	// For now, TTS service is initialized with a language.
	// To support dynamic language, we can update GenerateAudio to accept it.
	fileName, err := h.tts.GenerateAudio(message, language)
	if err != nil {
		slog.Error("TTS generation failed", "error", err)
		http.Error(w, "Failed to generate audio", http.StatusInternalServerError)
		return
	}

	localIP := GetLocalIP()
	mediaURL := fmt.Sprintf("http://%s%s/%s", localIP, h.cfg.Port, fileName)

	// 3. Play audio
	go func() {
		slog.Info("Playing from API", "url", mediaURL, "volume", volume)
		if err := h.castSvc.PlayMedia(context.Background(), client, mediaURL, volume); err != nil {
			slog.Error("API Playback failed", "error", err)
		}
	}()

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"status": "playing", "url": mediaURL})
}
