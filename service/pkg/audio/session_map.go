package audio

import (
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/thoas/go-funk"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
)

type SessionMap struct {
	sessions           []Session
	lock               sync.Locker
	sessionFinder      SessionFinder
	lastSessionRefresh time.Time
}

var SessionMapInstance *SessionMap

func InitAudioSessions() error {
	sessionFinder, err := newSessionFinder()
	if err != nil {
		clog.Logger.Errorf("Failed to create SessionFinder: %v", err)
		return fmt.Errorf("create new SessionFinder: %w", err)
	}

	SessionMapInstance, err = newSessionMap(sessionFinder)
	if err != nil {
		clog.Logger.Errorf("Failed to create sessionMap: %v", err)
		return fmt.Errorf("create new sessionMap: %w", err)
	}

	if err := SessionMapInstance.initialize(); err != nil {
		return fmt.Errorf("init session map: %w", err)
	}

	return nil
}

func GetSessionNames() []string {
	names := make([]string, 0)
	for _, session := range SessionMapInstance.sessions {
		names = append(names, session.Key())
	}
	return names
}

func GetSession(key string) (Session, bool) {
	index := -1
	for x, session := range SessionMapInstance.sessions {
		if key == session.Key() {
			index = x
		}
	}
	if index < 0 {
		return nil, false
	}
	return SessionMapInstance.sessions[index], true
}

const (
	masterSessionName = "master" // master device volume
	systemSessionName = "system" // system sounds volume
	inputSessionName  = "mic"    // microphone input level

	// some targets need to be transformed before their correct audio sessions can be accessed.
	// this prefix identifies those targets to ensure they don't contradict with another similarly-named process
	specialTargetTransformPrefix = "deej."

	// targets the currently active window (Windows-only, experimental)
	specialTargetCurrentWindow = "current"

	// targets all currently unmapped sessions (experimental)
	specialTargetAllUnmapped = "unmapped"

	// this threshold constant assumes that re-acquiring all sessions is a kind of expensive operation,
	// and needs to be limited in some manner. this value was previously user-configurable through a config
	// key "process_refresh_frequency", but exposing this type of implementation detail seems wrong now
	minTimeBetweenSessionRefreshes = time.Second * 5

	// determines whether the map should be refreshed when a slider moves.
	// this is a bit greedy but allows us to ensure sessions are always re-acquired, which is
	// especially important for process groups (because you can have one ongoing session
	// always preventing lookup of other processes bound to its slider, which forces the user
	// to manually refresh sessions). a cleaner way to do this down the line is by registering to notifications
	// whenever a new session is added, but that's too hard to justify for how easy this solution is
	maxTimeBetweenSessionRefreshes = time.Second * 45
)

// this matches friendly device names (on Windows), e.g. "Headphones (Realtek Audio)"
var deviceSessionKeyPattern = regexp.MustCompile(`^.+ \(.+\)$`)

func newSessionMap(sessionFinder SessionFinder) (*SessionMap, error) {
	m := &SessionMap{
		sessions:      make([]Session, 0),
		lock:          &sync.Mutex{},
		sessionFinder: sessionFinder,
	}

	clog.Logger.Debug("Created session map instance")

	return m, nil
}

func (m *SessionMap) initialize() error {
	if err := m.getAndAddSessions(); err != nil {
		return fmt.Errorf("get all sessions during init: %w", err)
	}

	return nil
}

func (m *SessionMap) Release() error {
	if err := m.sessionFinder.Release(); err != nil {
		return fmt.Errorf("release session finder during release: %w", err)
	}

	return nil
}

// assumes the session map is clean!
// only call on a new session map or as part of refreshSessions which calls reset
func (m *SessionMap) getAndAddSessions() error {

	m.clear()

	// mark that we're refreshing before anything else
	m.lastSessionRefresh = time.Now()

	sessions, err := m.sessionFinder.GetAllSessions()
	if err != nil {
		return fmt.Errorf("get sessions from SessionFinder: %w", err)
	}

	m.sessions = sessions

	clog.Logger.Debugf("Got all audio sessions successfully\r\nsessionMap: %v", m)

	return nil
}

// performance: explain why force == true at every such use to avoid unintended forced refresh spams
func (m *SessionMap) RefreshSessions(force bool) {

	// make sure enough time passed since the last refresh, unless force is true in which case always clear
	if !force && m.lastSessionRefresh.Add(minTimeBetweenSessionRefreshes).After(time.Now()) {
		return
	}

	// clear and release sessions first
	m.clear()

	if err := m.getAndAddSessions(); err != nil {
		clog.Logger.Alertf("Failed to re-acquire all audio sessions: %v", err)
	} else {
		clog.Logger.Debug("Re-acquired sessions successfully")
	}
}

