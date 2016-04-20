package gutrees

import "github.com/influx6/gu/guevents"

// Events provide an interface for markup event addition system
type Events interface {
	Events() []*Event
}

// EventHandler provides a custom event handler which allows access to the
// markup producing the event.
type EventHandler func(guevents.Event, Markup)

// Event provide a meta registry for helps in registering events for dom markups
// which is translated to the nodes themselves
type Event struct {
	Meta *guevents.EventMetable
	Fx   guevents.EventHandler
	tree Markup
}

// NewEvent returns a event object that allows registering events to eventlisteners
func NewEvent(etype, eselector string, efx EventHandler) *Event {
	ex := Event{
		Meta: &guevents.EventMetable{EventType: etype, EventTarget: eselector},
	}

	// wireup the function to get the ev and tree.
	ex.Fx = func(ev guevents.Event) {
		if efx != nil {
			efx(ev, ex.tree)
		}
	}

	return &ex
}

// StopImmediatePropagation will return itself and set StopPropagation to true
func (e *Event) StopImmediatePropagation() *Event {
	e.Meta.ShouldStopImmediatePropagation = true
	return e
}

// StopPropagation will return itself and set StopPropagation to true
func (e *Event) StopPropagation() *Event {
	e.Meta.ShouldStopPropagation = true
	return e
}

// PreventDefault will return itself and set PreventDefault to true
func (e *Event) PreventDefault() *Event {
	e.Meta.ShouldPreventDefault = true
	return e
}

// Events return the elements events
func (e *Element) Events() []*Event {
	return e.events
}

// Apply adds the event into the elements events lists
func (e *Event) Apply(ex Markup) {
	if em, ok := ex.(*Element); ok {
		if em.allowEvents {
			if e.Meta.EventTarget == "" {
				e.Meta.EventTarget = em.EventID()
			}
			e.tree = em
			em.events = append(em.events, e)
		}
	}
}

//Clone replicates the style into a unique instance
func (e *Event) Clone() *Event {
	return &Event{
		Meta: &guevents.EventMetable{EventType: e.Meta.EventType, EventTarget: e.Meta.EventTarget},
		Fx:   e.Fx,
	}
}
