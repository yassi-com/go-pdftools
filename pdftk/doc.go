/*
Package pdftk provides wrapper functions for calling PDFtk commands.

Expects command line executable of pdftk or pdftk-java to be installed.
Use [OptionExecutable] to name a specific binary.

Every function takes a [context.Context] and reports a failure as an [Error],
which carries the process stderr and exit code.

By: Patrick Brown
*/
package pdftk
