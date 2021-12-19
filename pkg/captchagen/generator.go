package captchagen

import (
	"bytes"
	"io"
	"log"
	"sync"
	"time"

	"github.com/dchest/captcha"
)

var customStore captcha.Store
var customStoreInitialized = false
var customStoreInitializedMtx = sync.Mutex{}

func Generate(id string) (io.Reader, error) {
	if !customStoreInitialized {
		initializeCustomStore()
	}

	// Create a state entry for the new image
	customStore.Set(id, captcha.RandomDigits(6))

	log.Println(customStore.Get(id, false))

	// Generate the image
	buf := bytes.Buffer{}
	err := captcha.WriteImage(&buf, id, 600, 400)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}

func Verify(id, test string) bool {
	return captcha.VerifyString(id, test)
}

func initializeCustomStore() {
	customStoreInitializedMtx.Lock()
	if !customStoreInitialized {
		customStore = captcha.NewMemoryStore(100, time.Hour)
		captcha.SetCustomStore(customStore)
		customStoreInitialized = true
	}

	customStoreInitializedMtx.Unlock()
}
