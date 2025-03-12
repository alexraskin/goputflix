# GoPutFlix

A command-line application to browse and stream videos from your Put.io account directly to VLC media player.

## Features

- Browse your Put.io files and folders
- Stream videos directly to VLC without downloading
- Cross-platform support (Windows, macOS, Linux)
- Simple CLI navigation

## Requirements

- [VLC Media Player](https://www.videolan.org/vlc/)
- [Put.io API token](https://app.put.io/settings/account/oauth/apps)

## Installation

### Using Homebrew (macOS & Linux)

```bash
brew tap alexraskin/goputflix
brew install goputflix
```

### Direct Download

Download the appropriate binary for your platform from the [Releases](https://github.com/alexraskin/goputflix/releases) page.

#### Linux and macOS

```bash
# Download the latest release (adjust the version and architecture as needed)
curl -L https://github.com/alexraskin/goputflix/releases/latest/download/goput-flix_Linux_x86_64.tar.gz -o goput-flix.tar.gz

# Extract the archive
tar -xzf goput-flix.tar.gz

# Make the binary executable
chmod +x goput-flix

# Optional: Move to a directory in your PATH
sudo mv goput-flix /usr/local/bin/
```

#### Windows

Download the ZIP file from the [Releases](https://github.com/alexraskin/goputflix/releases) page and extract it to a location of your choice.

### Building from Source

Requires [Go](https://golang.org/dl/) 1.24 or higher.

```bash
# Clone the repository
git clone https://github.com/alexraskin/goputflix.git
cd goputflix

# Build the application
go build -o goput-flix ./cmd/main.go
```

## Usage

You can run the application in two ways:

1. Using a command-line flag:
```bash
goput-flix -token=YOUR_PUTIO_API_TOKEN
```

2. Using an environment variable:
```bash
export PUTIO_TOKEN=YOUR_PUTIO_API_TOKEN
goput-flix
```

### Navigation

Once the application is running, you can navigate your Put.io files with these commands:

- `[d#]` - Enter a directory (e.g., `d1` to enter the first directory)
- `[v#]` - Play a video with VLC (e.g., `v1` to play the first video)
- `[b]` - Go back to the parent directory
- `[q]` - Quit the application
- `[h]` - Show help

## How to Get a Put.io API Token

1. Go to [Put.io OAuth Apps](https://app.put.io/settings/account/oauth/apps)
2. Click "Create a new OAuth App"
3. Fill in the required information
4. After creation, copy the OAuth token 
- In the callback URI - `localhost:3000` or anything

## License

MIT [LICENSE](LICENSE)

## Acknowledgements

- [Put.io](https://put.io) for their API
- [go-putio](https://github.com/putdotio/go-putio) for the Go SDK 