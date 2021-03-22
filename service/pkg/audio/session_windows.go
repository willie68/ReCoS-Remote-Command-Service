package audio

import (
	"errors"
	"fmt"

	ole "github.com/go-ole/go-ole"
	ps "github.com/mitchellh/go-ps"
	wca "github.com/moutend/go-wca"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
)

var errNoSuchProcess = errors.New("No such process")
var errRefreshSessions = errors.New("Trigger session refresh")

type wcaSession struct {
	baseSession

	pid         uint32
	processName string

	control *wca.IAudioSessionControl2
	volume  *wca.ISimpleAudioVolume

	eventCtx *ole.GUID
}

type masterSession struct {
	baseSession

	volume *wca.IAudioEndpointVolume

	eventCtx *ole.GUID

	stale bool // when set to true, we should refresh sessions on the next call to SetVolume
}

func newWCASession(control *wca.IAudioSessionControl2, volume *wca.ISimpleAudioVolume, pid uint32, eventCtx *ole.GUID) (*wcaSession, error) {

	s := &wcaSession{
		control:  control,
		volume:   volume,
		pid:      pid,
		eventCtx: eventCtx,
	}

	// special treatment for system sounds session
	if pid == 0 {
		s.system = true
		s.name = systemSessionName
		s.humanReadableDesc = "system sounds"
	} else {

		// find our session's process name
		process, err := ps.FindProcess(int(pid))
		if err != nil {
			defer s.Release()

			return nil, fmt.Errorf("find process name by pid: %w", err)
		}

		// this PID may be invalid - this means the process has already been
		// closed and we shouldn't create a session for it.
		if process == nil {
			return nil, errNoSuchProcess
		}

		s.processName = process.Executable()
		s.name = s.processName
		s.humanReadableDesc = fmt.Sprintf("%s (pid %d)", s.processName, s.pid)
	}

	clog.Logger.Debugf(sessionCreationLogMessage+"\r\nsession:%v", s)

	return s, nil
}

func newMasterSession(volume *wca.IAudioEndpointVolume, eventCtx *ole.GUID, key string) (*masterSession, error) {

	s := &masterSession{
		volume:   volume,
		eventCtx: eventCtx,
	}

	s.master = true
	s.name = key
	s.humanReadableDesc = key

	clog.Logger.Debugf(sessionCreationLogMessage+"\r\nsession: %v", s)

	return s, nil
}

func (s *wcaSession) GetVolume() float32 {
	var level float32

	if err := s.volume.GetMasterVolume(&level); err != nil {
		clog.Logger.Alertf("Failed to get session volume: %v", err)
	}

	return level
}

func (s *wcaSession) SetVolume(v float32) error {
	if err := s.volume.SetMasterVolume(v, s.eventCtx); err != nil {
		return fmt.Errorf("adjust session volume: %w", err)
	}

	// mitigate expired sessions by checking the state whenever we change volumes
	var state uint32

	if err := s.control.GetState(&state); err != nil {
		return fmt.Errorf("get session state: %w", err)
	}

	if state == wca.AudioSessionStateExpired {
		clog.Logger.Alert("Audio session expired, triggering session refresh")
		return errRefreshSessions
	}

	clog.Logger.Debugf("Adjusting session volume to %.2f", v)

	return nil
}

func (s *wcaSession) GetMute() bool {
	var isMuted bool

	if err := s.volume.GetMute(&isMuted); err != nil {
		clog.Logger.Alertf("Failed to get session mute: %v", err)
	}

	return isMuted
}

func (s *wcaSession) SetMute(mute bool) error {
	if err := s.volume.SetMute(mute, s.eventCtx); err != nil {
		return fmt.Errorf("set session mute: %w", err)
	}

	clog.Logger.Debugf("Setting session mute to %t", mute)

	return nil
}

func (s *wcaSession) Release() {
	clog.Logger.Debug("Releasing audio session")

	s.volume.Release()
	s.control.Release()
}

func (s *wcaSession) String() string {
	return fmt.Sprintf(sessionStringFormat, s.humanReadableDesc, s.IsInput(), s.GetMute(), s.GetVolume())
}

func (s *masterSession) GetVolume() float32 {
	var level float32

	if err := s.volume.GetMasterVolumeLevelScalar(&level); err != nil {
		clog.Logger.Alertf("Failed to get session volume: %v", err)
	}

	return level
}

func (s *masterSession) SetVolume(v float32) error {
	if s.stale {
		clog.Logger.Alert("Session expired because default device has changed, triggering session refresh")
		return errRefreshSessions
	}

	if err := s.volume.SetMasterVolumeLevelScalar(v, s.eventCtx); err != nil {
		return fmt.Errorf("adjust session volume: %w", err)
	}

	clog.Logger.Debugf("Adjusting session volume to %.2f", v)

	return nil
}

func (s *masterSession) GetMute() bool {
	var isMuted bool

	if err := s.volume.GetMute(&isMuted); err != nil {
		clog.Logger.Alertf("Failed to get session mute: %v", err)
	}

	return isMuted
}

func (s *masterSession) SetMute(mute bool) error {
	if err := s.volume.SetMute(mute, s.eventCtx); err != nil {
		return fmt.Errorf("set session mute: %w", err)
	}

	clog.Logger.Debugf("Setting session mute to %t", mute)

	return nil
}

func (s *masterSession) Release() {
	clog.Logger.Debug("Releasing audio session")

	s.volume.Release()
}

func (s *masterSession) String() string {
	return fmt.Sprintf(sessionStringFormat, s.humanReadableDesc, s.IsInput(), s.GetMute(), s.GetVolume())
}

func (s *masterSession) markAsStale() {
	s.stale = true
}
