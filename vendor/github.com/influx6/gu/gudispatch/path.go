package gudispatch

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-humble/detect"
	"github.com/gopherjs/gopherjs/js"
	"github.com/influx6/faux/pattern"
)

//==============================================================================

// PathDirective represent the current path and hash values.
type PathDirective struct {
	Host     string
	Hash     string
	Path     string
	Sequence string
}

// String returns the hash and path.
func (p PathDirective) String() string {
	return fmt.Sprintf("%s%s", pattern.TrimEndSlashe(p.Path), p.Hash)
}

//==============================================================================

// PathSequencer provides a function to convert either the path/hash into a
// dot seperated sequence string for use with States.
type PathSequencer func(path string, hash string) string

// HashSequencer provides a PathSequencer that returns the hash part of a url,
// as the path sequence.
func HashSequencer(path, hash string) string {
	cleanHash := strings.Replace(hash, "#", ".", -1)
	return strings.Replace(cleanHash, "/", ".", -1)
}

// URLPathSequencer provides a PathSequencer that returns the path part of a url,
// as the path sequence.
func URLPathSequencer(path, hash string) string {
	return strings.Replace(path, "/", ".", -1)
}

//==============================================================================

//ErrNotSupported is returned when a feature requested is not supported by the environment
var ErrNotSupported = errors.New("Feature not supported")

// PathObserver represent any continouse changing route path by the browser
type PathObserver struct {
	usingHash bool
	sequencer PathSequencer
}

// NewPathObserver returns a new PathObserver instance.
func NewPathObserver(ps PathSequencer) *PathObserver {
	if ps == nil {
		ps = HashSequencer
	}

	return &PathObserver{
		sequencer: ps,
	}
}

//==============================================================================

// GetLocation returns the path and hash of the browsers location api else
// panics if not in a browser.
func GetLocation() (host string, path string, hash string) {
	if !detect.IsBrowser() {
		return
	}

	loc := js.Global.Get("location")
	host = loc.Get("host").String()
	path = loc.Get("pathname").String()
	hash = loc.Get("hash").String()
	return
}

//==============================================================================

// HashChangePath returns a path observer path changes
func HashChangePath(ps PathSequencer) *PathObserver {
	panicBrowserDetect()
	path := NewPathObserver(ps)
	path.usingHash = true

	js.Global.Set("onhashchange", func() {
		path.Follow(GetLocation())
	})

	return path
}

//==============================================================================

// BrowserSupportsPushState checks if browser supports pushState
func BrowserSupportsPushState() bool {
	if !detect.IsBrowser() {
		return false
	}

	return (js.Global.Get("onpopstate") != js.Undefined) &&
		(js.Global.Get("history") != js.Undefined) &&
		(js.Global.Get("history").Get("pushState") != js.Undefined)
}

//==============================================================================

// PopStatePath returns a path observer path changes
func PopStatePath(ps PathSequencer) (*PathObserver, error) {
	panicBrowserDetect()

	if !BrowserSupportsPushState() {
		return nil, ErrNotSupported
	}

	path := NewPathObserver(ps)

	js.Global.Set("onpopstate", func() {
		path.Follow(GetLocation())
	})

	return path, nil
}

// Follow creates a Pathspec from the hash and path and sends it
func (p *PathObserver) Follow(host, path, hash string) {
	fmt.Printf("Dispatch route change %s->%s\n", path, hash)
	Dispatch(PathDirective{
		Host:     host,
		Hash:     hash,
		Path:     path,
		Sequence: p.sequencer(path, hash),
	})
}

//==============================================================================

// PushDOMState adds a new state the dom push history.
func PushDOMState(path string) {
	if !detect.IsBrowser() {
		return
	}

	// Use the advance pushState feature.
	js.Global.Get("history").Call("pushState", nil, "", path)

	// Set the browsers hash accordinly.
	js.Global.Get("location").Set("hash", path)
}

// SetDOMHash sets the dom location hash.
func SetDOMHash(hash string) {
	panicBrowserDetect()
	js.Global.Get("location").Set("hash", hash)
}

//==============================================================================

// HistoryProvider wraps the PathObserver with methods that allow easy control of
// client location
type HistoryProvider struct {
	*PathObserver
}

// History returns a new PathObserver and depending on browser support will
// either use the popState or HashChange.
func History(ps PathSequencer) *HistoryProvider {
	pop, err := PopStatePath(ps)

	if err != nil {
		pop = HashChangePath(ps)
	}

	return &HistoryProvider{pop}
}

// Go changes the path of the current browser location depending on wether
// its underline observer is hashed based or pushState based,
// it will use SetDOMHash or PushDOMState appropriately.
func (h *HistoryProvider) Go(path string) {
	if h.usingHash {
		SetDOMHash(path)
		return
	}
	PushDOMState(path)
}

//==============================================================================

// panicBrowserDetect panics if the current gc is not a browser based
// one.
func panicBrowserDetect() {
	if !detect.IsBrowser() {
		panic("expected to be used in a dom/browser env")
	}
}
