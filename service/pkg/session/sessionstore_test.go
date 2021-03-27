package session

import (
	"os"
	"testing"
)

func TestStore(t *testing.T) {

	Init("e:/temp/ReCoS")

	defer SessionCache.Destroy()

	if !SessionCache.initialised {
		t.Error("SessionCache not initialised")
	}

	SessionCache.Add("key", "Value")
	SessionCache.save()
	if _, err := os.Stat(SessionCache.sessionFile); os.IsNotExist(err) {
		t.Error("File not written")
	}

	value, ok := SessionCache.Get("key")
	if value != "Value" {
		t.Error("value not equal")
	}
	if !ok {
		t.Error("value not found")
	}

	SessionCache.Destroy()

	value, ok = SessionCache.Get("key")
	if ok {
		t.Error("value found")
	}

	Init("e:/temp/ReCoS")

	value, ok = SessionCache.Get("key")
	if value != "Value" {
		t.Error("value not equal")
	}
	if !ok {
		t.Error("value not found")
	}

	SessionCache.Destroy()

}
