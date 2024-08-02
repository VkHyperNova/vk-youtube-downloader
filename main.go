package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	

	"github.com/kkdai/youtube/v2"
)

// downloadAudioStream downloads the audio stream of a YouTube video.
func downloadAudioStream(video *youtube.Video, format *youtube.Format, outputPath string) error {
	client := youtube.Client{}

	resp, _, err := client.GetStream(video, format)
	if err != nil {
		return fmt.Errorf("error getting video stream: %v", err)
	}
	defer resp.Close()

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <youtube-video-url>")
		return
	}

	videoURL := os.Args[1]

	client := youtube.Client{}
	video, err := client.GetVideo(videoURL)
	if err != nil {
		log.Fatalf("Error getting video info: %v", err)
	}

	var audioFormat *youtube.Format
	for _, format := range video.Formats {
		if format.AudioChannels > 0 && format.ItagNo == 140 { // 140 typically represents m4a format
			audioFormat = &format
			break
		}
	}

	if audioFormat == nil {
		log.Fatalf("No suitable audio format found")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	
	tempAudioPath := homeDir + "/Desktop/" + video.Title + ".m4a"
	finalMp3Path := homeDir + "/Desktop/" + video.Title + ".mp3"

	

	err = downloadAudioStream(video, audioFormat, tempAudioPath)
	if err != nil {
		log.Fatalf("Error downloading audio stream: %v", err)
	}

	err = convertToMp3(tempAudioPath, finalMp3Path)
	if err != nil {
		log.Fatalf("Error converting to MP3: %v", err)
	}

	// Clean up the temporary file
	os.Remove(tempAudioPath)

	fmt.Printf("Downloaded and converted audio to %s\n", finalMp3Path)
}
