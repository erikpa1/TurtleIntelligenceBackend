// Command dwg2png renders a DWG file's flat 2D geometry to a PNG preview.
package main

import (
	"flag"
	"fmt"
	"os"

	"turtle/dwg"
)

func main() {
	width := flag.Int("width", 1600, "output image width")
	height := flag.Int("height", 1200, "output image height")
	flag.Parse()

	if flag.NArg() != 2 {
		fmt.Fprintf(os.Stderr, "usage: dwg2png <in.dwg> <out.png>\n")
		os.Exit(2)
	}
	inPath, outPath := flag.Arg(0), flag.Arg(1)

	doc, err := dwg.ParseFile(inPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dwg2png: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "dwg2png: parsed %d entities (version %s)\n", len(doc.Entities), doc.Version)

	img := dwg.Render(doc, dwg.RenderOptions{Width: *width, Height: *height})

	out, err := os.Create(outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dwg2png: %v\n", err)
		os.Exit(1)
	}
	defer out.Close()

	if err := dwg.EncodePNG(out, img); err != nil {
		fmt.Fprintf(os.Stderr, "dwg2png: %v\n", err)
		os.Exit(1)
	}
}
