package main

import "GoFast/shared"

func main() {
	doc := shared.NewParserDoc("example.json")
	if err := doc.Parse(); err != nil {
		panic(err)
	}
}