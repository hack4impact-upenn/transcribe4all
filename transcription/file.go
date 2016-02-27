package transcription

import (
	"log"
	"os/exec"
)

func convertAudioIntoRequiredFormat(fn string) {
	// http://cmusphinx.sourceforge.net/wiki/faq
	// -ar 16000 sets frequency to required 16khz
	// -ac 1 sets the number of audio channels to 1
	cmd := exec.Command("ffmpeg", "-i", fn, "-ar", "16000", "-ac", "1", "file.wav")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
