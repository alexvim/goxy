package core

import (
	"io"
)

func stream(dst io.Writer, src io.Reader) (int64, error) {
	return io.Copy(dst, src)
}
