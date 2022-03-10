package compressextract  
import (
	"archive/zip"  
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path" 
	"time"
	r "zip-service/internal/readdir"
) 

type Compressor func(destFilePath string, sourceFilePath string) error  

func CompressFileGZIP(destFilePath string, sourceFilePath string) error {

	sourceFile, err := os.Open(sourceFilePath)
	if err != nil {
		return err 
	}
	
	destFile, err := os.Create(destFilePath) 
	if err != nil {
		return err 
	} 

	w := gzip.NewWriter(destFile) 
	
	defer sourceFile.Close() 
	defer destFile.Close()
	defer w.Close() 

	io.Copy(w, sourceFile) 
	
	w.Close() 

	return nil 

}

func parallelCompressEachFile(destDirPath string, destExt string, sourceDirPath string, sourceFiles []r.FileInfo, threads int, compressor Compressor) (err error) { 

	numJobs := len(sourceFiles) 
	jobChan := make(chan r.FileInfo, numJobs) 
	resChan := make(chan error, numJobs) 
	
	if threads > numJobs {
		threads = numJobs  
	}

	for w := 1; w <= threads; w++ {
		go func(jobChan <-chan r.FileInfo, resChan chan<- error) {
			for file := range jobChan { 
				fmt.Println(time.Now())
				destFilePath := path.Join(destDirPath, file.Name) + "." + destExt 
				sourceFilePath := path.Join(sourceDirPath, file.Name) 
				fmt.Println(sourceFilePath)
				err := compressor(destFilePath, sourceFilePath)  
				resChan <- err 
				fmt.Println(time.Now())

			}
		}(jobChan, resChan)
	} 

	for _, file := range sourceFiles {
		jobChan <- file 
	}
	close(jobChan) 

	for j := 1; j <= numJobs; j++ {
		res := <- resChan 
		if res != nil {
			err = res  
		}
	}

	return 

}

func compressEachFile_(destDirPath string, destExt string, sourceDirPath string, sourceFiles []r.FileInfo, compressor Compressor) error { 
	for _, file := range sourceFiles { 
		destFilePath := path.Join(destDirPath, file.Name) + "." + destExt 
		sourceFilePath := path.Join(sourceDirPath, file.Name) 
		err := compressor(destFilePath, sourceFilePath)   
		if err != nil {
			return err 
		}
	} 
	return nil 
}

func CompressAndZipFiles(destDirPath string, destFilePath string, destExt string, sourceDirPath string, sourceFiles []r.FileInfo, threads int, compressor Compressor) (string, error) {

	dirPath, err := os.MkdirTemp(destDirPath, "TEMP") 
	defer os.RemoveAll(dirPath)
	
	if err != nil {
		return "", err  
	}

	err = parallelCompressEachFile(dirPath, destExt, sourceDirPath, sourceFiles, threads, compressor) 
	if err != nil {
		return "", err 
	}

	err = RawFilesToZIP(destFilePath, destExt, dirPath, sourceFiles) 
	if err != nil {
		return "", err 
	}
	
	return dirPath, nil 

}

func RawFilesToZIP(destFilePath string, destExt string, sourceDirPath string, sourceFiles []r.FileInfo) error {

	destFile, err := os.Create(destFilePath) 
	if err != nil {
		return err 
	} 

	zipWriter := zip.NewWriter(destFile)  
	
	defer destFile.Close() 
	defer zipWriter.Close() 

	for _, file := range sourceFiles {

		sourceFilePath := path.Join(sourceDirPath, file.Name) + "." + destExt 
		sourceFile, err := os.Open(sourceFilePath)
		if err != nil {
			return err 
		} 

		fileWriter, err := zipWriter.CreateRaw(&zip.FileHeader{Name: file.Name, CompressedSize64: uint64(file.Size), UncompressedSize64: uint64(file.Size)})  
		if err != nil {
			return err 
		} 

		io.Copy(fileWriter, sourceFile) 

	}

	zipWriter.Close()

	return nil 

}

