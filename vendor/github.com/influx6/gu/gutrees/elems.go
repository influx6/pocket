package gutrees

import (
	"fmt"
	"strings"

	"github.com/influx6/gu/guevents"
)

// Markup provide a basic specification type of how a element resolves its content
type Markup interface {
	Identity
	MarkupChildren
	Appliable
	Reconcilable
	Clonable
	Removable
}

// Element represent a concrete implementation of a element node
type Element struct {
	removed         bool
	autoclose       bool
	allowEvents     bool
	allowChildren   bool
	allowStyles     bool
	allowAttributes bool
	uid             string
	hash            string
	tagname         string
	textContent     string
	events          []*Event
	styles          []*Style
	attrs           []*Attribute
	children        []Markup
	eventManager    guevents.EventManagers
}

// NewText returns a new Text instance element
func NewText(txt string) *Element {
	em := NewElement("text", false)
	em.allowChildren = false
	em.allowAttributes = false
	em.allowStyles = false
	em.allowEvents = false
	em.textContent = txt
	return em
}

// NewElement returns a new element instance giving the specificed name
func NewElement(tag string, hasNoEndingTag bool) *Element {
	return &Element{
		uid:             RandString(8),
		hash:            RandString(10),
		tagname:         strings.ToLower(strings.TrimSpace(tag)),
		children:        make([]Markup, 0),
		styles:          make([]*Style, 0),
		attrs:           make([]*Attribute, 0),
		autoclose:       hasNoEndingTag,
		allowChildren:   true,
		allowStyles:     true,
		allowAttributes: true,
		allowEvents:     true,
	}
}

// AutoClosed returns true/false if this element uses a </> or a <></> tag convention
func (e *Element) AutoClosed() bool {
	return e.autoclose
}

//==============================================================================

// Eventers provide an interface type for elements able to register and load
// event managers.
type Eventers interface {
	LoadEvents()
	UseEventManager(guevents.EventManagers) bool
}

// UseEventManager adds a eventmanager into the markup and if not available before automatically registers
// the events with it,once an event manager is registered to it,it will and can not be changed
func (e *Element) UseEventManager(man guevents.EventManagers) bool {
	if man == nil {
		return true
	}

	if e.eventManager != nil {
		// e.eventManager.
		man.AttachManager(e.eventManager)
		return false
	}

	e.eventManager = man
	e.LoadEvents()
	return true
}

// LoadEvents loads up the events registered by this and by its children into each respective
// available events managers
func (e *Element) LoadEvents() {
	if e.eventManager != nil {
		e.eventManager.DisconnectRemoved()

		for _, ev := range e.events {
			if es, _ := e.eventManager.NewEventMeta(ev.Meta); es != nil {
				es.Q(ev.Fx)
			}
		}

	}

	//load up the children events also
	for _, em := range e.children {
		if ech, ok := em.(ElementalMarkup); ok {
			if !ech.UseEventManager(e.eventManager) {
				ech.LoadEvents()
			}
		}
	}
}

//==============================================================================

// EventID returns the selector used for tagging events for a markup.
func (e *Element) EventID() string {
	return fmt.Sprintf("%s[uid='%s']", strings.ToLower(e.Name()), e.UID())
}

// Empty resets the elements children list as 0 length
func (e *Element) Empty() {
	e.children = e.children[:0]
}

//==============================================================================

// Identity defines an interface for identifiable structures.
type Identity interface {
	Name() string
	UID() string
	Hash() string
}

// Name returns the tag name of the element
func (e *Element) Name() string {
	return e.tagname
}

// UID returns the current uid of the Element
func (e *Element) UID() string {
	return e.uid
}

// Hash returns the current hash of the Element
func (e *Element) Hash() string {
	return e.hash
}

//==============================================================================

// TextMarkup defines a interface for text based markup.
type TextMarkup interface {
	TextContent() string
}

// TextContent returns the elements text value if its a text
// type else an empty string.
func (e *Element) TextContent() string {
	return e.textContent
}

//==============================================================================

// Cleanable defines a interface for structures to self sanitize their contents.
type Cleanable interface {
	Clean()
}

// Clean cleans out all internal markup marked as removable.
func (e *Element) Clean() {
	for n, elm := range e.children {
		if elm.Removed() {
			copy(e.children[n:], e.children[n+1:])
			e.children = e.children[:len(e.children)-1]
		} else {
			if em, ok := elm.(Cleanable); ok {
				em.Clean()
			}
		}
	}
}

//==============================================================================

// Removable defines a self removal type structure.
type Removable interface {
	Remove()
	Removed() bool
}

// Remove sets the markup as removable and adds a 'haikuRemoved' attribute to it
func (e *Element) Remove() {
	if !e.Removed() {
		e.attrs = append(e.attrs, &Attribute{"haikuRemoved", ""})
		e.removed = true
	}
}

// Removed returns true/false if the Element is marked removed
func (e *Element) Removed() bool {
	return !!e.removed
}

//==============================================================================

// SwappableIdentity defines an interface that allows swapping a structures
// identity information.
type SwappableIdentity interface {
	SwapHash(string)
	SwapUID(string)
	UpdateHash()
}

// SwapUID swaps the uid of the internal Element.
func (e *Element) SwapUID(uid string) {
	e.uid = uid
}

// SwapHash swaps the hash of the internal Element.
func (e *Element) SwapHash(hash string) {
	e.hash = hash
}

// UpdateHash updates the Element hash value
func (e *Element) UpdateHash() {
	e.hash = RandString(10)
}

//==============================================================================

