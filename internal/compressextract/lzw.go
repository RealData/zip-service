package compressextract

import (
	"compress/lzw"
	"io"
	"os"
)

// CompressFileFlate compresses a file to LZW
func CompressFileLZW(destFilePath string, sourceFilePath string) error {

	sourceFile, err := os.Open(sourceFilePath)
	if err != nil {
		return err
	}

	destFile, err := os.Create(destFilePath)
	if err != nil {
		return err
	}

	writer := lzw.NewWriter(destFile, lzw.LSB, 8)

	defer sourceFile.Close()
	defer destFile.Close()
	defer writer.Close()

	io.Copy(writer, sourceFile)

	writer.Close()

	return nil

}

// ExtractFileFlate extracts compressed LZW file
func ExtractFileLZW(destFilePath string, sourceFilePath string) error {

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

	reader := lzw.NewReader(sourceFile, lzw.LSB, 8)
	defer reader.Close()

	io.Copy(destFile, reader)

	return nil

}