// returns true if a session is not currently mapped to any slider, false otherwise
// special sessions (master, system, mic) and device-specific sessions always count as mapped,
// even when absent from the config. this makes sense for every current feature that uses "unmapped sessions"
func (m *SessionMap) sessionMapped(session Session) bool {

	// count master/system/mic as mapped
	if funk.ContainsString([]string{masterSessionName, systemSessionName, inputSessionName}, session.Key()) {
		return true
	}

	// count device sessions as mapped
	if deviceSessionKeyPattern.MatchString(session.Key()) {
		return true
	}

	matchFound := false

	/*
		// look through the actual mappings
		m.deej.config.SliderMapping.iterate(func(sliderIdx int, targets []string) {
			for _, target := range targets {

				// ignore special transforms
				if m.targetHasSpecialTransform(target) {
					continue
				}

				// safe to assume this has a single element because we made sure there's no special transform
				target = m.resolveTarget(target)[0]

				if target == session.Key() {
					matchFound = true
					return
				}
			}
		})
	*/

	return matchFound
}

/*
func (m *sessionMap) handleSliderMoveEvent(event SliderMoveEvent) {

		// first of all, ensure our session map isn't moldy
		if m.lastSessionRefresh.Add(maxTimeBetweenSessionRefreshes).Before(time.Now()) {
			clog.Logger.Debug("Stale session map detected on slider move, refreshing")
			m.refreshSessions(true)
		}

		// get the targets mapped to this slider from the config
		targets, ok := m.deej.config.SliderMapping.get(event.SliderID)

		// if slider not found in config, silently ignore
		if !ok {
			return
		}

		targetFound := false
		adjustmentFailed := false

		// for each possible target for this slider...
		for _, target := range targets {

			// resolve the target name by cleaning it up and applying any special transformations.
			// depending on the transformation applied, this can result in more than one target name
			resolvedTargets := m.resolveTarget(target)

			// for each resolved target...
			for _, resolvedTarget := range resolvedTargets {

				// check the map for matching sessions
				sessions, ok := m.get(resolvedTarget)

				// no sessions matching this target - move on
				if !ok {
					continue
				}

				targetFound = true

				// iterate all matching sessions and adjust the volume of each one
				for _, session := range sessions {
					if session.GetVolume() != event.PercentValue {
						if err := session.SetVolume(event.PercentValue); err != nil {
							clog.Logger.Alertf("Failed to set target session volume: %v", err)
							adjustmentFailed = true
						}
					}
					if session.GetMute() != event.MuteValue {
						if err := session.SetMute(event.MuteValue); err != nil {
							clog.Logger.Alertf("Failed to set target session mute: %v", err)
							adjustmentFailed = true
						}
					}
				}
			}
		}

		// if we still haven't found a target or the volume adjustment failed, maybe look for the target again.
		// processes could've opened since the last time this slider moved.
		// if they haven't, the cooldown will take care to not spam it up
		if !targetFound {
			m.refreshSessions(false)
		} else if adjustmentFailed {

			// performance: the reason that forcing a refresh here is okay is that we'll only get here
			// when a session's SetVolume call errored, such as in the case of a stale master session
			// (or another, more catastrophic failure happens)
			m.refreshSessions(true)
		}
	}
*/

func (m *SessionMap) get(key string) (Session, bool) {
	m.lock.Lock()
	defer m.lock.Unlock()
	index := -1
	for x, session := range m.sessions {
		if session.Key() == key {
			index = x
		}
	}
	if index >= 0 {
		return m.sessions[index], true
	}
	return nil, false
}

func (m *SessionMap) clear() {
	m.lock.Lock()
	defer m.lock.Unlock()

	clog.Logger.Debug("Releasing and clearing all audio sessions")

	for _, session := range m.sessions {
		session.Release()
	}

	clog.Logger.Debug("Session map cleared")
}

func (m *SessionMap) String() string {
	m.lock.Lock()
	defer m.lock.Unlock()

	return fmt.Sprintf("<%d audio sessions>", len(m.sessions))
}

func (m *SessionMap) PrintSessionNames() {
	for _, session := range m.sessions {
		clog.Logger.Infof("Session sessionKey: %s, desc: %s", session.Key(), session.String())
	}
}
