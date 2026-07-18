package renderer

import "io"

// SaveToFile writes the PDF to the given file path and closes the document.
func (r *Renderer) SaveToFile(path string) error {
	return r.pdf.OutputFileAndClose(path)
}

// WriteTo writes the PDF to the given io.Writer.
func (r *Renderer) WriteTo(w io.Writer) (int64, error) {
	return 0, r.pdf.Output(w)
}
