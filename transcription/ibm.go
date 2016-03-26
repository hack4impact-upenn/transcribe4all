package transcription

import (
	"bufio"
	"encoding/json"
	"io"
	"net/url"
	"os"

	"golang.org/x/net/websocket"
)

// TranscribeWithIBM transcribes a given audio file using the IBM Watson
// Speech To Text API
func TranscribeWithIBM(filePath string, ibmAuthToken string) (string, error) {
	url, err := generateIBMURL(ibmAuthToken)
	if err != nil {
		return "", err
	}

	ws, err := websocket.Dial(url, "", "")
	if err != nil {
		return "", err
	}
	defer ws.Close()

	requestArgs, err := json.Marshal(map[string]string{
		"action":           "start",
		"content-type":     "audio/flac",
		"continuous":       "true",
		"word_confidence":  "true",
		"timestamps":       "true",
		"profanity_filter": "false",
		"interim_results":  "false",
	})
	if err != nil {
		return "", err
	}
	if _, err = ws.Write(requestArgs); err != nil {
		return "", err
	}

	if err = uploadBinaryWithWebsocket(ws, filePath); err != nil {
		return "", err
	}

	ws.Write([]byte{}) // write empty message to indicate end of uploading file

	transcriptionRes, err := pollForTranscriptionResult(ws)
	if err != nil {
		return "", err
	}

	return transcriptionRes, nil
}

func generateIBMURL(ibmAuthToken string) (string, error) {
	baseURL, err := url.Parse("wss://stream.watsonplatform.net/speech-to-text/api/v1/recognize")
	if err != nil {
		return "", err
	}
	params := url.Values{}
	params.Add("watson-token", ibmAuthToken)
	params.Add("model", "en-US_BroadbandModel")
	baseURL.RawQuery = params.Encode()
	return baseURL.String(), nil
}

func uploadBinaryWithWebsocket(ws *websocket.Conn, filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	r := bufio.NewReader(f)
	buffer := make([]byte, 2048)

	for {
		n, err := r.Read(buffer)
		if n == 0 {
			break
		}
		if err != nil && err != io.EOF {
			return err
		}
		ws.Write(buffer)
	}
	return nil
}

func pollForTranscriptionResult(ws *websocket.Conn) (string, error) {
	transcriptionRes := []byte{}
	for {
		n, err := ws.Read(transcriptionRes)
		if err != nil {
			return "", err
		}
		if n > 0 {
			break
		}
	}
	return string(transcriptionRes), nil
}
