package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
)

func parseArgs(args []string) {
	// Create new parser object
	parser := argparse.NewParser("MediumScraper", "Scraper for medium.com")

	// Create flags
	input := parser.String("i", "input", &argparse.Options{
		Required: true,
		Help:     "Input link"})

	randomUA := parser.Flag("", "random-ua", &argparse.Options{
		Required: false,
		Help:     "Randomize user agent on request"})

	// Parse input
	err := parser.Parse(args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		os.Exit(0)
	}

	// Fill arguments structure
	arguments.Input = *input
	arguments.RandomUA = *randomUA
}
