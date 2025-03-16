// CookieStore API support
// Chromium ✅
// Edge ✅
// Firefox ❌
// Webkit ❌

package cookie

import (
	"os"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/playwright-community/playwright-go"
)

const URL_PATH = "http://localhost:3000/cookie/public"

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		panic(err)
	}

	code := m.Run()

	os.Exit(code)

}

func setup() error {
	return playwright.Install()
}

func TestWasmIsLoaded(t *testing.T) {
	pw, err := playwright.Run()
	if err != nil {
		t.Fatalf("error running playwright, err: %v", err)
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch()
	if err != nil {
		t.Fatalf("error launching chromium browser, err: %v", err)
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		t.Fatalf("error creating a new page, err: %v", err)
	}

	if _, err = page.Goto(URL_PATH); err != nil {
		t.Fatalf("couldn't go to %s, err: %v", URL_PATH, err)
	}

	logSink := make(chan string)
	page.OnConsole(func(cm playwright.ConsoleMessage) {
		logSink <- cm.String()
	})

	go func() {
		page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
			State: playwright.LoadStateLoad,
		})
		time.Sleep(time.Second * 2) // give extra 2 seconds of time to collect logs
		close(logSink)
	}()

	var acutalLogs []string
	expectedLogMsg := "wasm loaded"
	found := false

	for log := range logSink {
		acutalLogs = append(acutalLogs, log)
	}

	if slices.Contains(acutalLogs, expectedLogMsg) {
		found = true
		return
	}

	if !found {
		t.Errorf("wasm wasn't loaded, expected log to contain %s, got logs: %v", expectedLogMsg, acutalLogs)
	}
}

func TestCookieStoreApi(t *testing.T) {
	pw, err := playwright.Run()
	if err != nil {
		t.Fatalf("error running playwright, err: %v", err)
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch()
	if err != nil {
		t.Fatalf("error launching chromium browser, err: %v", err)
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		t.Fatalf("error creating a new page, err: %v", err)
	}

	if _, err = page.Goto(URL_PATH); err != nil {
		t.Fatalf("couldn't go to %s, err: %v", URL_PATH, err)
	}

	logSink := make(chan string)
	page.OnConsole(func(cm playwright.ConsoleMessage) {
		logSink <- cm.String()
	})

	go func() {
		page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
			State: playwright.LoadStateLoad,
		})
		time.Sleep(time.Second * 2) // give extra 2 seconds of time to collect logs
		close(logSink)
	}()

	var acutalLogs []string
	panicMessage := "panic"

	for log := range logSink {
		acutalLogs = append(acutalLogs, log)
	}

	for _, log := range acutalLogs {
		if strings.Contains(log, panicMessage) {
			t.Errorf("program panicked, got logs: %v", acutalLogs)
		}
	}
}
