package compressextract

import (
	"compress/flate"
	"io"
	"os"
)

// CompressFileFlate compresses a file to DEFLATE
func CompressFileFlate(destFilePath string, sourceFilePath string) error {

	sourceFile, err := os.Open(sourceFilePath)
	if err != nil {
		return err
	}

	destFile, err := os.Create(destFilePath)
	if err != nil {
		return err
	}

	writer, err := flate.NewWriter(destFile, 1)
	if err != nil {
		return err
	}

	defer sourceFile.Close()
	defer destFile.Close()
	defer writer.Close()

	io.Copy(writer, sourceFile)

	writer.Close()

	return nil

}

// ExtractFileFlate extracts compressed DEFLATE file
func ExtractFileFlate(destFilePath string, sourceFilePath string) error {

	sourceFile, err := os.Open(sourceFilePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(destFilePath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	reader := flate.NewReader(sourceFile)
	defer reader.Close()

	io.Copy(destFile, reader)

	return nil

}
