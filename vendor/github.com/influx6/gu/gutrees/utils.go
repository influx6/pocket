package gutrees

import (
	"crypto/rand"
	"strings"
)

// RandString generates a set of random numbers of a set length
func RandString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

// Augment adds new markup to an the root if its Element
func Augment(root Markup, m ...Markup) {
	if el, ok := root.(*Element); ok {
		for _, mo := range m {
			mo.Apply(el)
		}
	}
}

// ReconcileEvents checks through two markup events against each other and if it finds any disparity marks
// event objects as Removed
func ReconcileEvents(e, em Events) {
	oldevents := em.Events()
	newevents := e.Events()

	if len(oldevents) <= 0 && len(newevents) <= 0 {
		return
	}

	if len(newevents) <= 0 && len(oldevents) > 0 {
		for _, ev := range oldevents {
			ev.Meta.Remove()
		}
		return
	}

	checkOut := func(ev *Event) bool {
		for _, evs := range newevents {
			if evs.Meta.EventType == ev.Meta.EventType {
				return true
			}
		}
		return false

	}
	//set to equal as the logic will try to assert its falsiness

	// outerfind:
	for _, ev := range oldevents {
		if checkOut(ev) {
			continue
		}

		ev.Meta.Remove()
	}

}

// EqualStyles returns true/false if the style values are all equal attribute.
func EqualStyles(e, em Styles) bool {
	oldAttrs := em.Styles()

	if len(oldAttrs) <= 0 {
		if len(e.Styles()) <= 0 {
			return true
		}
		return false
	}

	//set to equal as the logic will try to assert its falsiness
	var equal = true

	for _, oa := range oldAttrs {
		//lets get the styles type from the element, if it exists then check the value if its equal
		// continue the loop and check the rest, else we found a contention point, attribute of old markup
		// does not exists in new markup, so we break and mark as different,letting the new markup keep its hash
		// but if the loop finishes and all are equal then we swap the hashes
		if ta, err := GetStyle(e, oa.Name); err == nil {
			if ta.Value == oa.Value {
				continue
			}

			equal = false
			break
		} else {
			equal = false
			break
		}
	}

	return equal
}

// EqualAttributes returns true/false if the elements and the giving markup have equal attribute
func EqualAttributes(e, em Attributes) bool {
	oldAttrs := em.Attributes()

	if len(oldAttrs) <= 0 {
		if len(e.Attributes()) <= 0 {
			return true
		}
		return false
	}

	//set to equal as the logic will try to assert its falsiness
	var equal = true

	for _, oa := range oldAttrs {
		//lets get the attribute type from the element, if it exists then check the value if its equal
		// continue the loop and check the rest, else we found a contention point, attribute of old markup
		// does not exists in new markup, so we break and mark as different,letting the new markup keep its hash
		// but if the loop finishes and all are equal then we swap the hashes
		if ta, err := GetAttr(e, oa.Name); err == nil {
			if ta.Value == oa.Value {
				continue
			}

			equal = false
			break
		} else {
			equal = false
			break
		}
	}

	return equal
}

// GetStyles returns the styles that contain the specified name and if not empty that contains the specified value also, note that strings
// NOTE: string.Contains is used when checking value parameter if present
func GetStyles(e Styles, f, val string) []*Style {
	var found []*Style
	var styles = e.Styles()

	for _, as := range styles {
		if as.Name != f {
			continue
		}

		if val != "" {
			if !strings.Contains(as.Value, val) {
				continue
			}
		}

		found = append(found, as)
	}

	return found
}

// GetStyle returns the style with the specified tag name
func GetStyle(e Styles, f string) (*Style, error) {
	styles := e.Styles()
	for _, as := range styles {
		if as.Name == f {
			return as, nil
		}
	}
	return nil, ErrNotFound
}

// StyleContains returns the styles that contain the specified name and if the val is not empty then
// that contains the specified value also, note that strings
// NOTE: string.Contains is used
func StyleContains(e Styles, f, val string) bool {
	styles := e.Styles()
	for _, as := range styles {
		if !strings.Contains(as.Name, f) {
			continue
		}

		if val != "" {
			if !strings.Contains(as.Value, val) {
				continue
			}
		}

		return true
	}

	return false
}

