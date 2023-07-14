package main

import (
	"errors"
	"image"
	"image/color"

	"github.com/gocarina/gocsv"
)

// Converter is an interface for converting csv file
type (
	Decoder interface {
		// DecodeStyle returns single style
		DecodeStyles(data []byte) (Styles, error)
		// DecodeCard returns single card
		DecodeCards(data []byte, styles Styles) (*Cards, error)
	}
	// csvConverter is a struct for converting csv file
	CsvDecoder struct{}

	StyleCSV struct {
		Name     string  `csv:"name"`
		FontSize float64 `csv:"font_size"`
		X        int     `csv:"x"`
		Y        int     `csv:"y"`
		ColorR   int     `csv:"color_r"`
		ColorG   int     `csv:"color_g"`
		ColorB   int     `csv:"color_b"`
		ColorA   int     `csv:"color_a"`
	}

	CardCSV struct {
		Name  string `csv:"name"`
		Style string `csv:"style"`
		Text  string `csv:"text"`
	}
)

var (
	CSVDecoder Decoder = CsvDecoder{}
)

// DecodeStyle returns single style
func (c CsvDecoder) DecodeStyles(data []byte) (Styles, error) {
	s := make(Styles)
	if err := gocsv.UnmarshalBytesToCallback(data, func(sc StyleCSV) error {
		p := image.Point{sc.X, sc.Y}
		rgba := color.RGBA{
			uint8(sc.ColorR),
			uint8(sc.ColorG),
			uint8(sc.ColorB),
			uint8(sc.ColorA),
		}
		if _, ok := s[sc.Name]; ok {
			return errors.New("style name is duplicated")
		}
		s[sc.Name] = &Style{
			name:     sc.Name,
			position: p,
			fontsize: sc.FontSize,
			rgba:     rgba,
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return s, nil
}

// DecodeCard returns single card
func (c CsvDecoder) DecodeCards(data []byte, styles Styles) (Cards, error) {
	cards := make(Cards, 0)
	if err := gocsv.UnmarshalBytesToCallback(data, func(cc CardCSV) error {
		sp, err := styles.Pointer(cc.Style)
		if err != nil {
			return errors.New("style is undefined")
		}
		return cards.Add(cc.Name, cc.Text, sp)
	}); err != nil {
		return nil, err
	}
	return cards, nil
}
