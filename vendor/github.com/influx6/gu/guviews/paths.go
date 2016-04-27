package guviews

import (
	"sync"

	"github.com/influx6/gu/gudispatch"
)

//==============================================================================

// Path defines a representation of a location path matching a specific sequence.
type Path struct {
	gudispatch.PathDirective
	Param map[string]string
	ID    string
}

//==============================================================================

// pathMl provides a mutex to control read and write to internal view cache store.
var pathMl sync.RWMutex

// pathMl2 provides a mutex to control write and read to internal cache maps.
var pathMl2 sync.RWMutex

// pathWatch registers a view with selected watch routes to reduce unnecessary
// multiwatching of same routes and helps manage state of views.
var pathWatch = make(map[Views]map[string]bool)

// AttachView allows setting a specific view to become active for a specific URL
// route pattern. This allows to control the active and inactive state and also
// the visibility of the view dependent on the current location URL path.
func AttachView(v Views, pattern string) {
	gudispatch.URIMatcher(pattern)

	new := addView(v)

	pathMl.RLock()
	vcache := pathWatch[v]
	pathMl.RUnlock()

	// If we are already watching for this specific route then skip this.
	pathMl.RLock()
	hasOk := vcache[pattern]
	pathMl.RUnlock()

	if hasOk {
		return
	}

	pathMl2.Lock()
	vcache[pattern] = true
	pathMl2.Unlock()

	if !new {
		return
	}

	gudispatch.Subscribe(func(p gudispatch.PathDirective) {
		pathMl2.RLock()
		defer pathMl2.RUnlock()

		var found bool
		var params map[string]string

		for key := range vcache {
			// Get the matcher for this key.
			watcher := gudispatch.URIMatcher(key)

			pm, ok := watcher.Validate(p.String())
			if !ok {
				continue
			}

			params = pm
			found = true
			break
		}

		if !found {
			v.Hide()
			return
		}

		pu := Path{
			PathDirective: p,
			ID:            v.UUID(),
			Param:         params,
		}

		gudispatch.Dispatch(&pu)
	})

	gudispatch.Follow(gudispatch.GetLocation())
}

// addView attaches view into the pathWatch match.
func addView(v Views) bool {
	// Get the view route cache.
	pathMl.RLock()
	_, ok := pathWatch[v]
	pathMl.RUnlock()

	// If no cache is found for this view then make one and store it.
	if !ok {
		vcache := make(map[string]bool)
		pathMl.Lock()
		pathWatch[v] = vcache
		pathMl.Unlock()
		return true
	}

	return false
}

//==============================================================================
