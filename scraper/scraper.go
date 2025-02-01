package scraper

import (
	"context"
	"fmt"
	"strings"
)

var emptyLine = ""

func clearLine() {
	if emptyLine == "" {
		emptyLine = strings.Repeat(" ", getTerminalWidth())
	}
	fmt.Printf("\r%s\r", emptyLine)
}

func moveCursorUp(n int) {
	if n > 0 {
		fmt.Printf("\033[%dA", n)
	}
}

func moveCursorDown(n int) {
	if n > 0 {
		fmt.Printf("\033[%dB", n)
	}
}

func hideCursor() {
	fmt.Print("\033[?25l")
}

func showCursor() {
	fmt.Print("\033[?25h")
}

func DownloadAllFilesFromItemID(ctx context.Context, itemID string, destPath string, noChecksum bool) error {
	metadata := GetMetadata(itemID)

	totalSize := int64(0)
	for _, file := range metadata.Files {
		totalSize += file.Size
	}

	totalDownloaded := int64(0)
	totalExisting := int64(0)

	existingCount := 0
	successCount := 0
	errorCount := 0

	logEntries := []string{}

	// Hide the cursor to avoid flickering
	hideCursor()

	fmt.Printf("\n")

	headerLineCount := 5
	printHeader := func(processing int) {
		clearLine()
		fmt.Printf("Progress               : %s / %s (%d%%)\n", humanReadableSize(totalDownloaded+totalExisting), humanReadableSize(totalSize), int(100*(totalDownloaded+totalExisting)/totalSize))

		fmt.Printf("Files Already Present  : %d\n", existingCount)
		fmt.Printf("Successfully Downloaded: %d\n", successCount)
		fmt.Printf("Failed Downloads       : %d\n", errorCount)

		// no need to clean the line, the successCount &errorCount only increase
		if processing >= 0 {
			fmt.Printf("Processing             : %d / %d\n", processing+1, len(metadata.Files))
		} else {
			clearLine()
			fmt.Printf("Processing             : \n")
		}
	}

	defer func() {
		printHeader(-1)
		clearLine()
		fmt.Printf("\n")
		for _, logEntry := range logEntries {
			fmt.Println(logEntry)
		}
		showCursor()
	}()

	for i, file := range metadata.Files {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			printHeader(i)
			logPrefix := fmt.Sprintf("Downloading            : %s", file.Name)

			if file.Size > 0 {
				logPrefix += fmt.Sprintf(" (%s)", humanReadableSize(file.Size))
			}

			clearLine()
			fmt.Printf(logPrefix)

			downloadChan := make(chan int64)

			go func() {
				defer close(downloadChan)

				exist := isFileDownloadedAndValid(destPath, file, noChecksum)

				if exist {
					logEntries = append(logEntries, fmt.Sprintf("🟢 %s", file.Name))
					totalExisting += file.Size
					existingCount++
					return
				}

				err := downloadFile(ctx, downloadChan, itemID, destPath, file, noChecksum)
				if err != nil {
					logEntries = append(logEntries, fmt.Sprintf("❌ %s", file.Name))
					logEntries = append(logEntries, fmt.Sprintf("  ➜ %s", err.Error()))
					errorCount++
				} else {
					logEntries = append(logEntries, fmt.Sprintf("✅ %s", file.Name))
					totalDownloaded += file.Size
					successCount++
				}
			}()

			previousLog := logPrefix
			for progress := range downloadChan {
				log := ""
				if file.Size > 0 {
					log = fmt.Sprintf("%s (%d%%%%)", logPrefix, int(100*progress/file.Size))
				} else {
					log = fmt.Sprintf("%s (%s)", logPrefix, humanReadableSize(progress))
				}
				if log != previousLog {
					clearLine()
					fmt.Printf(log)
					previousLog = log
				}
			}

			moveCursorUp(headerLineCount)
		}
	}
	return nil
}
