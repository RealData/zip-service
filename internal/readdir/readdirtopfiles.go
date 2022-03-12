package readdir 

import "zip-service/internal/filelist"

// ReadDirTopFiles reads file info from a directory, filter out subdirectories, filter files by name pattern, and returns `top` largest files  
func ReadDirTopFiles(dirPath string, pattern string, top int, threads int) ([]filelist.FileInfo, error) { 

	files, err := readDir(dirPath) 
	if err != nil {
		return nil, err 
	}

	if pattern != "" {
		files, err = filelist.FilterFiles(files, pattern) 
	} 
	if err != nil {
		return nil, err 
	}
	
	return filelist.ParallelFindTop(files, top, threads), nil  
	
} 
