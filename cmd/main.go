package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"goputflix/internal"

	"github.com/putdotio/go-putio"
)

var (
	version string = "0.0.1"
	commit  string = ""
)

func main() {
	var token string
	var showHelp bool
	var showVersion bool

	flag.StringVar(&token, "token", "", "Put.io API token")
	flag.BoolVar(&showHelp, "help", false, "Show help message")
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.Parse()

	if showHelp {
		printHelp()
		return
	}

	if showVersion {
		fmt.Println("goput-flix version", version)
		fmt.Println("commit", commit)
		return
	}

	if token == "" {
		token = os.Getenv("PUTIO_TOKEN")
	}

	if token == "" {
		fmt.Println("Error: Put.io API token is required")
		fmt.Println("You can provide it using the -token flag or set the PUTIO_TOKEN environment variable")
		fmt.Println("\nFor more information, run: goput-flix -help")
		return
	}

	client := internal.InitPutio(internal.PutIoOptions{
		Token: token,
	})

	browsePutio(client, 0)
}

func printHelp() {
	fmt.Println("goput-flix: A CLI tool to browse your Put.io account and play videos using VLC")
	fmt.Println("\nUsage:")
	fmt.Println("  goput-flix [flags]")
	fmt.Println("\nFlags:")
	fmt.Println("  -token string   Put.io API token (can also be set via PUTIO_TOKEN environment variable)")
	fmt.Println("  -help           Show this help message")
	fmt.Println("\nNavigation:")
	fmt.Println("  [d#]   Navigate into a folder")
	fmt.Println("  [v#]   Play a video using VLC")
	fmt.Println("  [b]    Go back to the parent folder")
	fmt.Println("  [q]    Quit the application")
}

func browsePutio(client *putio.Client, folderID int64) {
	ctx := context.Background()
	var currentFolderPath string

	for {
		fmt.Println("\n----- Put.io Files -----")

		if folderID != 0 {
			folder, err := client.Files.Get(ctx, folderID)
			if err == nil {
				currentFolderPath = folder.Name
			}
		} else {
			currentFolderPath = "Root"
		}
		fmt.Printf("Location: %s\n\n", currentFolderPath)

		files, _, err := client.Files.List(ctx, folderID)
		if err != nil {
			log.Printf("Error fetching files: %v\n", err)
			fmt.Println("Press Enter to continue or 'q' to quit")
			if getUserInput() == "q" {
				return
			}
			continue
		}

		if len(files) == 0 {
			fmt.Println("No files found in this folder")
			fmt.Println("[b] Go back")
			fmt.Println("[q] Quit")

			choice := getUserInput()
			if choice == "q" {
				return
			} else if choice == "b" && folderID != 0 {
				// Go to parent folder if not in root
				parent, err := client.Files.Get(ctx, folderID)
				if err != nil {
					log.Printf("Error fetching parent folder: %v\n", err)
					folderID = 0 // Go to root if there's an error
				} else {
					folderID = parent.ParentID
				}
			} else if choice == "b" {
				// If we're at root, just continue to show root again
				continue
			}
			continue
		}

		fmt.Println("Options:")
		if folderID != 0 {
			fmt.Println("[b] Go back to parent folder")
		}
		fmt.Println("[q] Quit")
		fmt.Println("[h] Help")
		fmt.Println()

		var folders []putio.File
		var videos []putio.File
		var others []putio.File

		for _, file := range files {
			if file.ContentType == "application/x-directory" {
				folders = append(folders, file)
			} else if isVideoFile(file) {
				videos = append(videos, file)
			} else {
				others = append(others, file)
			}
		}

		if len(folders) > 0 {
			fmt.Println("Folders:")
			for i, folder := range folders {
				fmt.Printf("  [d%d] %s\n", i+1, folder.Name)
			}
			fmt.Println()
		}

		if len(videos) > 0 {
			fmt.Println("Videos:")
			for i, video := range videos {
				size := float64(video.Size) / 1048576
				fmt.Printf("  [v%d] %s (%.2f MB)\n", i+1, video.Name, size)
			}
			fmt.Println()
		}

		// Display other files
		if len(others) > 0 {
			fmt.Println("Other Files:")
			for i, file := range others {
				size := float64(file.Size) / 1048576
				fmt.Printf("  [o%d] %s (%.2f MB)\n", i+1, file.Name, size)
			}
			fmt.Println()
		}

		fmt.Print("Enter your choice: ")
		choice := getUserInput()

		if choice == "q" {
			return
		}

		if choice == "h" {
			printHelp()
			fmt.Println("\nPress Enter to continue")
			getUserInput()
			continue
		}

		if choice == "b" {
			if folderID == 0 {
				continue // Already at root
			}
			// go to parent folder
			parent, err := client.Files.Get(ctx, folderID)
			if err != nil {
				log.Printf("Error fetching parent folder: %v\n", err)
				folderID = 0 // Go to root if there's an error
			} else {
				folderID = parent.ParentID
			}
			continue
		}

		// directory selection
		if strings.HasPrefix(choice, "d") {
			index, err := strconv.Atoi(choice[1:])
			if err != nil || index < 1 || index > len(folders) {
				fmt.Println("Invalid folder selection")
				continue
			}
			folderID = folders[index-1].ID
			continue
		}

		if strings.HasPrefix(choice, "v") {
			index, err := strconv.Atoi(choice[1:])
			if err != nil || index < 1 || index > len(videos) {
				fmt.Println("Invalid video selection")
				continue
			}

			selectedVideo := videos[index-1]
			playVideo(client, selectedVideo)
		}

		if strings.HasPrefix(choice, "o") {
			fmt.Println("This file type cannot be played with VLC")
		}
	}
}

func getUserInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func playVideo(client *putio.Client, video putio.File) {
	fmt.Printf("\nPreparing to play: %s\n", video.Name)

	ctx := context.Background()

	hlsStream, err := client.Files.HLSPlaylist(ctx, video.ID, "all")
	if err != nil {
		log.Printf("Error getting HLS playlist: %v\n", err)
		return
	}
	defer hlsStream.Close()

	playlistContent, err := io.ReadAll(hlsStream)
	if err != nil {
		log.Printf("Error reading HLS playlist: %v\n", err)
		return
	}

	tmpFile, err := os.CreateTemp("", "putio-playlist-*.m3u8")
	if err != nil {
		log.Printf("Error creating temporary file: %v\n", err)
		return
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(playlistContent); err != nil {
		tmpFile.Close()
		log.Printf("Error writing to temporary file: %v\n", err)
		return
	}
	if err := tmpFile.Close(); err != nil {
		log.Printf("Error closing temporary file: %v\n", err)
		return
	}

	fmt.Printf("Playing: %s\n", video.Name)
	err = internal.PlayVLC(tmpFile.Name())
	if err != nil {
		log.Printf("Error playing video: %v\n", err)
		return
	}

	fmt.Println("Video playback started in VLC")
	fmt.Println("Press Enter to continue browsing")
	getUserInput()
}

func isVideoFile(file putio.File) bool {
	if strings.HasPrefix(file.ContentType, "video/") {
		return true
	}

	name := strings.ToLower(file.Name)
	extensions := []string{".mp4", ".mkv", ".avi", ".mov", ".wmv", ".m4v", ".flv", ".webm"}

	for _, ext := range extensions {
		if strings.HasSuffix(name, ext) {
			return true
		}
	}

	return false
}
