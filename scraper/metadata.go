package scraper

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type MetadataFile struct {
	Name  string `json:"name"`
	Size  int64  `json:"size"`
	SHA1  string `json:"sha1"`
	MD5   string `json:"md5"`
	CRC32 string `json:"crc32"`
}
type MetadataResponse struct {
	Files []MetadataFile `json:"files"`
}

const baseMetadataURL = "https://archive.org/metadata/"

// Custom Unmarshal function to convert size (string → int64)
func (f *MetadataFile) UnmarshalJSON(data []byte) error {
	// Temporary struct to match JSON structure but keep size as a string
	type Alias struct {
		Name  string `json:"name"`
		Size  string `json:"size"`
		SHA1  string `json:"sha1"`
		MD5   string `json:"md5"`
		CRC32 string `json:"crc32"`
	}
	aux := &Alias{}

	// Unmarshal into the temporary struct
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Convert size from string to int64
	if aux.Size != "" {
		sizeInt, err := strconv.ParseInt(aux.Size, 10, 64)
		if err != nil {
			return err
		}
		f.Size = sizeInt
	}

	// Assign values to the main struct
	f.Name = aux.Name
	f.SHA1 = aux.SHA1
	f.MD5 = aux.MD5
	f.CRC32 = aux.CRC32

	return nil
}

func GetMetadata(itemID string) MetadataResponse {
	metadataURL := baseMetadataURL + itemID + "/"
	resp, err := http.Get(metadataURL)
	if err != nil {
		log.Fatalf("Error fetching metadata: %v", err)
	}
	defer resp.Body.Close()

	var metadata MetadataResponse
	err = json.NewDecoder(resp.Body).Decode(&metadata)
	if err != nil {
		log.Fatalf("Error parsing metadata: %v", err)
	}
	return metadata
}
