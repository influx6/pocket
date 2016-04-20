package styles

import (
	"strconv"

	"github.com/influx6/gu/gutrees"
)

// Size presents a basic stringifed unit
type Size string

// Px returns the value in "%px" format
func Px(pixels int) Size {
	return Size(strconv.Itoa(pixels) + "px")
}

// Color provides the color style value
func Color(value string) *gutrees.Style {
	return &gutrees.Style{Name: "color", Value: value}
}

// Display provides the style setter that sets the css display value.
func Display(ops string) *gutrees.Style {
	return &gutrees.Style{Name: "display", Value: ops}
}

// Height provides the height style value
func Height(size Size) *gutrees.Style {
	return &gutrees.Style{Name: "height", Value: string(size)}
}

// FontSize provides the margin style value
func FontSize(size Size) *gutrees.Style {
	return &gutrees.Style{Name: "font-size", Value: string(size)}
}

// Padding provides the margin style value
func Padding(size Size) *gutrees.Style {
	return &gutrees.Style{Name: "padding", Value: string(size)}
}

// Margin provides the margin style value
func Margin(size Size) *gutrees.Style {
	return &gutrees.Style{Name: "margin", Value: string(size)}
}

// Width provides the width style value
func Width(size Size) *gutrees.Style {
	return &gutrees.Style{Name: "width", Value: string(size)}
}
