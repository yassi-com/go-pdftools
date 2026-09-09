package pdftk

import (
	"context"
	"io"
)

// Concatenate PDF files and write to out.
// Specify pageRanges to manipulate ordering, ranges, and rotation of pages.
// Omit pageRanges to concatenate every input in file handle order.
func Cat(ctx context.Context, out io.Writer, fileNames InputFileMap, pageRanges []PageRange, options ...Option) error {
	return newCommand(out, nil, options, catArgs(fileNames, pageRanges)...).run(ctx)
}

// catArgs builds the argument list Cat invokes pdftk with.
func catArgs(fileNames InputFileMap, pageRanges []PageRange) []string {
	args := append(fileNames.parameterize(), "cat")
	args = append(args, pageRangesToStrings(pageRanges)...)

	return append(args, "output", "-")
}

// Fill a PDF file with FDF data and write to out.
func FillForm(ctx context.Context, out io.Writer, inputFileName string, fdf io.Reader, options ...Option) error {
	return newCommand(out, fdf, options, inputFileName, "fill_form", "-", "output", "-").run(ctx)
}

// Applies a PDF watermark to the background of a each page of input pdf.
// If background is multiple pages, only first page is used for the background.
func Background(ctx context.Context, out io.Writer, inputFileName string, background io.Reader, options ...Option) error {
	return newCommand(out, background, options, inputFileName, "background", "-", "output", "-").run(ctx)
}

// Applies a PDF watermark to the background of a each page of input pdf.
// Similar to background but applies each page of background to corresponding
// input page. If input pdf has more pages than background, last page of
// background is used for the remainder of the input pages.
func MultiBackground(ctx context.Context, out io.Writer, inputFileName string, background io.Reader, options ...Option) error {
	return newCommand(out, background, options, inputFileName, "multibackground", "-", "output", "-").run(ctx)
}

// Stamp (overlay) each page of input pdf with a stamp PDF and write to out.
// If stamp is multiple pages, only the first page is used for the stamp.
func Stamp(ctx context.Context, out io.Writer, inputFileName string, stamp io.Reader, options ...Option) error {
	return newCommand(out, stamp, options, inputFileName, "stamp", "-", "output", "-").run(ctx)
}

// MultiStamp (overlay) each page of input pdf with stamp PDF and write to out.
// Similar to stamp but applies each page of stamp to corresponding input page.
// If input pdf has more pages than stamp, last page of stamp is used for the
// remainder of the input pages.
func MultiStamp(ctx context.Context, out io.Writer, inputFileName string, stamp io.Reader, options ...Option) error {
	return newCommand(out, stamp, options, inputFileName, "multistamp", "-", "output", "-").run(ctx)
}
