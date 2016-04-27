package gutrees

import (
	"fmt"
	"strings"
)

//==============================================================================

// Attributes interface defines a type that has Attributes
type Attributes interface {
	Attributes() []*Attribute
}

// Attribute define the struct  for attributes
type Attribute struct {
	Name  string
	Value string
}

// NewAttr returns a new attribute instance
func NewAttr(name, val string) *Attribute {
	a := Attribute{Name: name, Value: val}
	return &a
}

// Apply applies a set change to the giving element attributes list
func (a *Attribute) Apply(e Markup) {
	if em, ok := e.(*Element); ok {
		if em.allowAttributes {
			em.attrs = append(em.attrs, a)
		}
	}
}

//Clone replicates the attribute into a unique instance
func (a *Attribute) Clone() *Attribute {
	return &Attribute{Name: a.Name, Value: a.Value}
}

// Reconcile checks if the attribute matches then upgrades its value.
func (a *Attribute) Reconcile(m *Attribute) bool {
	if strings.TrimSpace(a.Name) == strings.TrimSpace(m.Name) {
		a.Value = m.Value
		return true
	}
	return false
}

//==============================================================================

// Styles interface defines a type that has Styles
type Styles interface {
	Styles() []*Style
}

// Style define the style specification for element styles
type Style struct {
	Name  string
	Value string
}

// NewStyle returns a new style instance
func NewStyle(name, val string) *Style {
	s := Style{Name: name, Value: val}
	return &s
}

//Clone replicates the style into a unique instance
func (s *Style) Clone() *Style {
	return &Style{Name: s.Name, Value: s.Value}
}

// Apply applies a set change to the giving element style list
func (s *Style) Apply(e Markup) {
	if em, ok := e.(*Element); ok {
		if em.allowStyles {
			em.styles = append(em.styles, s)
		}
	}
}

// Reconcile checks if the style matches then upgrades its value.
func (s *Style) Reconcile(m *Style) bool {
	if strings.TrimSpace(s.Name) == strings.TrimSpace(m.Name) {
		s.Value = m.Value
		return true
	}
	return false
}

//==============================================================================

// ClassList defines the list type for class lists.
type ClassList []string

// Add adds a class name into the lists.
func (c *ClassList) Add(class string) {
	*c = append(*c, class)
}

// Apply checks for a class attribute
func (c *ClassList) Apply(em Markup) {
	if len(*c) == 0 {
		return
	}

	e, ok := em.(*Element)
	if !ok {
		return
	}

	list := strings.Join(*c, " ")

	a, err := GetAttr(e, "class")

	if err != nil {
		(&Attribute{Name: "class", Value: "list"}).Apply(e)
		return
	}

	a.Value = fmt.Sprintf("%s %s", a.Value, list)
}

// Clone replicates the lists of classnames.
func (c *ClassList) Clone() *ClassList {
	newlist := new(ClassList)
	*newlist = append(*newlist, (*c)...)
	return newlist
}

// Reconcile checks each item against the given lists
// and replaces/add any missing item.
func (c *ClassList) Reconcile(m *ClassList) bool {
	var added bool

	maxlen := len(*c)

	for ind, val := range *c {

		if ind >= maxlen {
			added = true
			c.Add(val)
			continue
		}

		if (*c)[ind] == val {
			continue
		}

		added = true
		(*c)[ind] = val
	}

	return added
}
