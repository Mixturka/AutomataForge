package sliceutils

import (
	"bytes"
	"encoding/gob"
)

func HashSlice[T any](slice []T) string {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(slice)
	return buf.String()
}
