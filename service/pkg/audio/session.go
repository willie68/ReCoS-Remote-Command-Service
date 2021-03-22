package audio

import (
	"strings"
)

// Session represents a single addressable audio session
type Session interface {
	GetVolume() float32
	SetVolume(v float32) error

	GetMute() bool
	SetMute(m bool) error

	IsInput() bool

	Key() string
	Release()

	String() string
}

const (

	// ideally these would share a common ground in baseSession
	// but it will not call the child GetVolume correctly :/
	sessionCreationLogMessage = "Created audio session instance"

	// format this with s.humanReadableDesc and whatever the current volume is
	sessionStringFormat = "<session: %s, Input: %t, mute: %t, vol: %.2f>"
)

type baseSession struct {
	inputDevice bool
	system      bool
	master      bool

	// used by Key(), needs to be set by child
	name string

	// used by String(), needs to be set by child
	humanReadableDesc string
}

func (s *baseSession) Key() string {
	if s.system {
		return systemSessionName
	}

	if s.master {
		return strings.ToLower(s.name) // could be master or mic, or any device's friendly name
	}

	return strings.ToLower(s.name)
}

func (s *baseSession) IsInput() bool {
	return s.inputDevice
}
