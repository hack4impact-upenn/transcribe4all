package transcription

import (
	"bufio"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/net/websocket"
)

// TranscribeWithIBM transcribes a given audio file using the IBM Watson
// Speech To Text API
func TranscribeWithIBM(filePath string, IBMUsername string, IBMPassword string) (string, error) {
	IBMAuthToken, err := getIBMAuthToken(IBMUsername, IBMPassword)

	url, err := generateIBMURL(IBMAuthToken)
	if err != nil {
		return "", err
	}
	ws, err := websocket.Dial(url, "", "http://localhost:8000")
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
	log.Println("File uploaded")

	ws.Write([]byte{}) // write empty message to indicate end of uploading file

	transcriptionRes, err := pollForTranscriptionResult(ws)
	if err != nil {
		return "", err
	}

	return transcriptionRes, nil
}

func getIBMAuthToken(IBMUsername string, IBMPassword string) (string, error) {
	baseURL, err := url.Parse("https://stream.watsonplatform.net/authorization/api/v1/token")
	if err != nil {
		return "", err
	}
	params := url.Values{}
	params.Add("url", "https://stream.watsonplatform.net/speech-to-text/api")
	baseURL.RawQuery = params.Encode()

	req, err := http.NewRequest(http.MethodGet, baseURL.String(), nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(IBMUsername, IBMPassword)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	token, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(token), nil
}

func generateIBMURL(IBMAuthToken string) (string, error) {
	baseURL, err := url.Parse("wss://stream.watsonplatform.net/speech-to-text/api/v1/recognize")
	if err != nil {
		return "", err
	}
	params := url.Values{}
	params.Add("watson-token", IBMAuthToken)
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
		log.Println("Waiting...")
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
