package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func downloadAudioWithYtDlp(videoURL, outputPath string) error {
	cmd := exec.Command("yt-dlp",
		"--no-playlist",
		"--newline",
		"--progress",
		"-x",
		"--audio-format", "m4a",
		"-o", outputPath,
		videoURL,
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %v", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start yt-dlp: %v", err)
	}

	printOutput := func(pipeReader io.Reader) {
		scanner := bufio.NewScanner(pipeReader)
		re := regexp.MustCompile(`\[download\]\s+(\d+\.\d+)%\s+of\s+([\d\.]+[KMG]iB)\s+at\s+([\d\.]+[KMG]iB/s)`)

		writer := bufio.NewWriter(os.Stdout)

		for scanner.Scan() {
			line := scanner.Text()
			if matches := re.FindStringSubmatch(line); len(matches) == 4 {
				// Overwrite current line
				fmt.Fprintf(writer, "\rDownloading: %s%% of %s at %s   ", matches[1], matches[2], matches[3])
				writer.Flush()
			} else {
				// Print other messages normally
				fmt.Fprintln(writer, line)
				writer.Flush()
			}
		}
		fmt.Fprintln(writer) // finish with newline
		writer.Flush()
	}

	go printOutput(stdout)
	go printOutput(stderr)

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("yt-dlp finished with error: %v", err)
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

