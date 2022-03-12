package compressextract 

import (
	"compress/zlib" 
	"io"
	"os"
)

// CompressFileFlate compresses a file to ZLIB    
func CompressFileZLIB(destFilePath string, sourceFilePath string) error {

	sourceFile, err := os.Open(sourceFilePath)
	if err != nil {
		return err 
	}
	
	destFile, err := os.Create(destFilePath) 
	if err != nil {
		return err 
	} 

	writer := zlib.NewWriter(destFile) 

	defer sourceFile.Close() 
	defer destFile.Close()
	defer writer.Close() 

	io.Copy(writer, sourceFile) 
	
	writer.Close() 

	return nil 

} 

// ExtractFileFlate extracts compressed ZLIB file  
func ExtractFileZLIB(destFilePath string, sourceFilePath string) error {
	
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

	reader, err := zlib.NewReader(sourceFile) 
	if err != nil {
		return err 
	} 
	defer reader.Close() 

	io.Copy(destFile, reader) 

	return nil 

}
