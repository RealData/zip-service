package compressextract 
import (
	"archive/zip"
	"compress/gzip"
	"io" 
	"os" 
	"path"
	r "zip-service/internal/readdir"
)

type Extractor func(destFilePath string, sourceFilePath string) error  

// ExtractFileInfo returns list of file info for files in a ZIP file 
func ExtractFileInfo(sourceFilePath string) ([]r.FileInfo, error)  {

	reader, err := zip.OpenReader(sourceFilePath) 
	if err != nil {
		return nil, err  
	} 

	var files []r.FileInfo 

	for _, f := range reader.File {
		files = append(files, r.FileInfo{Name: f.Name, Size: int64(f.CompressedSize64)})
	}

	return files, nil 

}

// ZIPToRawFiles reads files from ZIP file 
func ZIPToRawFiles(destDirPath string, sourceFilePath string) ([]r.FileInfo, error) {

	zipReader, err := zip.OpenReader(sourceFilePath) 
	if err != nil {
		return nil, err 
	} 
	defer zipReader.Close() 
	
	var files []r.FileInfo  

	for _, f := range zipReader.File {
		reader, err := f.OpenRaw() 
		writer, err := os.Create(path.Join(destDirPath, f.Name)) 
		_, err = io.Copy(writer, reader) 
		if err != nil {
			return nil, err 
		} 
		files = append(files, r.FileInfo{f.Name, int64(f.CompressedSize64)})
		writer.Close()
	}

	return files, nil 

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

// ExtractEachFile extracts files from a list using specified compressor 
func ExtractEachFile(destDirPath string, sourceDirPath string, sourceFiles []r.FileInfo, extractor Extractor) error { 
	for _, file := range sourceFiles { 
		destFilePath := path.Join(destDirPath, file.Name)  
		sourceFilePath := path.Join(sourceDirPath, file.Name) 
		err := extractor(destFilePath, sourceFilePath)   
		if err != nil {
			return err 
		}
	} 
	return nil 
} 

// UnzipAndExtractFiles unzip ZIP file, writes them into temporary directory, and then extracts each file 
func UnzipAndExtractFiles(destDirPath string, sourceFilePath string, extractor Extractor) error {
	
	tempDirPath, err := os.MkdirTemp(destDirPath, "TEMP") 
	if err != nil {
	 	return err 
	}

	files, err := ZIPToRawFiles(tempDirPath, sourceFilePath)  
	if err != nil {
		return err 
    }

	err = ExtractEachFile(destDirPath, tempDirPath, files, extractor) 
	if err != nil {
		return err 
	}

	return nil  

}