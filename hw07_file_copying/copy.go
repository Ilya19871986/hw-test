package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	from, to                 string
	limit, offset            int64
	ErrIllegalArgument       = errors.New("errors illegal arguments")
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

// Проверяет входные параметры.
func validateInput(fromPath, toPath string) error {
	if fromPath == "" || toPath == "" {
		return ErrIllegalArgument
	}
	return nil
}

// Информацию о файле.
func getFileInfo(file *os.File) (os.FileInfo, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	return fileInfo, nil
}

// Проверяет, что файл поддерживается.
func validateFile(fileInfo os.FileInfo) error {
	if fileInfo.Size() == 0 {
		return ErrUnsupportedFile
	}
	return nil
}

// Вычисляет лимит для копирования.
func calculateLimit(fileSize, offset, limit int64) int64 {
	available := fileSize - offset
	if limit == 0 || limit > available {
		return available
	}
	return limit
}

// Создает и настраивает прогресс-бар.
func setupProgressBar(limit int64) *pb.ProgressBar {
	bar := pb.Full.Start64(limit)
	bar.Set(pb.Bytes, true)
	bar.Set(pb.SIBytesPrefix, true)
	return bar
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	if err := validateInput(fromPath, toPath); err != nil {
		return err
	}

	srcFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	fileInfo, err := getFileInfo(srcFile)
	if err != nil {
		return err
	}

	if err := validateFile(fileInfo); err != nil {
		return err
	}

	if offset > fileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	limit = calculateLimit(fileInfo.Size(), offset, limit)

	if _, err = srcFile.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	dstFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	bar := setupProgressBar(limit)
	defer bar.Finish()

	reader := bar.NewProxyReader(srcFile)
	if _, err = io.CopyN(dstFile, reader, limit); err != nil {
		if err != io.EOF {
			return fmt.Errorf("copy error: %w", err)
		}
	}

	return nil
}
