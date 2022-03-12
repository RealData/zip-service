package readdir 

import (
	"os" 
	"zip-service/internal/filelist" 
)


// ReadDir reads file info from a directory filtering out subdirectories
func readDir(dirPath string) ([]filelist.FileInfo, error) {
	
	dir, err := os.Open(dirPath) 
	if err != nil {
		return nil, err 
	}
	defer dir.Close() 
	
	dirContent, err := dir.ReadDir(-1) 
	if err != nil {
		return nil, err 
	}

	files := make([]filelist.FileInfo, 0, len(dirContent)) 
	for _, file := range dirContent { 
		info, err := file.Info() 
		if err != nil {
			return nil, err 
		} 
		if !info.IsDir() { 
			files = append(files, filelist.FileInfo{info.Name(), info.Size()}) 
		}
	} 

	return files, nil 

}

