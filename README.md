# Simple Google Speaker

A simple Go-based service that allows you to send text-to-speech (TTS) messages to Google Cast speakers via a REST API.

## Features

- **REST API**: Send messages using a simple POST request.
- **Google TTS**: High-quality speech generation using Google's TTS engine.
- **Caching**: Audio files are cached based on content and language to minimize network requests.
- **Auto-Discovery**: Automatically finds Google Cast devices in your local network using mDNS.
- **Configurable**: Manage default settings via environment variables.

## Getting Started

### Prerequisites

- A Google Cast compatible device (Google Home, Nest Mini, etc.) on the same local network.
- Docker installed (recommended) or Go 1.23+.

### Running with Docker (Recommended)

To ensure the service can discover devices on your network, it's best to use the host network mode.

#### Using Docker CLI
```bash
# Build the image
docker build -t simple-google-speaker .

# Run the container
docker run -d \
  --name simple-google-speaker \
  --network host \
  -e MESSAGE_TEXT="Good morning" \
  -e VOLUME=70 \
  simple-google-speaker
```

#### Using Docker Compose
Create a `docker-compose.yml` file:
```yaml
services:
  speaker:
    image: vladikamira/simple-google-speaker:latest
    container_name: simple-google-speaker
    network_mode: host
    restart: unless-stopped
    environment:
      - PORT=:8080
      - VOLUME=70
    volumes:
      - ./audio:/app/audio
```

Then run:
```bash
docker-compose up -d
```

> **Note**: `network_mode: host` is required for mDNS discovery to work correctly from inside the container.

### Running with Go

```bash
go run main.go
```

## Configuration

The service can be configured using the following environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | The port the HTTP server listens on | `:8080` |
| `VOLUME` | Speaker volume in percent (0-100) | `100` |
| `AUDIO_FOLDER` | Directory to store generated audio files | `audio` |

## API Usage

### Speak a Message

**Endpoint:** `POST /speak`

**Request Body:**
```json
{
  "message": "Привет! Это проверочное сообщение.",
  "language": "ru"
}
```

**Example via cURL:**
```bash
curl -X POST http://localhost:8080/speak \
-H "Content-Type: application/json" \
-d '{
  "message": "Hello from Google Speaker Service",
  "language": "en"
}'
```

The service will:
1. Generate an MP3 file using Google TTS.
2. Search for a Google Cast device in the local network.
3. Serve the MP3 file via its built-in server.
4. Instruct the speaker to play the audio.

## Important Nuances

1. **Network Connectivity**: The Google Cast device must be able to reach the IP address of the machine running this service. Ensure your firewall allows connections on the configured `PORT`.
2. **First Device Only**: Currently, the service connects to the first Google Cast device it discovers. 
3. **mDNS**: Discovery relies on mDNS. In some network environments (like complex VLANs or some Docker setups), mDNS traffic might be blocked.
