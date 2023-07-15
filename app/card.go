package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
)

type (
	// Card is a struct for generating card image
	// @TODO CARDの数が多すぎる場合はPoolを検討する
	Card struct {
		styledTexts []StyledText
		styles      map[string]struct{} // for duplicate style check
	}
	Cards map[string]*Card

	// StyledText is a struct for text and style
	// This value is used for generating Card
	StyledText struct {
		text  string
		style *Style // have pointer to style to save memory
	}

	Styles map[string]*Style
	// Style is a struct for text style definition
	// can define position, fontsize and color
	Style struct {
		name     string
		position image.Point
		fontsize float64
		rgba     color.RGBA
	}
)

func (s Styles) Pointer(name string) (*Style, error) {
	if _, e := s[name]; !e {
		return nil, errors.New("style is undefined")
	}
	return s[name], nil
}

func (s Styles) String() string {
	str := ""
	for k, v := range s {
		str += "key: " + k + ", {" + v.String() + "}\n"
	}
	return str
}

func (s Style) String() string {
	ss := "name: " + s.name
	ss += ", position: " + s.position.String()
	ss += fmt.Sprintf(", fontsize: %.1f", s.fontsize)
	ss += fmt.Sprintf(", rgba: ( %d, %d, %d, %d )", s.rgba.R, s.rgba.G, s.rgba.B, s.rgba.A)
	return ss
}

func (c *Cards) Add(name string, text string, style *Style) error {
	if style == nil {
		return errors.New("style is undefined")
	}
	if _, ok := (*c)[name]; !ok {
		(*c)[name] = &Card{}
	}
	return (*c)[name].Add(text, style)
}

func (c *Card) Add(text string, style *Style) error {
	if _, ok := (*c).styles[style.name]; ok {
		return errors.New("style is duplicated")
	}
	if (*c).styles == nil {
		(*c).styles = make(map[string]struct{})
	}
	(*c).styles[style.name] = struct{}{}
	if (*c).styledTexts == nil {
		(*c).styledTexts = make([]StyledText, 0)
	}
	(*c).styledTexts = append((*c).styledTexts, StyledText{
		text:  text,
		style: style,
	})
	return nil
}

func (c Cards) String() string {
	st := ""
	for k, v := range c {
		st += "key: " + k + "\n"
		st += "used_styles: "
		for name := range v.styles {
			st += name + ", "
		}
		st += "\n"
		for _, stxt := range v.styledTexts {
			st += fmt.Sprint(stxt, "\n")
		}
	}
	return st
}

func (st StyledText) String() string {
	return "text: " + st.text + fmt.Sprintf(", style: %p", st.style)
}
