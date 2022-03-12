package compressextract 

import (
	"compress/gzip" 
	"io"
	"os"
)

// CompressFileGZIP compresses a file to GZIP 
func CompressFileGZIP(destFilePath string, sourceFilePath string) error {

	sourceFile, err := os.Open(sourceFilePath)
	if err != nil {
		return err 
	}
	
	destFile, err := os.Create(destFilePath) 
	if err != nil {
		return err 
	} 

	writer := gzip.NewWriter(destFile) 
	
	defer sourceFile.Close() 
	defer destFile.Close()
	defer writer.Close() 

	io.Copy(writer, sourceFile) 
	
	writer.Close() 

	return nil 

} 

// ExtractFileGZIP extracts compressed GZIP file 
func ExtractFileGZIP(destFilePath string, sourceFilePath string) error {
	
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

	reader, err := gzip.NewReader(sourceFile) 
	if err != nil {
		return err 
	} 
	defer reader.Close() 

	io.Copy(destFile, reader) 

	return nil 

}
