package filters

import "io"

type IFilter interface {
	Filter(io.Reader, io.Writer) error
}
