package util

import (
	"io"
	"io/ioutil"
	"strings"

	"github.com/google/uuid"
)

func GetUUID() string {
	id := uuid.Must(uuid.NewRandom())

	return id.String()
}

func StrToReadCloser(str string) io.ReadCloser {
	r := ioutil.NopCloser(strings.NewReader(str))

	return r
}
