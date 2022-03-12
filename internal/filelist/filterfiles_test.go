package filelist 

import "testing" 

func TestFilterFiles(t *testing.T) {

	var tests = []struct { 
		scenario string 
		files []FileInfo  
		pattern string 
		filtered []FileInfo 
	} {
		{"file FILE pattern FILE", []FileInfo{{"FILE", 10}}, "FILE", []FileInfo{{"FILE", 10}}}, 
		{"file FILE pattern F*", []FileInfo{{"FILE", 10}}, "F*", []FileInfo{{"FILE", 10}}}, 
		{"file FILE pattern F", []FileInfo{{"FILE", 10}}, "F", []FileInfo{}}, 
		{"file FILE,FILE2 pattern FILE", []FileInfo{{"FILE", 10}, {"FILE2", 10}}, "FILE", []FileInfo{{"FILE", 10}}}, 
		{"file FILE,FILE2 pattern F*", []FileInfo{{"FILE", 10}, {"FILE2", 10}}, "F*", []FileInfo{{"FILE", 10}, {"FILE2", 10}}},  
	}

	for _, test := range tests { 
		t.Run(test.scenario, func(t *testing.T) {
			filtered, err := filterFiles(test.files, test.pattern) 
			if err != nil {
				t.Error(err) 
			}
			if len(filtered) != len(test.filtered) {
				t.Errorf("Length of the result should be %d, found %d", len(test.filtered), len(filtered)) 
				return  
			}
			for i, file := range test.filtered {
				if file.Name != filtered[i].Name {
					t.Errorf("Name of file #%d should be %s, found %s", i, file.Name, filtered[i].Name)
				}
			} 
		})
	}
	
}
