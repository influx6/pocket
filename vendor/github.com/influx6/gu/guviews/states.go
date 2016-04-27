package guviews

import (
	"strings"

	"github.com/influx6/gu/gutrees"
	"github.com/influx6/gu/gutrees/styles"
)

// ViewStates defines the two possible behavioral state of a view's markup
type ViewStates interface {
	Render(gutrees.Markup)
}

// HideView provides a ViewStates for Views inactive state
type HideView struct{}

// Render marks the given markup as display:none
func (v HideView) Render(m gutrees.Markup) {
	// if we are allowed to query for styles then check and change display
	if mm, ok := m.(gutrees.MarkupProps); ok {
		if ds, err := gutrees.GetStyle(mm, "display"); err == nil {
			if !strings.Contains(ds.Value, "none") {
				ds.Value = "none"
			}
			return
		}

	}
	styles.Display("none").Apply(m)
}

// ShowView provides a ViewStates for Views active state
type ShowView struct{}

// Render marks the given markup with a display: block
func (v ShowView) Render(m gutrees.Markup) {
	//if we are allowed to query for styles then check and change display
	if mm, ok := m.(gutrees.MarkupProps); ok {
		if ds, err := gutrees.GetStyle(mm, "display"); err == nil {
			if strings.Contains(ds.Value, "none") {
				ds.Value = "block"
			}
			return
		}

		styles.Display("block").Apply(m)
	}
}

//==============================================================================
