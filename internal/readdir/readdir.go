package readdir 

import ( 
	"path"
	"sort"
	"os"
)

type FileInfo struct {	
	Name string 
	Size int64  
}

func filterNames(files []FileInfo, pattern string) ([]FileInfo, error) {  

	filtered := make([]FileInfo, 0, len(files)) 
	for _, file := range files {
		if matched, err := path.Match(pattern, file.Name); err == nil {
			if matched {
				filtered = append(filtered, file) 
			}
		} else {
			return nil, err 
		}
	}

	return filtered, nil 

}

func ReadDir(dirStr string, pattern string, count int) ([]FileInfo, error) {
	
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
		files, err = filterNames(files, pattern) 
	} 
	if err != nil {
		return nil, err 
	}
	
	//TODO There is more effective implementation for large directories  
	sort.Slice(files, func(i, j int) bool { return files[i].Size > files[j].Size }) 
	
	if len(files) >= count {
		return files[:count], nil  
	} else {
		return files, nil 
	}

}