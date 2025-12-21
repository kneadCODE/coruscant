package main

import (
	"log"

	"github.com/kneadCODE/go-diagrams/diagram"
)

func main() {
	if err := genGovernance(); err != nil {
		log.Fatal(err)
	}
	if err := genHub(); err != nil {
		log.Fatal(err)
	}
}

func defaultDiagramAttributes() diagram.Option {
	return diagram.WithAttributes(map[string]string{
		"nodesep":   "2.0",      // Horizontal spacing between nodes (increased)
		"ranksep":   "1.5",      // Vertical spacing between ranks/levels (increased)
		"splines":   "polyline", // Polyline edges for better balanced tree layout
		"overlap":   "false",    // Remove node overlaps using Prism algorithm
		"fixedsize": "false",    // Allow nodes to auto-size based on label
		"ordering":  "out",      // Preserve edge order for balanced tree layout
	})
}
