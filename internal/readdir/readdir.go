package readdir 

import "os"

// ReadDir reads file info from a directory, filter out subdirectories, filter files by name pattern, and returns `top` largest files  
func ReadDir(dirStr string, pattern string, top int, threads int) ([]FileInfo, error) {
	
	dir, err := os.Open(dirStr) 
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

	if pattern != "" {
		files, err = filterFiles(files, pattern) 
	} 
	if err != nil {
		return nil, err 
	}
	
	files = parallelFindTop(files, top, threads)
	
	if len(files) >= top {
		return files[:top], nil  
	} else {
		return files, nil 
	}

}