// ElementalMarkup defines a markup for elemental structures, which provide the
// concrete structure for dom nodes.
type ElementalMarkup interface {
	Markup
	Events
	Styles
	Attributes
	Eventers
	SwappableIdentity
	TextMarkup
	Cleanable

	AutoClosed() bool

	EventID() string
	Empty()
}

// Reconcilable defines the interface of markups that can reconcile their content against another
type Reconcilable interface {
	Reconcile(Markup) bool
}

// Reconcile takes a old markup and reconciles its uid and its children with
// these information,it returns a true/false telling the parent if the children
// swapped hashes.
// The reconcilation uses the order in which elements are added, if the order
// and element types are same then the uid are swapped, else it firsts checks the
// element type, but if not the same adds the old one into the new list as removed
// then continues the check. The system takes position of elements in the old and
// new as very important and I cant stress this enough, "Element Positioning" in
// the markup are very important, If a Anchor was the first element in the old
// render and the next pass returns a Div in the position for that Anchor in the
// new render, the old Anchor will be marked as removed and will be removed from
// the dom and ignored by the writers.
// When two elements position are same and their types are the same then a checkup
// process is doing using the elements attributes, this is done to determine if the
// hash value of the new should be swapped with the old. We cant use style properties
// here because they are the most volatile of the set and will periodically be
// either changed and returned to normal values eg display: none to display: block
// and vise-versa, so only attributes are used in the check process.
func (e *Element) Reconcile(m Markup) bool {
	em, ok := m.(ElementalMarkup)
	if !ok {
		return false
	}

	// are we reconciling the proper elements type ? if not skip (i.e different types cant reconcile eachother)]
	// TODO: decide if we should mark the markup as removed in this case as a catchall system
	if e.Name() != em.Name() {
		return false
	}

	em.Clean()

	//since the tagname are the same, swap uids
	// olduid := em.UID()
	e.SwapUID(em.UID())

	//since the tagname are the same and we have swapped uid, to determine who gets or keeps
	// its hash we will check the attributes against each other, but also the hash is dependent on the
	// children also, if the children observered there was a change
	oldHash := em.Hash()
	// newHash := e.Hash()

	// if we have a special case for text element then we do things differently
	if e.Name() == "text" {
		//if the contents are equal,keep the prev hash
		if e.TextContent() == em.TextContent() {
			e.SwapHash(oldHash)
			return false
		}
		return true
	}

	newChildren := e.Children()
	oldChildren := em.Children()
	maxSize := len(newChildren)
	oldMaxSize := len(oldChildren)

	// if the element had no children too, swap hash.
	if maxSize <= 0 {
		if oldMaxSize <= 0 {
			e.SwapHash(oldHash)
			return false
		}

		return true
	}

	var childChanged bool

	for n, och := range oldChildren {
		if maxSize > n {

			nch := newChildren[n]

			if nch.Name() == och.Name() {

				if nch.Reconcile(och) {
					childChanged = true
				}

			} else {

				och.Remove()
				e.AddChild(och)
			}

			continue
		}

		och.Remove()
		e.AddChild(och)
	}

	ReconcileEvents(e, em)

	if e.eventManager != nil {
		e.eventManager.DisconnectRemoved()
	}

	//if the sizes of the new node is more than the old node then ,we definitely changed
	if maxSize > oldMaxSize {
		return true
	}

	if !childChanged && EqualAttributes(e, em) && EqualStyles(e, em) {
		e.SwapHash(oldHash)
		return false
	}

	return true
}

//==============================================================================

// MarkupChildren defines the interface of an element that has children
type MarkupChildren interface {
	AddChild(...Markup)
	Children() []Markup
}

// AddChild adds a new markup as the children of this element
func (e *Element) AddChild(em ...Markup) {
	if e.allowChildren {
		for _, mm := range em {

			if mm == nil {
				continue
			}

			if m, ok := mm.(ElementalMarkup); ok {
				e.children = append(e.children, m)
				//if this are free elements, then use this event manager
				m.UseEventManager(e.eventManager)
			}

		}
	}
}

// Children returns the children list for the element
func (e *Element) Children() []Markup {
	return e.children
}

//==============================================================================

// Styles return the internal style list of the element
func (e *Element) Styles() []*Style {
	return e.styles
}

// Attributes return the internal attribute list of the element
func (e *Element) Attributes() []*Attribute {
	return e.attrs
}

//==============================================================================

// Appliable define the interface specification for applying changes to elements elements in tree
type Appliable interface {
	Apply(Markup)
}

//Apply adds the giving element into the current elements children tree
func (e *Element) Apply(em Markup) {
	if mm, ok := em.(MarkupChildren); ok {
		mm.AddChild(e)
	}
}

//==============================================================================

// Clonable defines an interface for objects that can be cloned
type Clonable interface {
	Clone() Markup
}

// Clone makes a new copy of the markup structure
func (e *Element) Clone() Markup {
	co := NewElement(e.Name(), e.AutoClosed())

	//copy over the textContent
	co.textContent = e.textContent

	//copy over the attribute lockers
	co.allowChildren = e.allowChildren
	co.allowEvents = e.allowEvents
	co.allowAttributes = e.allowAttributes
	co.eventManager = e.eventManager

	if e.Removed() {
		co.Removed()
	}

	//clone the internal styles
	for _, so := range e.styles {
		so.Clone().Apply(co)
	}

	co.allowStyles = e.allowStyles

	//clone the internal attribute
	for _, ao := range e.attrs {
		ao.Clone().Apply(co)
	}

	// co.allowAttributes = e.allowAttributes
	//clone the internal children
	for _, ch := range e.children {
		ch.Clone().Apply(co)
	}

	for _, ch := range e.events {
		ch.Clone().Apply(co)
	}

	return co
}
