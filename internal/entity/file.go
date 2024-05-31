package entity

import (
	"bytes"
	"io"
)

type File struct {
	Buffer   *bytes.Buffer
	Filename string
	Size     int64
}

type FileInfo struct {
	Id       string
	Filename string
}

type FileResult struct {
	io.ReadCloser
}
