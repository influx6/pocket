package guviews

import (
	"fmt"
	"sync"

	"github.com/go-humble/detect"
	"github.com/gopherjs/gopherjs/js"
	"github.com/influx6/faux/maker"
	"github.com/influx6/gu/gujs"
)

//==============================================================================

/*
	We want to have the ability to register and initialize views either on the
	backend or frontend in this manner.

	Registering a component should follow:

	guviews.Register("app.component/profile",func () guview.Renderable {
		...
	})

	guviews.Register("app.component/buttons",func () guviews.Renderables {
		...
	})

	Creating components should follow:

	guviews.Create({
		Name: "app.component/profile",
		Elem: "body",
		ID: "profile.453435",
		Paths: []string{"/app/profiles","/app"},
		Param: PConfig{...},
	},{
		Name: "app.component/buttons",
		Elem: "body",
		ID: "32-223-323232-453435",
		Paths: []string{"/app/selected","/app"},
		Param: BConfig{...},
	})

	This allows us to initialize the same view types with the appropriate
	Ids, paths and names without issuesing when moving from backend rendering
	to frontend rendering, we need the views to be able to take control on
	the frontend will updating their respective parts of the DOM without hassle
	when the backend switches control to the frontend in rendering.

*/

//==============================================================================

// ViewConfig provides a registery system for initializing views which have
// been registered within our makeregistery.
type ViewConfig struct {
	Name  string
	Elem  string
	ID    string
	Paths []string
	Param interface{}
}

//==============================================================================

var vox struct {
	maker  maker.Makers
	rw     sync.RWMutex
	builds map[string]Views
}

func init() {
	vox.maker = maker.New(nil)
	vox.builds = make(map[string]Views)
}

//==============================================================================

// Register adds the giving ViewMaker aliased with the giving name.
// Returns an error if the name is already registered. Ensure your names
// are unique.
func Register(name string, mux interface{}) error {
	return vox.maker.Register(name, mux)
}

// MustRegister works as the Register() function but panics on error.
func MustRegister(name string, mux interface{}) {
	if err := vox.maker.Register(name, mux); err != nil {
		panic(err)
	}
}

// Create intializes the giving views into the view instance registery
// provided the view Name does exist in the registery else return an
// error if any view was not found.
func Create(vcs ...ViewConfig) error {
	for _, vc := range vcs {
		res, err := vox.maker.Create(vc.Name, vc.Param)
		if err != nil {
			return err
		}

		var pass bool
		var view Views

		if render, ok := res.(Renderable); ok {

			// Creat a new view, set up the path directives and store this.
			view = NewWithID(vc.ID, render)
			pass = true

		}

		if renders, ok := res.(Renderables); ok {

			// Creat a new view, set up the path directives and store this.
			view = NewWithID(vc.ID, renders...)
			pass = true

		}

		if !pass {
			return fmt.Errorf("Invalid Type returned, Expected Renderable or Renderables: %+v", res)
		}

		// Attach the provided paths for this new view.
		for _, path := range vc.Paths {
			AttachView(view, path)
		}

		// Attach to the specified element as given.
		if detect.IsBrowser() {
			doc := js.Global.Get("document")
			if node := gujs.QuerySelector(doc, vc.Elem); node != nil {
				view.Mount(node)
			}
		}

		vox.rw.Lock()
		vox.builds[vc.ID] = view
		vox.rw.Unlock()
	}

	return nil
}

// MustCreate works as the Create() function but panics on error.
func MustCreate(vc ...ViewConfig) {
	if err := Create(vc...); err != nil {
		panic(err)
	}
}

// Get retrieves the registerd views with the giving ID else returns an error.
func Get(viewID string) (Views, error) {
	vox.rw.RLock()
	defer vox.rw.RUnlock()

	view, ok := vox.builds[viewID]
	if !ok {
		return nil, fmt.Errorf("No View found with ID[%q]", viewID)
	}

	return view, nil
}

// MustGet works like Get() function but panics when there is an error.
func MustGet(viewID string) Views {
	vox.rw.RLock()
	defer vox.rw.RUnlock()

	view, ok := vox.builds[viewID]
	if !ok {
		panic(fmt.Errorf("No View found with ID[%q]", viewID))
	}

	return view
}

//==============================================================================
