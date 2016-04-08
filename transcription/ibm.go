package transcription

import (
	"bufio"
	"encoding/base64"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

// TranscribeWithIBM transcribes a given audio file using the IBM Watson
// Speech To Text API
func TranscribeWithIBM(filePath string, IBMUsername string, IBMPassword string) (string, error) {
	url := "wss://stream.watsonplatform.net/speech-to-text/api/v1/recognize?model=es-ES_BroadbandModel"
	header := http.Header{}
	header.Set("Authorization", "Basic "+basicAuth(IBMUsername, IBMPassword))

	dialer := websocket.DefaultDialer
	ws, _, err := dialer.Dial(url, header)
	if err != nil {
		return "", err
	}
	defer ws.Close()

	requestArgs := map[string]interface{}{
		"action":           "start",
		"content-type":     "audio/flac",
		"continuous":       true,
		"word_confidence":  true,
		"timestamps":       true,
		"profanity_filter": false,
		"interim_results":  false,
	}
	if err = ws.WriteJSON(requestArgs); err != nil {
		return "", err
	}
	if err = uploadBinaryWithWebsocket(ws, filePath); err != nil {
		return "", err
	}

	ws.WriteMessage(websocket.BinaryMessage, []byte{}) // write empty message to indicate end of uploading file
	transcriptionRes, err := pollForTranscriptionResult(ws)
	if err != nil {
		return "", err
	}
	return transcriptionRes, nil
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
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
		ws.WriteMessage(websocket.BinaryMessage, buffer)
	}
	return nil
}

func pollForTranscriptionResult(ws *websocket.Conn) (string, error) {
	for {
		_, transcriptionRes, err := ws.ReadMessage()
		if err != nil {
			return "", err
		}
		// BUG(sandlerben): This is a hack which will not work if the transcription contains "listening"
		if len(transcriptionRes) > 0 && !strings.Contains(string(transcriptionRes), "listening") {
			return string(transcriptionRes), nil
		}
	}
}
