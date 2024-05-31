package filereader

import (
	"io"
	"log/slog"
)

type Reader interface {
	Size() int64
	Filename() string
	io.Reader
}

type FileReader struct {
	size     int64
	filename string
	io.Reader
}

func New(r io.Reader, size int64, filename string) *FileReader {
	slog.Info("creating file reader", slog.String("filename", filename), slog.Int64("size", size), slog.Any("reader", r))
	return &FileReader{
		size:     size,
		filename: filename,
		Reader:   r,
	}
}

func (r FileReader) Filename() string {
	return r.filename
}

func (r FileReader) Size() int64 {
	return r.size
}
