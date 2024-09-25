package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"

	"github.com/eryalito/vigo-bus-core/pkg/api"
)

func PngToBase64(img image.Image) (string, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func GenerateImageWithMarkers(apiKey string, origin struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}, stops []api.Stop) (image.Image, error) {
	// Construct the Google Maps Static API URL
	baseURL := "https://maps.googleapis.com/maps/api/staticmap"
	size := "600x400"
	markers := fmt.Sprintf("markers=color:red|label:O|%f,%f", origin.Lat, origin.Lon)

	for index, stop := range stops {
		markers += fmt.Sprintf("&markers=color:green|label:%d|%f,%f", index+1, stop.Location.Lat, stop.Location.Lon)
	}

	url := fmt.Sprintf("%s?size=%s&%s&key=%s", baseURL, size, markers, apiKey)

	// Make the HTTP request to the Google Maps Static API
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image from Google Maps Static API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Google Maps Static API returned status code %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Decode the image
	img, _, err := image.Decode(bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}

	return img, nil
}
