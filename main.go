package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// downloadAudioWithYtDlp downloads the audio stream using yt-dlp
func downloadAudioWithYtDlp(videoURL, outputPath string) error {
    cmd := exec.Command("yt-dlp", "--no-playlist", "-x", "--audio-format", "m4a", "-o", outputPath, videoURL)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("yt-dlp error: %v\n%s", err, string(output))
    }
    return nil
}

// convertToMp3 converts the given audio file to MP3 format using ffmpeg.
func convertToMp3(inputPath, outputPath string) error {
	cmd := exec.Command("ffmpeg", "-i", inputPath, "-q:a", "0", "-map", "a", outputPath)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error converting to MP3: %v", err)
	}
	return nil
}

// getVideoTitle fetches the title of a YouTube video using yt-dlp
func getVideoTitle(videoURL string) (string, error) {
    cmd := exec.Command("yt-dlp", "--no-playlist", "--print", "%(title)s", videoURL)
    output, err := cmd.Output()
    if err != nil {
        return "", fmt.Errorf("yt-dlp error: %v\n%s", err, string(output))
    }
    return string(output), nil
}

// sanitizeFilename removes or replaces unsafe characters from filenames
func sanitizeFilename(name string) string {
	// Remove invalid characters such as slashes, colons, spaces, newlines, etc.
	re := regexp.MustCompile(`[<>:"/\\|?*\x00-\x1F\x7F]+`)
	// Replace the special characters with underscores
	name = re.ReplaceAllString(name, "")

	// Optionally, remove leading and trailing spaces
	name = strings.TrimSpace(name)

	return name
}


func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <youtube-video-url>")
		return
	}

	videoURL := os.Args[1]

	// Get video title
	videoTitle, err := getVideoTitle(videoURL)
	if err != nil {
		log.Fatalf("Error getting video title: %v", err)
	}

	// Sanitize video title for safe filename
	videoTitle = sanitizeFilename(videoTitle)

	// Determine the user's desktop path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	// Define output file paths
	tempAudioPath := fmt.Sprintf("%s/Desktop/%s.m4a", homeDir, videoTitle)
	finalMp3Path := fmt.Sprintf("%s/Desktop/%s.mp3", homeDir, videoTitle)

	fmt.Println("Downloading audio using yt-dlp...")

	err = downloadAudioWithYtDlp(videoURL, tempAudioPath)
	if err != nil {
		log.Fatalf("Error downloading audio: %v", err)
	}

	fmt.Println("Converting audio to MP3...")

	err = convertToMp3(tempAudioPath, finalMp3Path)
	if err != nil {
		log.Fatalf("Error converting to MP3: %v", err)
	}

	// Clean up the temporary audio file
	os.Remove(tempAudioPath)

	fmt.Printf("Downloaded and converted audio to %s\n", finalMp3Path)
}

