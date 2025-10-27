// helpers/api.go
package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

const (
	apiBaseLegacy = "https://www.dnd5eapi.co/api"
	apiBase2014   = "https://www.dnd5eapi.co/api/2014"
)

var httpClient = &http.Client{
	Timeout: 12 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
	},
}

// -------- logging that callers can capture --------

type LogSink interface{ Add(line string) }
type logBuf struct{ rows []string }

func (b *logBuf) Add(s string)    { b.rows = append(b.rows, s) }
func NewLogBuf() *logBuf          { return &logBuf{} }
func (b *logBuf) Lines() []string { return b.rows }

// StdoutLog implements LogSink by writing lines to stdout. Use this from
// CLI callers to stream progress to the user.
type StdoutLog struct{}

func (s *StdoutLog) Add(line string) {
	fmt.Println(line)
}

// -------- HTTP helpers --------

func getJSON(ctx context.Context, u string, dst any, logs LogSink) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "DDsheetfinal/1.0")

	t0 := time.Now()
	resp, err := httpClient.Do(req)
	if err != nil {
		if logs != nil {
			logs.Add(fmt.Sprintf("GET %s -> transport error: %v", u, err))
		}
		return err
	}
	defer resp.Body.Close()
	dur := time.Since(t0)

	if resp.StatusCode != http.StatusOK {
		if logs != nil {
			logs.Add(fmt.Sprintf("GET %s -> %s (%v)", u, resp.Status, dur))
		}
		return fmt.Errorf("GET %s: %s", u, resp.Status)
	}
	if logs != nil {
		logs.Add(fmt.Sprintf("GET %s -> 200 (%v)", u, dur))
	}
	return json.NewDecoder(resp.Body).Decode(dst)
}

func getJSONWithFallback(ctx context.Context, path string, dst any, logs LogSink) error {
	start := time.Now()
	u1 := apiBaseLegacy + path
	if err := getJSON(ctx, u1, dst, logs); err == nil {
		if logs != nil {
			logs.Add(fmt.Sprintf("OK legacy %s (%v)", path, time.Since(start)))
		}
		return nil
	}
	if logs != nil {
		logs.Add("legacy miss -> trying 2014 for " + path)
	}
	u2 := apiBase2014 + path
	if err := getJSON(ctx, u2, dst, logs); err == nil {
		if logs != nil {
			logs.Add(fmt.Sprintf("OK 2014 %s (%v)", path, time.Since(start)))
		}
		return nil
	}
	return fmt.Errorf("both legacy and 2014 failed for %s", path)
}

// runner

func retryWithBackoff[T any](ctx context.Context, fn func(T) error, value T) error {
	backoff := 200 * time.Millisecond
	for attempt := 0; attempt < 3; attempt++ {
		if err := fn(value); err == nil {
			return nil
		}
		if attempt < 2 {
			t := time.NewTimer(backoff)
			select {
			case <-ctx.Done():
				t.Stop()
				return ctx.Err()
			case <-t.C:
				backoff *= 2
			}
		}
	}
	return nil
}

func runWorker[T any](ctx context.Context, jobs <-chan T, ticker *time.Ticker, errs chan<- error, fn func(T) error) {
	for j := range jobs {
		select {
		case <-ctx.Done():
			errs <- ctx.Err()
			return
		case <-ticker.C:
			if err := retryWithBackoff(ctx, fn, j); err != nil {
				errs <- err
				return
			}
		}
	}
}

func RunParallel[T any](ctx context.Context, items []T, perSec int, maxInFlight int, fn func(T) error) error {
	if len(items) == 0 {
		return nil
	}

	if perSec <= 0 {
		perSec = 6
	}
	if maxInFlight <= 0 {
		maxInFlight = 6
	}

	interval := time.Second / time.Duration(perSec)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	jobs := make(chan T)
	errs := make(chan error, len(items))

	var wg sync.WaitGroup
	wg.Add(maxInFlight)
	for w := 0; w < maxInFlight; w++ {
		go func() {
			defer wg.Done()
			runWorker(ctx, jobs, ticker, errs, fn)
		}()
	}

	go func() {
		defer close(jobs)
		for _, it := range items {
			select {
			case <-ctx.Done():
				return
			case jobs <- it:
			}
		}
	}()

	wg.Wait()
	close(errs)

	var anyErr error
	for e := range errs {
		if e != nil && anyErr == nil {
			anyErr = e
		}
	}
	return anyErr
}

// -------- indexing helpers (shared) --------

func sanitizeIndex(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	repl := []struct{ old, new string }{
		{"+", " plus "}, {"’", "'"}, {"–", "-"}, {"—", "-"}, {"/", "-"},
	}
	for _, r := range repl {
		s = strings.ReplaceAll(s, r.old, r.new)
	}
	reParen := regexp.MustCompile(`\s*\(.*?\)`)
	s = reParen.ReplaceAllString(s, "")
	s = strings.TrimSpace(s)
	s = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

func UniqueIndices(names []string, spell bool) map[string]string {
	m := map[string]string{}
	for _, n := range names {
		n = strings.TrimSpace(n)
		if n == "" {
			continue
		}
		idx := sanitizeIndex(n)
		if spell {
			switch strings.ToLower(n) {
			case "tasha's hideous laughter":
				idx = "tashas-hideous-laughter"
			}
		} else {
			switch strings.ToLower(n) {
			case "scale mail":
				idx = "scale-mail"
			case "leather armor":
				idx = "leather-armor"
			case "chain shirt":
				idx = "chain-shirt"
			case "half plate":
				idx = "half-plate"
			case "ring mail":
				idx = "ring-mail"
			case "chain mail":
				idx = "chain-mail"
			case "splint armor":
				idx = "splint"
			case "plate armor":
				idx = "plate"
			}
		}
		m[n] = idx
	}
	return m
}

// Small helper for building API paths
func SpellPath(idx string) string     { return "/spells/" + url.PathEscape(idx) }
func EquipmentPath(idx string) string { return "/equipment/" + url.PathEscape(idx) }
