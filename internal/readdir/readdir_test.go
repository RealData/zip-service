package readdir 

import ( 
	"os"
	"path/filepath"
	"testing" 
)

func TestReadDir(t *testing.T) {
	
	t.Run("One file", func(t *testing.T) {
		dirStr := "TEST" 
		dir := dirStr 
		err :=  os.Mkdir(dir, 0777) 
		// dir, err := os.MkdirTemp(".", dirStr) 
		if err != nil {
			t.Error("Cannot create temporary directory")
		} else {
			defer os.RemoveAll(dir) 
		}

		file := filepath.Join(dir, "file") 
		if err := os.WriteFile(file, []byte("Content"), 0777); err != nil {
			t.Error("Cannot create file") 
		}

		files, err := ReadDir(dir, "", 10, 1) 
		if err != nil {
			t.Error(err)  
		}
		if len(files) != 1 { 
			t.Error("Length should be equal to 1") 
		} else {
			if files[0].Name != "file" {
				t.Error("File name should be 'file'") 
			}
			if files[0].Size != 7 {
				t.Error("Incorrect size of the file")
			}
		}
	})

	t.Run("No files", func(t *testing.T) {
		dirStr := "TEST" 
		dir := dirStr 
		err :=  os.Mkdir(dir, 0777) 
		// dir, err := os.MkdirTemp(".", dirStr) 
		if err != nil {
			t.Error("Cannot create temporary directory")
		} else {
			defer os.RemoveAll(dir) 
		}

		files, err := ReadDir(dir, "", 10, 1) 
		if err != nil {
			t.Error(err)  
		}
		if len(files) > 0 { 
			t.Error("Should return no files")  
		} 	
	})

}