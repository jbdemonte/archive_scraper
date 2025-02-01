package scraper

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"hash/crc32"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

const tmpFile = "_.tmp"
const bufferSize = 4 * 1024 * 1024
const baseDownloadURL = "https://archive.org/download/"

var createdDirs = make(map[string]bool)

// Ensure the directory exists, but only create it once
func ensureDirCached(filePath string) error {
	dir := filepath.Dir(filePath)

	// Skip if already created
	if createdDirs[dir] {
		return nil
	}

	// Create directory
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	// Cache the created directory
	createdDirs[dir] = true
	return nil
}

func isFileDownloadedAndValid(destPath string, file MetadataFile, noChecksum bool) bool {
	filePath := path.Join(destPath, file.Name)
	if !fileExists(filePath) {
		return false
	}
	if noChecksum {
		return true
	}
	err := verifyOrDeleteFile(filePath, file)
	return err == nil
}

func downloadFile(ctx context.Context, downloadChan chan int64, itemID, destPath string, file MetadataFile, noChecksum bool) error {
	fileURL := baseDownloadURL + itemID + "/" + url.PathEscape(file.Name)
	tmpFilePath := path.Join(destPath, tmpFile)
	destFilePath := path.Join(destPath, file.Name)

	err := deleteFile(tmpFilePath)
	if err != nil {
		return err
	}

	err = ensureDirCached(destFilePath)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", fileURL, nil)
	if err != nil {
		return err
	}

	// Set headers to mimic a real browser request
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "https://archive.org/")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check if request is successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error: %d %s", resp.StatusCode, resp.Status)
	}

	// Create destination file
	out, err := os.Create(tmpFilePath)
	if err != nil {
		return err
	}

	closeAndRemoveTmpFile := func() {
		_ = out.Close()
		_ = os.Remove(tmpFilePath)
	}

	closeAndRenameTmpFile := func() error {
		err := out.Close()
		if err != nil {
			return err
		}
		err = os.Rename(tmpFilePath, destFilePath)
		if err != nil {
			_ = os.Remove(tmpFilePath)
			return err
		}
		return nil
	}

	buffer := make([]byte, bufferSize)
	var downloaded int64

	// Copy data in chunks
	var n int
	var readErr error
	for {
		select {
		case <-ctx.Done():
			// task is cancelled
			closeAndRemoveTmpFile()
			return errors.New("download cancelled")
		default:
			// Read from response body into buffer
			if file.Size == 0 || file.Size-downloaded < bufferSize {
				n, readErr = resp.Body.Read(buffer)
			} else {
				n, readErr = io.ReadFull(resp.Body, buffer)
			}
			if n > 0 {
				// Write to file
				if _, writeErr := out.Write(buffer[:n]); writeErr != nil {
					closeAndRemoveTmpFile()
					return writeErr
				}
				downloaded += int64(n)

				// Send progress update if size is known
				downloadChan <- downloaded
			}

			// Stop reading at EOF
			if readErr == io.EOF {
				// verify and finalise the file
				err = closeAndRenameTmpFile()
				if err == nil && noChecksum == false {
					err = verifyOrDeleteFile(destFilePath, file)
				}
				return err
			}

			if readErr != nil {
				closeAndRemoveTmpFile()
				return readErr
			}
		}
	}
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func deleteFile(filePath string) error {
	if fileExists(filePath) {
		return os.Remove(filePath)
	}
	return nil
}

func verifyOrDeleteFile(filePath string, file MetadataFile) error {
	err := verifyFile(filePath, file)
	if err != nil {
		deleteFile(filePath)
		return err
	}
	return nil
}

func verifyFile(filePath string, file MetadataFile) error {
	var expectedHash string
	var hashType string

	if file.MD5 != "" {
		expectedHash = file.MD5
		hashType = "md5"
	} else if file.SHA1 != "" {
		expectedHash = file.SHA1
		hashType = "sha1"
	} else if file.CRC32 != "" {
		expectedHash = file.CRC32
		hashType = "crc32"
	}
	result, err := verifyFileChecksum(filePath, expectedHash, hashType)
	if err != nil {
		return err
	}
	if result == false {
		return fmt.Errorf("checksum verification failed (%s)", hashType)
	}
	return nil
}

func verifyFileChecksum(filePath, expectedHash string, hashType string) (bool, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	var hashFunc hash.Hash

	// Select SHA1 or MD5 based on input
	switch hashType {
	case "sha1":
		hashFunc = sha1.New()
	case "md5":
		hashFunc = md5.New()
	case "crc32":
		hashFunc = crc32.NewIEEE() // CRC32 IEEE standard
	default:
		return false, fmt.Errorf("unsupported hash type: %s", hashType)
	}

	// Compute the hash
	if _, err := io.Copy(hashFunc, file); err != nil {
		return false, err
	}
	computedHash := hex.EncodeToString(hashFunc.Sum(nil))

	// Compare with expected hash
	return computedHash == expectedHash, nil
}
