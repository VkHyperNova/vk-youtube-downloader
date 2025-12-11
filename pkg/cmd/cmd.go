package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
)

func DownloadAudioWithYtDlp(videoURL, outputPath string) error {
	
	cmd := exec.Command("yt-dlp", "-U",
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