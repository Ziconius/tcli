package utils

import (
	"io"
	"log/slog"
	"os"
)

func FileContents(filePath string) ([]byte, error) {
	if _, err := os.Stat(filePath); err != nil {
		slog.Error("No file found at target location", "file path", filePath, "error", err)

		return []byte{}, err
	}
	fh, err := os.Open(filePath)
	if err != nil {
		slog.Error("Failed to open file", "error", err)

		return []byte{}, err
	}
	defer fh.Close()
	contents, err := io.ReadAll(fh)
	if err != nil {
		slog.Error("Failed to read file contents.", "error", err)

		return []byte{}, err
	}

	return contents, nil
}

func WriteBytes(filePath string, data []byte) error {
	slog.Debug("Writing bytes to file", "file", filePath, "byte length", len(data))
	fh, err := os.Create(filePath)
	if err != nil {
		return err
	}
	_, err = fh.Write(data)
	if err != nil {
		return err
	}
	
	return nil
}
