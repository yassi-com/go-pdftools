# go-pdftools

PDF utilities to fill PDFs via FDF and manipulating them with PDFtk.

YASSI-maintained fork of [patiek/go-pdftools](https://github.com/patiek/go-pdftools),
Apache-2.0. Consumed by the `pdfkit` service in the main repo.

## FDF

Package `fdf` generates FDF (Forms Data Format) for filling PDF forms with data.

## PDFtk

Package `pdftk` wraps the PDFtk command line tool. Every call takes a
`context.Context` and returns a `*pdftk.Error` carrying the process stderr and
exit code on failure.

Expects the `pdftk` or `pdftk-java` executable to be installed. Use
`pdftk.OptionExecutable` to select a different binary.

```go
var out bytes.Buffer
err := pdftk.Cat(ctx, &out, pdftk.NewInputFileMap("a.pdf", "b.pdf"), nil)
```
