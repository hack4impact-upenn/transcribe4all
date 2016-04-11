package transcription

import (
	"bufio"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

// IBMResult is the result of an IBM transcription. See
// https://www.ibm.com/smarterplanet/us/en/ibmwatson/developercloud/doc/speech-to-text/output.shtml
// for details.
type IBMResult struct {
	ResultIndex int              `json:"result_index"`
	Results     []ibmResultField `json:"results"`
}
type ibmResultField struct {
	Alternatives []ibmAlternativesField `json:"alternatives"`
	Final        bool                   `json:"final"`
}
type ibmAlternativesField struct {
	WordConfidence    []ibmWordConfidence `json:"word_confidence"`
	OverallConfidence float64             `json:"confidence"`
	Transcript        string              `json:"transcript"`
	Timestamps        []ibmWordTimestamp  `json:"timestamps"`
}
type ibmWordConfidence []interface{}
type ibmWordTimestamp []interface{}

// TranscribeWithIBM transcribes a given audio file using the IBM Watson
// Speech To Text API
func TranscribeWithIBM(filePath string, IBMUsername string, IBMPassword string) (IBMResult, error) {
	result := IBMResult{}

	url := "wss://stream.watsonplatform.net/speech-to-text/api/v1/recognize?model=en-US_BroadbandModel"
	header := http.Header{}
	header.Set("Authorization", "Basic "+basicAuth(IBMUsername, IBMPassword))

	dialer := websocket.DefaultDialer
	ws, _, err := dialer.Dial(url, header)
	if err != nil {
		return result, err
	}
	defer ws.Close()

	requestArgs := map[string]interface{}{
		"action":             "start",
		"content-type":       "audio/flac",
		"continuous":         true,
		"word_confidence":    true,
		"timestamps":         true,
		"profanity_filter":   false,
		"interim_results":    false,
		"inactivity_timeout": -1,
	}
	if err = ws.WriteJSON(requestArgs); err != nil {
		return result, err
	}
	if err = uploadFileWithWebsocket(ws, filePath); err != nil {
		return result, err
	}

	ws.WriteMessage(websocket.BinaryMessage, []byte{}) // write empty message to indicate end of uploading file
	log.Println("File uploaded")

	// IBM must receive a message every 30 seconds or it will close the websocket.
	// This code concurrently writes a message every 5 second until returning.
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go keepConnectionOpen(ws, ticker, quit)
	defer close(quit)

	for {
		err := ws.ReadJSON(&result)
		if err != nil {
			return result, err
		}
		if len(result.Results) > 0 {
			return result, nil
		}
	}
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func uploadFileWithWebsocket(ws *websocket.Conn, filePath string) error {
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

func keepConnectionOpen(ws *websocket.Conn, ticker *time.Ticker, quit chan struct{}) {
	for {
		select {
		case <-ticker.C:
			err := ws.WriteJSON(map[string]string{
				"action": "no-op",
			})
			if err != nil {
				return
			}
		case <-quit:
			ticker.Stop()
			return
		}
	}
}
