package grpc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math"
	"object-storage/api/objectstorage"
	"object-storage/internal/entity"
	"object-storage/pkg/filereader"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const BUFFER_WINDOW = 512

// var _ objectstorage.ObjectStorageServer = (*Server)(nil)

type Storage interface {
	Put(ctx context.Context, reader filereader.Reader) (*entity.FileInfo, error)
	Get(ctx context.Context, filename string) (*entity.FileResult, uint64, error)
}

type Server struct {
	storage Storage

	objectstorage.UnimplementedObjectStorageServer
}

// Get implements objectstorage.ObjectStorageServer.
func (s *Server) Get(request *objectstorage.GetRequest, stream objectstorage.ObjectStorage_GetServer) error {
	f, size, err := s.storage.Get(stream.Context(), request.FileName)
	if err != nil {
		return status.Errorf(codes.Internal, "Unexpected: %s", err.Error())
	}

	for {
		chunkSize := int(math.Min(BUFFER_WINDOW, float64(size)))
		buffer := make([]byte, chunkSize)
		_, err := f.Read(buffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return status.Errorf(codes.Internal, "Unexpected: %s", err.Error())
		}

		size -= BUFFER_WINDOW

		err = stream.Send(&objectstorage.File{
			Chunk: buffer,
		})
		if err != nil {
			return status.Errorf(codes.Internal, "Unexpected: %s", err.Error())
		}
	}
	return nil
}

func (s *Server) Upload(stream objectstorage.ObjectStorage_UploadServer) error {
	file := &entity.File{
		Buffer: bytes.NewBuffer(nil),
	}

	for {
		object, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return status.Error(codes.Internal, fmt.Sprintf("Unexpected error: %s", err.Error()))
		}

		if object.Meta != nil {
			file.Filename = object.Meta.Filename
		}

		chunk := object.Chunk

		file.Size += int64(len(chunk))
		n, err := file.Buffer.Write(chunk)
		if err != nil {
			fmt.Printf("%s", err.Error())
			return err
		}

		fmt.Printf("writted %d bytes\n", n)
		// fmt.Printf("chunk: %v\n", chunk)
	}

	slog.Info("new upload", slog.String("filename", file.Filename), slog.Int64("size", file.Size), slog.Int("buffer len", file.Buffer.Len()))

	reader, err := s.storage.Put(stream.Context(), filereader.New(file.Buffer, file.Size, file.Filename))
	if err != nil {
		return err
	}

	return stream.SendAndClose(&objectstorage.UploadResponse{
		FileName: reader.Filename,
	})
}

func New(storage Storage) *Server {
	return &Server{
		storage: storage,
	}
}
