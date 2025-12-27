package internal

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/barnybug/go-cast"
	"github.com/barnybug/go-cast/controllers"
	"github.com/hashicorp/mdns"
)

type CastService struct{}

func NewCastService() *CastService {
	return &CastService{}
}

func (s *CastService) FindDevice() *cast.Client {
	slog.Info("Searching for Google Cast devices...")
	entriesCh := make(chan *mdns.ServiceEntry, 10)

	go func() {
		params := &mdns.QueryParam{
			Service:     "_googlecast._tcp",
			Domain:      "local",
			Timeout:     time.Second * 5,
			Entries:     entriesCh,
			DisableIPv6: true,
		}
		_ = mdns.Query(params)
		close(entriesCh)
	}()

	for entry := range entriesCh {
		if !strings.Contains(entry.Name, "_googlecast") {
			continue
		}
		slog.Info("Found device", "name", entry.Name, "address", entry.AddrV4)
		client := cast.NewClient(entry.AddrV4, entry.Port)
		return client
	}
	return nil
}

func (s *CastService) PlayMedia(ctx context.Context, client *cast.Client, url string, volume float64) error {
	if err := client.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	_, _ = client.Receiver().SetVolume(ctx, &controllers.Volume{Level: &volume})

	media, err := client.Media(ctx)
	if err != nil {
		return fmt.Errorf("failed to get media controller: %w", err)
	}

	item := controllers.MediaItem{
		ContentId:   url,
		StreamType:  "BUFFERED",
		ContentType: "audio/mpeg",
	}

	if _, err := media.LoadMedia(ctx, item, 0, true, nil); err != nil {
		return fmt.Errorf("failed to load media: %w", err)
	}

	slog.Info("Playback started")
	s.waitForCompletion(ctx, media)
	return nil
}

func (s *CastService) waitForCompletion(ctx context.Context, media *controllers.MediaController) {
	for i := 0; i < 30; i++ {
		status, err := media.GetStatus(ctx)
		if err == nil && len(status.Status) > 0 {
			if status.Status[0].PlayerState == "IDLE" && i > 2 {
				break
			}
		}
		time.Sleep(1 * time.Second)
	}
}
