package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/napcs/opml2pptx/pkg/opml"
	"github.com/napcs/opml2pptx/pkg/pptx"
)

var version = "dev" // default fallback; overridden at build time

func main() {
	inputPath := flag.String("input", "", "Path to the input OPML file")
	outputPath := flag.String("output", "", "Path to the output PPTX file")
	showVersion := flag.Bool("version", false, "Show version and exit")

	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}
	if *inputPath == "" || *outputPath == "" {
		fmt.Println("Usage: opml2pptx -input input.opml -output output.pptx")
		os.Exit(1)
	}

	if _, err := os.Stat(*inputPath); err != nil {
		log.Fatalf("Input file does not exist: %s", *inputPath)
	}

	// Ensure output path ends in .pptx
	if filepath.Ext(*outputPath) != ".pptx" {
		log.Fatalf("Output file must have .pptx extension")
	}

	// Parse OPML
	pres, err := opml.ParseOPMLFile(*inputPath)
	if err != nil {
		log.Fatalf("Failed to parse OPML: %v", err)
	}

	// Generate PPTX
	if err := pptx.BuildPPTX(pres, *outputPath); err != nil {
		log.Fatalf("Failed to build PPTX: %v", err)
	}

	fmt.Printf("Presentation generated: %s\n", *outputPath)
}
