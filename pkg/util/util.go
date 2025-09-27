package util

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

// convertToMp3 converts the given audio file to MP3 format using ffmpeg.
func ConvertToMp3(inputPath, outputPath string) error {
	cmd := exec.Command("ffmpeg", "-i", inputPath, "-q:a", "0", "-map", "a", outputPath)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error converting to MP3: %v", err)
	}
	return nil
}

// getVideoTitle fetches the title of a YouTube video using yt-dlp
func GetVideoTitle(videoURL string) (string, error) {
    cmd := exec.Command("yt-dlp", "--no-playlist", "--print", "%(title)s", videoURL)
    output, err := cmd.Output()
    if err != nil {
        return "", fmt.Errorf("yt-dlp error: %v\n%s", err, string(output))
    }
    return string(output), nil
}

// sanitizeFilename removes or replaces unsafe characters from filenames
func SanitizeFilename(name string) string {
	// Remove invalid characters such as slashes, colons, spaces, newlines, etc.
	re := regexp.MustCompile(`[<>:"/\\|?*\x00-\x1F\x7F]+`)
	// Replace the special characters with underscores
	name = re.ReplaceAllString(name, "")

	// Optionally, remove leading and trailing spaces
	name = strings.TrimSpace(name)

	return name
}

func ClearScreen() {

	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error clearing screen:", err)
	}
}

func CreateDirectory(dirName string) error {
	// MkdirAll creates the directory along with parents if needed
	err := os.MkdirAll(dirName, 0700)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dirName, err)
	}
	return nil
}