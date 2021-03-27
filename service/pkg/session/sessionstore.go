package session

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
)

type sessionStore struct {
	configDir   string
	sessionFile string
	initialised bool
	c           *cache.Cache
	lock        sync.Mutex
}

var SessionCache sessionStore

func Init(configDir string) {
	if SessionCache.initialised {
		SessionCache.Destroy()
	}

	SessionCache = sessionStore{
		c:           cache.New(5*time.Minute, 10*time.Minute),
		configDir:   configDir + "/sessions",
		initialised: false,
	}

	SessionCache.init()
}

func (s *sessionStore) init() {
	if _, err := os.Stat(s.configDir); os.IsNotExist(err) {
		os.MkdirAll(s.configDir, 0777)
	}

	s.sessionFile = s.configDir + "/sessions.map"

	if _, err := os.Stat(s.sessionFile); !os.IsNotExist(err) {
		jsonFile, err := os.Open(s.sessionFile)
		// if we os.Open returns an error then handle it
		if err != nil {
			clog.Logger.Errorf("Error loading session file: %v", err)
			return
		}
		defer jsonFile.Close()
		byteValue, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			clog.Logger.Errorf("Error loading session file: %v", err)
			return
		}
		var sessions map[string]cache.Item
		json.Unmarshal(byteValue, &sessions)
		s.c = cache.NewFrom(5*time.Minute, 10*time.Minute, sessions)
	}
	s.initialised = true
}

func (s *sessionStore) save() {
	s.lock.Lock()
	defer s.lock.Unlock()
	sessions := s.c.Items()
	jsonString, err := json.Marshal(sessions)
	if err != nil {
		clog.Logger.Errorf("Error saving session file: %v", err)
		return
	}
	err = ioutil.WriteFile(s.sessionFile, jsonString, 0644)
	if err != nil {
		clog.Logger.Errorf("Error saving session file: %v", err)
		return
	}
}

func (s *sessionStore) StoreCommandData(profile string, action string, commandName string, value interface{}) {
	key := s.makeKey(profile, action, commandName)
	s.Add(key, value)
}

func (s *sessionStore) RetrieveCommandData(profile string, action string, commandName string) (interface{}, bool) {
	key := s.makeKey(profile, action, commandName)
	return s.Get(key)
}

func (s *sessionStore) makeKey(profile string, action string, commandName string) string {
	return fmt.Sprintf("%s_%s_%s", profile, action, commandName)
}

func (s *sessionStore) Add(key string, value interface{}) {
	err := s.c.Add(key, value, cache.NoExpiration)
	if err != nil {
		err = s.c.Replace(key, value, cache.NoExpiration)
		if err != nil {
			clog.Logger.Errorf("error saving value for %s: %v", key, err)
		}
	}
	go func() {
		s.save()
	}()
}

func (s *sessionStore) Get(key string) (interface{}, bool) {
	return s.c.Get(key)
}

func (s *sessionStore) Destroy() {
	s.save()
	s.initialised = false
	s.c.Flush()
}
