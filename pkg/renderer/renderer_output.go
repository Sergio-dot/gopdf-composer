package renderer

import "io"

func (r *Renderer) SaveToFile(path string) error {
	return r.pdf.OutputFileAndClose(path)
}

func (r *Renderer) WriteTo(w io.Writer) (int64, error) {
	return 0, r.pdf.Output(w)
}
