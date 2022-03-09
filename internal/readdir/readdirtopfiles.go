package readdir 

// ReadDirTopFiles reads file info from a directory, filter out subdirectories, filter files by name pattern, and returns `top` largest files  
func ReadDirTopFiles(dirPath string, pattern string, top int, threads int) ([]FileInfo, error) { 

	files, err := readDir(dirPath) 
	if err != nil {
		return nil, err 
	}

	if pattern != "" {
		files, err = filterFiles(files, pattern) 
	} 
	if err != nil {
		return nil, err 
	}
	
	return parallelFindTop(files, top, threads), nil 
	
} 
