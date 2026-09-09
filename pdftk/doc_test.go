package pdftk_test

import (
	"bytes"
	"context"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/yassi-com/go-pdftools/fdf"
	"github.com/yassi-com/go-pdftools/pdftk"
)

func ExampleCat() {
	ctx := context.Background()

	// file to write output into
	f, err := os.Create("out.pdf")
	if err != nil {
		// handle error
	}

	err = pdftk.Cat(ctx, f, pdftk.NewInputFileMap("first.pdf", "second.pdf", "third.pdf"), []pdftk.PageRange{
		{
			FileHandleName: pdftk.InputHandleNameFromInt(1),
			Rotation:       pdftk.East,
		},
		{
			FileHandleName: pdftk.InputHandleNameFromInt(0),
		},
	}, pdftk.OptionFlatten())
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleCat_all() {
	ctx := context.Background()

	// file to write output into
	f, err := os.Create("out.pdf")
	if err != nil {
		// handle error
	}

	// omitting page ranges concatenates every input in file handle order
	err = pdftk.Cat(ctx, f, pdftk.NewInputFileMap("first.pdf", "second.pdf"), nil)
	if err != nil {
		// handle error
	}
}

func ExampleCat_pageRanges() {
	ctx := context.Background()

	// file to write output into
	f, err := os.Create("out.pdf")
	if err != nil {
		// handle error
	}

	inputFiles := pdftk.InputFileMap{
		"A": "first.pdf",
		"B": "second.pdf",
		"C": "third.pdf",
	}

	pageRanges := []pdftk.PageRange{
		// first we take page 2 from C
		{
			FileHandleName: "C",
			BeginPage:      2,
			EndPage:        2,
		},
		// then we take pages 4 until end of A and rotate them all
		{
			FileHandleName: "A",
			BeginPage:      4,
			Rotation:       pdftk.East,
		},
		// then we take all odd pages of B
		{
			FileHandleName: "B",
			Qualifier:      pdftk.Odd,
		},
		// finally we add in the first page from C
		{
			FileHandleName: "C",
			BeginPage:      1,
			EndPage:        1,
		},
	}

	err = pdftk.Cat(ctx, f, inputFiles, pageRanges)
	if err != nil {
		// handle error
	}
}

func ExampleFillForm() {
	ctx := context.Background()

	var b bytes.Buffer
	if err := fdf.Write(&b, fdf.Inputs{
		"first field":  "hello",
		"second field": "world",
	}); err != nil {
		// handle error
	}

	f, err := os.Create("out.pdf")
	if err != nil {
		// handle error
	}
	defer f.Close()

	// fill form with FDF data from buffer b and flatten
	err = pdftk.FillForm(ctx, f, "test.pdf", &b, pdftk.OptionFlatten())
	if err != nil {
		// handle error
	}
}

func ExampleFillForm_file() {
	ctx := context.Background()

	fdfFile, err := os.Open("input.fdf")
	if err != nil {
		// handle error
	}
	defer fdfFile.Close()

	outFile, err := os.Create("out.pdf")
	if err != nil {
		// handle error
	}
	defer outFile.Close()

	// fill form with FDF data from FDF file and flatten
	err = pdftk.FillForm(ctx, outFile, "test.pdf", fdfFile, pdftk.OptionFlatten())
	if err != nil {
		// handle error
	}
}

func ExampleBackground() {
	ctx := context.Background()

	backgroundFile, err := os.Open("background.pdf")
	if err != nil {
		// handle error
	}
	defer backgroundFile.Close()

	outFile, err := os.Create("out.pdf")
	if err != nil {
		// handle error
	}
	defer outFile.Close()

	// add backgroundFile to background of input.pdf
	err = pdftk.Background(ctx, outFile, "input.pdf", backgroundFile)
	if err != nil {
		// handle error
	}
}

func ExampleStamp() {
	ctx := context.Background()

	stampFile, err := os.Open("stamp.pdf")
	if err != nil {
		// handle error
	}
	defer stampFile.Close()

	outFile, err := os.Create("out.pdf")
	if err != nil {
		// handle error
	}
	defer outFile.Close()

	// stamp input.pdf with stampFile
	err = pdftk.Stamp(ctx, outFile, "input.pdf", stampFile)
	if err != nil {
		// handle error
	}
}

func ExampleError() {
	ctx := context.Background()

	var out bytes.Buffer
	err := pdftk.Cat(ctx, &out, pdftk.NewInputFileMap("encrypted.pdf"), nil)

	var toolErr *pdftk.Error
	if errors.As(err, &toolErr) {
		// classify the failure from what pdftk reported
		if strings.Contains(strings.ToLower(toolErr.Stderr), "owner password") {
			log.Print("input is password protected")
		}
	}
}
