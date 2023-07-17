package main

import (
	"flag"
	"fmt"
)

func main() {
	stylePath := flag.String("style", "", "Path to the style file")
	cardsPath := flag.String("cards", "", "Path to the cards file")

	flag.Parse()

	fmt.Println("Style File Path:", *stylePath)
	fmt.Println("Cards File Path:", *cardsPath)
}
