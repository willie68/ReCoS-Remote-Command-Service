package audio

import (
	"os"
	"strconv"
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
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

var AudioPlayerIntegInfo = models.IntegInfo{
	Category:    "Audio-Video",
	Name:        "audioplayer",
	Description: "AudioPlayer is a integration for playing audio files. For the audioplayer i need simply the sample rate to work with. <br />For convinience you can only switch between 44,1kHz and 48kHz. <br />",
	Image:       "speaker.svg",
	Parameters: []models.ParamInfo{
		{
			Name:           "active",
			Type:           "bool",
			Description:    "activate the open hardwaremonitor",
			WizardPossible: false,
			List:           make([]string, 0),
		},
		{
			Name:           "samplerate",
			Type:           "string",
			Description:    "the samplerate to use. 44100 or 48000 are ok.",
			WizardPossible: false,
			List:           []string{"48000", "44100"},
		},
	},
}

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
				mysampleRate := 48000
				mysampleRate, ok := config["samplerate"].(int)
				if !ok {
					valuestr, ok := config["samplerate"].(string)
					if ok {
						mysampleRate, _ = strconv.Atoi(valuestr)
					}
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
