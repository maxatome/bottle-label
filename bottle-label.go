package main

import (
	"flag"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	pageMargin     = 10 // mm
	labelWidth     = 47 // mm
	labelHeight    = 55 // mm
	labelMargin    = 7  // mm
	circleRay      = 16 // mm
	textLineHeight = 4  // mm
)

func usage() {
	flag.Usage()
	os.Exit(1)
}

func checkError(pdf *gofpdf.Fpdf) {
	if err := pdf.Error(); err != nil {
		fmt.Printf("An error occured: %s\n", err)
		os.Exit(1)
	}
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [OPTIONS] (NUMx)LABEL ...\n",
			os.Args[0])
		flag.PrintDefaults()
	}

	var fontFile string
	flag.StringVar(&fontFile, "font", "",
		"JSON font file, relative to current directory")

	var output string
	flag.StringVar(&output, "output", "",
		"PDF output file (defaults to output.pdf)")

	flag.Parse()

	if flag.NArg() == 0 {
		usage()
	}

	parse := regexp.MustCompile(`^(\d+)\s*x\s*(.*)\z`)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetAcceptPageBreakFunc(func() bool { return false })

	tr := pdf.UnicodeTranslatorFromDescriptor("")

	if fontFile != "" {
		pdf.AddFont("MyFont", "", fontFile)
		checkError(pdf)
		pdf.SetFont("MyFont", "", 11)
	} else {
		pdf.SetFont("helvetica", "", 11)
	}
	checkError(pdf)

	var labels [][]string
	for _, arg := range flag.Args() {
		num := 1
		if match := parse.FindStringSubmatch(arg); len(match) == 3 {
			num, _ = strconv.Atoi(match[1])
			arg = match[2]
		}

		lines := strings.Split(arg, `\n`)
		if len(lines) > 3 {
			fmt.Fprintf(os.Stderr, "3 lines per label max: `%s'\n", arg)
			os.Exit(1)
		}

		for _, line := range lines {
			dlines := pdf.SplitLines([]byte(tr(line)), 47)
			if len(dlines) > 1 {
				fmt.Fprintf(os.Stderr, "Line too long: `%s'\n", line)
				os.Exit(1)
			}
		}

		for i := 0; i < num; i++ {
			labels = append(labels, lines)
		}
	}

	if len(labels) == 0 {
		usage()
	}

	idxLabel := 0
main:
	for {
		pdf.AddPage()

		for ny, y := 1, float64(pageMargin); ny <= 5; ny++ {
			for nx, x := 1, float64(pageMargin); nx <= 4; nx++ {
				pdf.Rect(x, y, labelWidth, labelHeight, "")
				pdf.Circle(x+labelMargin+circleRay+.5, y+labelMargin+circleRay,
					circleRay, "")

				pdf.SetXY(x, y+labelMargin+circleRay*2+labelMargin/2)
				pdf.MultiCell(labelWidth, textLineHeight,
					tr(strings.Join(labels[idxLabel], "\n")), "", "C", false)

				idxLabel++
				if idxLabel == len(labels) {
					break main
				}
				x += labelWidth
			}
			y += labelHeight
		}
	}

	if output == "" {
		output = "output.pdf"
	}
	err := pdf.OutputFileAndClose(output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot write %s: %s\n", output, err)
		os.Exit(1)
	}
	fmt.Printf("%s generated.\n", output)
}
