package readdir 

import "testing" 


func TestFindTop(t *testing.T) {

	files := []FileInfo{
		FileInfo{"file1", 5}, 
		FileInfo{"file2", 6}, 
		FileInfo{"file3", 10},  
		FileInfo{"file4", 16},  
		FileInfo{"file5", 17}, 
	}

	top := 10 

	res := findTop(files, top) 

	if len(files) >= top && len(res) != top { 
		t.Errorf("Should return top=%d files", top)
	}

	if len(res) > 0  {
	
		if res[0].Name != "file5" || res[0].Size != 17 {
			t.Error("Incorrect top file")
		}
	}

}

func TestParallelFindTop(t *testing.T) {

	files := []FileInfo{
		FileInfo{"file1", 5}, 
		FileInfo{"file2", 6}, 
		FileInfo{"file3", 10},  
		FileInfo{"file4", 16},  
		FileInfo{"file5", 17}, 
	}

	threads := 2 
	top := 2 

	res := parallelFindTop(files, top, threads) 

	if len(files) >= top && len(res) != top { 
		t.Errorf("Should return top=%d files", top)
	}

	if len(res) > 0  {
	
		if res[0].Name != "file5" || res[0].Size != 17 {
			t.Error("Incorrect top file")
		}
	}

}
