package filter_drivers

import (
	"io"

	aes2htm "github.com/romiras/aes2htm/pkg/converter"
	"github.com/romiras/bv/filters"
)

type AnsiTerminalToHTML struct{}

func NewAnsiTerminalToHTML() filters.IFilter {
	return &AnsiTerminalToHTML{}
}

func (ah *AnsiTerminalToHTML) Filter(r io.Reader, w io.Writer) error {
	conv, err := aes2htm.NewAes2Htm(w)
	if err != nil {
		return err
	}

	return conv.WriteHTML(r)
}
