package scraper

import (
	"fmt"
	"golang.org/x/term"
	"os"
)

func humanReadableSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	suffixes := []string{"kB", "MB", "GB", "TB", "PB", "EB"}

	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.2f %s", float64(bytes)/float64(div), suffixes[exp])
}

func getTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 80 // Default to 80 columns if unknown
	}
	return width
}
