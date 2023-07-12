package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	// style test definition
	correctStyleCSV = `
	name,font_size,x,y,color_r,color_g,color_b,color_a
	title,30,100,100,0,0,0,255
	description,15,200,200,100,150,200,255
	cost,15,10,10,255,255,255,100`
	failedStyleCSV = `
	name,font_size,x,y,color_r,color_g,color_b,color_a
	title,30,100,100,0,0,0,255
	description,15,200,200,100,150,200,255
	title,35,100,100,0,0,0,255`

	// card test definition
	correctCardCSV = `
	name,style,text
	main,title,some title
	main,description,test
	sub1,title,sub1 title
	sub1,description
	main,cost,12
	`
	duplicateCardCSV = `
	name,style,text
	main,title,some title
	main,title,test
	sub1,sub1 title, title
	sub1,sub1 description, description
	`
	undefinedStyleCardCSV = `
	name,style,text
	main,undefined,some title
	main,title,some title
	`
)

func TestCSVDecoder(t *testing.T) {
	_, err := CSVDecoder.DecodeStyles([]byte(failedStyleCSV))
	assert.NotNil(t, err, "duplicate card name should be error")
	styles, err := CSVDecoder.DecodeStyles([]byte(correctStyleCSV))
	assert.Nil(t, err, "style csv should be decoded")

	// TODO: add style is expected check
	// cards, err := CSVDecoder.DecodeCards([]byte(correctCardCSV), styles)
	// TODO: add card is expected check
	if _, err := CSVDecoder.DecodeCards([]byte(duplicateCardCSV), styles); err == nil {
		t.Error("duplicate card name should be error")
	}
	if _, err := CSVDecoder.DecodeCards([]byte(undefinedStyleCardCSV), styles); err == nil {
		t.Error("undefined style assign should be error")
	}
}
