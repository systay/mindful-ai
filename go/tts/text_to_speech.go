package tts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Replace with your ElevenLabs API key
const apiKey = "YOUR_ELEVENLABS_API_KEY"

// APIEndpoint is the ElevenLabs TTS API endpoint
const APIEndpoint = "https://api.elevenlabs.io/v1/text-to-speech/YOUR_VOICE_ID"

// TTSRequest represents the request payload for the TTS API
type TTSRequest struct {
	Text          string `json:"text"`
	VoiceSettings struct {
		Stability  float32 `json:"stability"`
		Similarity float32 `json:"similarity_boost"`
	} `json:"voice_settings"`
}

func main() {
	// Example text to convert to speech
	text := "Hello, this is a sample text to speech conversion using ElevenLabs API and Go."

	// Call the function to get audio and save to a file
	err := TextToSpeech(text, "output_audio.mp3")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Audio saved successfully as output_audio.mp3")
	}
}

// TextToSpeech takes a string and converts it to speech using ElevenLabs API
// It saves the resulting audio to the specified filename
func TextToSpeech(text, filename string) error {
	// Create the request payload
	var reqPayload TTSRequest
	reqPayload.Text = text
	reqPayload.VoiceSettings.Stability = 0.75  // Adjust as needed
	reqPayload.VoiceSettings.Similarity = 0.75 // Adjust as needed

	payloadBytes, err := json.Marshal(reqPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal request payload: %v", err)
	}

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", APIEndpoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set the necessary headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("xi-api-key", apiKey)

	// Initialize HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("non-200 response: %d - %s", resp.StatusCode, string(bodyBytes))
	}

	// Create the output file
	outFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outFile.Close()

	// Stream the response body to the file
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write audio to file: %v", err)
	}

	return nil
}
