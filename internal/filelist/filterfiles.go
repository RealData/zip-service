package filelist

import "path"

// FilterFiles filter a list of FileInfo in accordance with the pattern
func FilterFiles(files []FileInfo, pattern string) ([]FileInfo, error) {

	filtered := make([]FileInfo, 0, len(files))
	for _, file := range files {
		if matched, err := path.Match(pattern, file.Name); err != nil {
			return nil, err
		} else if matched {
			filtered = append(filtered, file)
		}
	}

	return filtered, nil

}
