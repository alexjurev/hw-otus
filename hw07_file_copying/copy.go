package main

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/cheggaaa/pb"
	"github.com/spf13/afero"
)

func closeFile(file afero.File, fromPath string, finError *error) {
	if file != nil {
		err := file.Close()
		if err != nil {
			*finError = fmt.Errorf("can`t close file %s: %w", fromPath, err)
		}
	}
}

func CopyFile(fs afero.Fs, fromPath, toPath string, limit, offset, chunkSize int64) (finError error) {
	if limit < 0 {
		limit = 0
	}

	if offset < 0 {
		offset = 0
	}

	src, err := fs.Open(fromPath)
	defer closeFile(src, fromPath, &finError)

	if err != nil {
		return fmt.Errorf("can`t open file %s: %w", fromPath, err)
	}

	fileStats, err := src.Stat()
	if err != nil {
		return fmt.Errorf("can`t stat file %s: %w", fromPath, err)
	}

	fileSize := fileStats.Size()
	if fileSize < offset {
		return fmt.Errorf("offset exceeds file size (%d bytes)", fileSize)
	}
	if limit > fileSize {
		limit = fileSize
	}

	dst, err := fs.Create(toPath)
	defer closeFile(dst, toPath, &finError)

	if err != nil {
		return fmt.Errorf("can`t create file %s: %w", toPath, err)
	}

	buf := make([]byte, chunkSize)
	initialOffset := offset
	var isLastChunk bool
	var totalRead int64
	bar := pb.StartNew(int(fileSize / 1024))
	for offset < fileSize {
		read, err := src.ReadAt(buf, offset)
		if !errors.Is(err, io.EOF) && err != nil {
			bar.Finish()
			return fmt.Errorf("can`t read from file %s: %w", fromPath, err)
		}
		if errors.Is(err, io.EOF) {
			isLastChunk = true
		}

		if limit > 0 && totalRead+int64(read) > limit {
			_, err = dst.WriteAt(buf[:(limit-totalRead)], offset-initialOffset)
			isLastChunk = true
		} else {
			_, err = dst.WriteAt(buf[:read], offset-initialOffset)
		}

		if err != nil {
			bar.Finish()
			return fmt.Errorf("can`t write to file %s: %w", toPath, err)
		}

		offset += int64(read)
		totalRead += int64(read)

		bar.Increment()
		time.Sleep(time.Millisecond) // можно поменять на секунду, тогда успеет отобразиться движение прогресс бара

		if isLastChunk {
			bar.Finish()
			break
		}
	}

	return nil
}
