package compressextract  
import (
	"archive/zip"  
	"compress/gzip"
	"io"
	"os"
	"path" 
	r "zip-service/internal/readdir"
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

// parallelCompressEachFile concurently compresses each file from a list using specified compressor 
func parallelCompressEachFile(destDirPath string, sourceDirPath string, sourceFiles []r.FileInfo, threads int, compressor Compressor) (err error) { 

	numJobs := len(sourceFiles) 
	jobChan := make(chan r.FileInfo, numJobs) 
	resChan := make(chan error, numJobs) 
	
	if threads > numJobs {
		threads = numJobs  
	}

	for w := 1; w <= threads; w++ {
		go func(jobChan <-chan r.FileInfo, resChan chan<- error) {
			for file := range jobChan { 
				destFilePath := path.Join(destDirPath, file.Name)  
				sourceFilePath := path.Join(sourceDirPath, file.Name) 
				err := compressor(destFilePath, sourceFilePath)  
				resChan <- err 
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

// func compressEachFile_(destDirPath string, sourceDirPath string, sourceFiles []r.FileInfo, compressor Compressor) error { 
// 	for _, file := range sourceFiles { 
// 		destFilePath := path.Join(destDirPath, file.Name)  
// 		sourceFilePath := path.Join(sourceDirPath, file.Name) 
// 		err := compressor(destFilePath, sourceFilePath)   
// 		if err != nil {
// 			return err 
// 		}
// 	} 
// 	return nil 
// }

// CompressAndZipFiles concurently compresses files from a list using specified compressor, writes them into temporary directory, and then writes them into ZIP 
func CompressAndZipFiles(destFilePath string, sourceDirPath string, sourceFiles []r.FileInfo, threads int, compressor Compressor) error {

	destDirPath, _ := path.Split(destFilePath)

	tempDirPath, err := os.MkdirTemp(destDirPath, "TEMP") 
	defer os.RemoveAll(tempDirPath)
	
	if err != nil {
		return err  
	}

	err = parallelCompressEachFile(tempDirPath, sourceDirPath, sourceFiles, threads, compressor) 
	if err != nil {
		return err 
	}

	err = RawFilesToZIP(destFilePath, tempDirPath, sourceFiles) 
	if err != nil {
		return err 
	}
	
	return nil 

}

// RawFilesToZIP writes files to ZIP file 
func RawFilesToZIP(destFilePath string, sourceDirPath string, sourceFiles []r.FileInfo) error {

	destFile, err := os.Create(destFilePath) 
	if err != nil {
		return err 
	} 

	zipWriter := zip.NewWriter(destFile)  
	
	defer destFile.Close() 
	defer zipWriter.Close() 

	for _, file := range sourceFiles {

		sourceFilePath := path.Join(sourceDirPath, file.Name) 
		sourceFile, err := os.Open(sourceFilePath)
		if err != nil {
			return err 
		} 

		fileInfo, err := sourceFile.Stat() 
		if err != nil {
			return err 
		} 
		compressedSize := fileInfo.Size()


		fileWriter, err := zipWriter.CreateRaw(&zip.FileHeader{Name: file.Name, CompressedSize64: uint64(compressedSize), UncompressedSize64: uint64(file.Size)})  
		if err != nil {
			return err 
		} 

		io.Copy(fileWriter, sourceFile) 

	}

	zipWriter.Close()

	return nil 

}

