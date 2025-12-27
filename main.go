package main

import (
	"log/slog"
	"net/http"

	"github.com/Vladikamira/simple-google-speaker/internal"
)

func main() {
	cfg := internal.LoadConfig()

	tts := internal.NewTTSService(cfg.AudioFolder, cfg.Language)
	castSvc := internal.NewCastService()
	api := internal.NewAPIHandler(cfg, tts, castSvc)

	// Route for speaking
	http.HandleFunc("/speak", api.SpeakHandler)

	// Serve audio files
	http.Handle("/", http.FileServer(http.Dir(cfg.AudioFolder)))

	slog.Info("Server starting", "port", cfg.Port)
	if err := http.ListenAndServe(cfg.Port, nil); err != nil {
		slog.Error("Server failed", "error", err)
	}
}
