package compressextract  
import (
	"archive/zip"  
	"io"
	"os"
	"path"
	r "zip-service/internal/readdir"
) 

func CompressFiles(zipFileStr string, dirStr string, files []r.FileInfo) error {

	zipFile, err := os.Create(zipFileStr) 
	if err != nil {
		return err 
	}
	defer zipFile.Close() 

	writer := zip.NewWriter(zipFile) 
	defer writer.Close() 
	
	for _, file := range files {
		
		f, err := os.Open(path.Join(dirStr, file.Name)) 
		if err != nil {
			return err 
		}
		defer f.Close() 
		
		w, err := writer.Create(file.Name) 
		if err != nil {
			return err 
		} 
		defer writer.Close()

		io.Copy(w, f)

	} 

	return nil 

} 

