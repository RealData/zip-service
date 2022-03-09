package readdir 

import "os"

// ReadDir reads file info from a directory filtering out subdirectories
func readDir(dirPath string) ([]FileInfo, error) {
	
	dir, err := os.Open(dirPath) 
	if err == nil {
		defer dir.Close()  
	} else {
		return nil, err 
	}
	
	dirContent, err := dir.ReadDir(-1) 
	if err != nil {
		return nil, err 
	}

	files := make([]FileInfo, 0, len(dirContent)) 
	for _, file := range dirContent { 
		info, err := file.Info() 
		if err != nil {
			return nil, err 
		} 
		if !info.IsDir() { 
			files = append(files, FileInfo{info.Name(), info.Size()}) 
		}
	} 

	return files, nil 

}