// GetAttrs returns the attributes that have the specified text within the naming
// convention and if it also contains the set val if not an empty "",
// NOTE: string.Contains is used
func GetAttrs(e Attributes, f, val string) []*Attribute {

	var found []*Attribute

	for _, as := range e.Attributes() {
		if as.Name != f {
			continue
		}

		if val != "" {
			if !strings.Contains(as.Value, val) {
				continue
			}
		}

		found = append(found, as)
	}

	return found
}

// AttrContains returns the attributes that have the specified text within the naming
// convention and if it also contains the set val if not an empty "",
// NOTE: string.Contains is used
func AttrContains(e Attributes, f, val string) bool {
	for _, as := range e.Attributes() {
		if !strings.Contains(as.Name, f) {
			continue
		}

		if val != "" {
			if !strings.Contains(as.Value, val) {
				continue
			}
		}

		return true
	}

	return false
}

// GetAttr returns the attribute with the specified tag name
func GetAttr(e Attributes, f string) (*Attribute, error) {
	for _, as := range e.Attributes() {
		if as.Name == f {
			return as, nil
		}
	}
	return nil, ErrNotFound
}

//==============================================================================

// MarkupProps defines a custom type that combines the Markup, Styles and
// Attributes interfaces.
type MarkupProps interface {
	Markup
	Styles
	Attributes
}

// ElementsUsingStyle returns the children within the element matching the
// stlye restrictions passed.
// NOTE: is uses StyleContains
func ElementsUsingStyle(e MarkupProps, f, val string) []Markup {
	return DeepElementsUsingStyle(e, f, val, 1)
}

// ElementsWithAttr returns the children within the element matching the
// stlye restrictions passed.
// NOTE: is uses AttrContains
func ElementsWithAttr(e MarkupProps, f, val string) []Markup {
	return DeepElementsWithAttr(e, f, val, 1)
}

// DeepElementsUsingStyle returns the children within the element matching the
// style restrictions passed allowing control of search depth
// NOTE: is uses StyleContains
// WARNING: depth must start at 1
func DeepElementsUsingStyle(e MarkupProps, f, val string, depth int) []Markup {
	if depth <= 0 {
		return nil
	}

	var found []Markup

	for _, c := range e.Children() {
		if ch, ok := c.(MarkupProps); ok {
			if StyleContains(ch, f, val) {
				found = append(found, ch)
				cfo := DeepElementsUsingStyle(ch, f, val, depth-1)
				if len(cfo) > 0 {
					found = append(found, cfo...)
				}
			}
		}
	}

	return found
}

// DeepElementsWithAttr returns the children within the element matching the
// attributes restrictions passed allowing control of search depth
// NOTE: is uses Element.AttrContains
// WARNING: depth must start at 1
func DeepElementsWithAttr(e MarkupProps, f, val string, depth int) []Markup {
	if depth <= 0 {
		return nil
	}

	var found []Markup

	for _, c := range e.Children() {
		if ch, ok := c.(MarkupProps); ok {
			if AttrContains(ch, f, val) {
				found = append(found, ch)
				cfo := DeepElementsWithAttr(ch, f, val, depth-1)
				if len(cfo) > 0 {
					found = append(found, cfo...)
				}
			}
		}
	}

	return found
}

// ElementsWithTag returns elements matching the tag type in the parent markup children list
// only without going deeper into children's children lists
func ElementsWithTag(e MarkupProps, f string) []Markup {
	return DeepElementsWithTag(e, f, 1)
}

// DeepElementsWithTag returns elements matching the tag type in the parent markup
// and depending on the depth will walk down other children within the children.
// WARNING: depth must start at 1
func DeepElementsWithTag(e MarkupProps, f string, depth int) []Markup {
	if depth <= 0 {
		return nil
	}

	f = strings.TrimSpace(strings.ToLower(f))

	var found []Markup

	for _, c := range e.Children() {
		if ch, ok := c.(MarkupProps); ok {
			if ch.Name() == f {
				found = append(found, ch)
				cfo := DeepElementsWithTag(ch, f, depth-1)
				if len(cfo) > 0 {
					found = append(found, cfo...)
				}
			}
		}
	}

	return found
}

//==============================================================================
