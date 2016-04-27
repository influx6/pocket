package gudispatch

import (
	"sync"

	"github.com/influx6/faux/pattern"
)

// matchers define a lists of pattern associated with url match validators.
var matchers = make(map[string]pattern.URIMatcher)
var matchlock sync.RWMutex

// URIMatcher returns a new uri matcher if it has not being already creatd.
func URIMatcher(path string) pattern.URIMatcher {
	matchlock.RLock()
	mk, ok := matchers[path]
	matchlock.RUnlock()

	if !ok {
		m := pattern.New(path)
		matchlock.Lock()
		matchers[path] = m
		matchlock.Unlock()
		return m
	}

	return mk
}
