package compressextract 
import (
	"archive/zip"
	"io" 
	"os" 
	"path"
	"zip-service/internal/filelist"
)

// ExtractFileInfo returns list of file info for files in a ZIP file 
func ExtractFileInfo(sourceFilePath string) ([]filelist.FileInfo, error)  {

	reader, err := zip.OpenReader(sourceFilePath) 
	if err != nil {
		return nil, err  
	} 

	var files []filelist.FileInfo 

	for _, f := range reader.File {
		files = append(files, filelist.FileInfo{Name: f.Name, Size: int64(f.CompressedSize64)})
	}

	return files, nil 

}

// zipToRawFiles reads files from ZIP file 
func zipToRawFiles(destDirPath string, sourceFilePath string) ([]filelist.FileInfo, error) {

	zipReader, err := zip.OpenReader(sourceFilePath) 
	if err != nil {
		return nil, err 
	} 
	defer zipReader.Close() 
	
	var files []filelist.FileInfo  

	for _, f := range zipReader.File {
		reader, err := f.OpenRaw() 
		writer, err := os.Create(path.Join(destDirPath, f.Name)) 
		_, err = io.Copy(writer, reader) 
		if err != nil {
			return nil, err 
		} 
		files = append(files, filelist.FileInfo{f.Name, int64(f.CompressedSize64)})
		writer.Close()
	}

	return files, nil 

}

func parallelExtractEachFile(destDirPath string, sourceDirPath string, sourceFiles []filelist.FileInfo, threads int, extractor Extractor) (err error) { 

	numJobs := len(sourceFiles) 
	jobChan := make(chan filelist.FileInfo, numJobs) 
	resChan := make(chan error, numJobs) 
	
	if threads > numJobs {
		threads = numJobs  
	}

	for w := 1; w <= threads; w++ {
		go func(jobChan <-chan filelist.FileInfo, resChan chan<- error) {
			for file := range jobChan { 
				destFilePath := path.Join(destDirPath, file.Name)  
				sourceFilePath := path.Join(sourceDirPath, file.Name) 
				err := extractor(destFilePath, sourceFilePath)  
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

	return nil 

}

// extractEachFile extracts files from a list using specified compressor 
func extractEachFile(destDirPath string, sourceDirPath string, sourceFiles []filelist.FileInfo, extractor Extractor) error { 
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
func UnzipAndExtractFiles(destDirPath string, sourceFilePath string, threads int, extractor Extractor) error {
	
	tempDirPath, err := os.MkdirTemp(destDirPath, "TEMP") 
	if err != nil {
	 	return err 
	}
	defer os.RemoveAll(tempDirPath)

	files, err := zipToRawFiles(tempDirPath, sourceFilePath)  
	if err != nil {
		return err 
    }

	err = parallelExtractEachFile(destDirPath, tempDirPath, files, threads, extractor) 
	if err != nil {
		return err 
	}

	return nil  

}