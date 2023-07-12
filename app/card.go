package main

import (
	"errors"
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

func (c *Cards) Add(name string, text string, style *Style) error {
	if style == nil {
		return errors.New("style is undefined")
	}
	if _, e := (*c)[name]; !e {
		(*c)[name] = &Card{}
	}
	return (*c)[name].Add(text, style)
}

func (c *Card) Add(text string, style *Style) error {

	if _, e := (*c).styles[style.name]; e {
		return errors.New("style is duplicated")
	}
	(*c).styles[style.name] = struct{}{}
	(*c).styledTexts = append((*c).styledTexts, StyledText{
		text:  text,
		style: style,
	})
	return nil
}
