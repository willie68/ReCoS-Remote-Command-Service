package audio

import (
	"os"
	"strings"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/beep/wav"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
)

var (
	sampleRate beep.SampleRate
	oneSpeaker sync.Once
)

// InitAudioplayer initialise the open hardware monitor connection
func InitAudioplayer(extconfig map[string]interface{}) error {
	var err error
	value, ok := extconfig["audioplayer"]
	if ok {
		config := value.(map[string]interface{})
		if config != nil {
			clog.Logger.Debug("audio:audioplayer: found config")
			active, ok := config["active"].(bool)
			if !ok {
				active = false
			}
			if active {
				clog.Logger.Debug("audio:audioplayer: active")
				mysampleRate, ok := config["samplerate"].(int)
				if !ok {
					mysampleRate = 48000
				}
				sampleRate = beep.SampleRate(mysampleRate)
				oneSpeaker.Do(func() {
					speaker.Init(sampleRate, sampleRate.N(time.Second/10))
				})
			} else {
				return nil
			}
		}
	}
	return err
}

func PlayAudio(file string) {
	var format beep.Format
	var streamer beep.StreamCloser
	f, err := os.Open(file)
	if err != nil {
		clog.Logger.Errorf("error open audio file: %v", err)
		return
	}
	defer f.Close()

	if strings.HasSuffix(file, ".mp3") {
		streamer, format, err = mp3.Decode(f)
		if err != nil {
			clog.Logger.Errorf("error decoding mp3 audio file: %v", err)
			return
		}
	}

	if strings.HasSuffix(file, ".wav") {
		streamer, format, err = wav.Decode(f)
		if err != nil {
			clog.Logger.Errorf("error decoding wav audio file: %v", err)
			return
		}
	}

	if strings.HasSuffix(file, ".flac") {
		streamer, format, err = flac.Decode(f)
		if err != nil {
			clog.Logger.Errorf("error decoding flac audio file: %v", err)
			return
		}
	}

	if strings.HasSuffix(file, ".ogg") {
		streamer, format, err = vorbis.Decode(f)
		if err != nil {
			clog.Logger.Errorf("error decoding ogg audio file: %v", err)
			return
		}
	}

	defer streamer.Close()

	done := make(chan bool)
	if format.SampleRate == sampleRate {
		speaker.Play(beep.Seq(streamer, beep.Callback(func() {
			done <- true
		})))
	} else {
		resampled := beep.Resample(4, format.SampleRate, sampleRate, streamer)
		speaker.Play(beep.Seq(resampled, beep.Callback(func() {
			done <- true
		})))
	}
	<-done
}
