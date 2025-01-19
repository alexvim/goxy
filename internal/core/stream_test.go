package core

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoutine(t *testing.T) {

	rdata := []byte{0, 0, 1, 1, 0, 0}

	r := bytes.NewReader(rdata)
	w := bytes.NewBuffer(nil)

	stream(w, r)

	assert.Equal(t, rdata, w.Bytes())
}
