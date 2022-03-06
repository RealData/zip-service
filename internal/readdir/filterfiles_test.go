package readdir 

import "testing" 

func TestFilterFiles(t *testing.T) {

	t.Run("file FILE pattern FILE", func(t *testing.T) {
		files := []FileInfo{{"FILE", 10}}
		filtered, err := filterFiles(files, "FILE") 
		if err != nil {
			t.Error(err) 
		} 
		if len(filtered) == 1 {
			if filtered[0].Name != "FILE" {
				t.Error("Name should be 'FILE'")
			}
		} else {
			t.Error("Should return one file")
		}
	})

	t.Run("file FILE pattern *F*", func(t *testing.T) {
		files := []FileInfo{{"FILE", 10}}
		filtered, err := filterFiles(files, "*F*") 
		if err != nil {
			t.Error(err) 
		} 
		if len(filtered) == 1 {
			if filtered[0].Name != "FILE" {
				t.Error("Name should be 'FILE'")
			}
		} else {
			t.Error("Should return one file")
		}
	})
	
	t.Run("file FILE pattern F", func(t *testing.T) {
		files := []FileInfo{{"F", 10}}
		filtered, err := filterFiles(files, "FILE") 
		if err != nil {
			t.Error(err) 
		} 
		if len(filtered) > 0 {
			t.Error("Should return no files")
		}
	})

}
