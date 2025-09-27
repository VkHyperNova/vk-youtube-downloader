package main

import (
	"fmt"
	"log"
	"os"
	"vk-youtube-downloader/pkg/cmd"
	"vk-youtube-downloader/pkg/util"
)


func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <youtube-video-url>")
		return
	}

	videoURL := os.Args[1]

	// Get video title
	videoTitle, err := util.GetVideoTitle(videoURL)
	if err != nil {
		log.Fatalf("Error getting video title: %v", err)
	}

	// Sanitize video title for safe filename
	videoTitle = util.SanitizeFilename(videoTitle)

	// Determine the user's desktop path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	util.CreateDirectory(fmt.Sprintf("%s/Desktop/NewMusic", homeDir))
	tempAudioPath := fmt.Sprintf("%s/Desktop/NewMusic/%s.m4a", homeDir, videoTitle)
	finalMp3Path := fmt.Sprintf("%s/Desktop//NewMusic/%s.mp3", homeDir, videoTitle)

	fmt.Println("Downloading audio using yt-dlp...")

	err = cmd.DownloadAudioWithYtDlp(videoURL, tempAudioPath)
	if err != nil {
		log.Fatalf("Error downloading audio: %v", err)
	}

	fmt.Println("Converting audio to MP3...")

	err = util.ConvertToMp3(tempAudioPath, finalMp3Path)
	if err != nil {
		os.Remove(tempAudioPath)
		log.Fatalf("Error converting to MP3: %v", err)
	}

	// Clean up the temporary audio file
	os.Remove(tempAudioPath)

	fmt.Printf("Downloaded and converted audio to %s\n", finalMp3Path)
}

