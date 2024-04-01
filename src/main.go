package main

import (
	"fmt"
	"os"

	"gecko/argparse"
	"gecko/colors"
)

func main() {
	args := os.Args[1:]
	// fmt.Printf("All: '%s'\n", args[:])

	conf := argparse.ParseArguments(args)
	//fmt.Println(conf.TextInput)

	pcolors := colors.LoadPossibleColors(os.Getenv("TERM"), conf)
	fmt.Println(colors.MarkupText(pcolors, conf))
}
