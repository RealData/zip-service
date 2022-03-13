package readdir

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
	"zip-service/internal/filelist"
)

type File struct {
	Name    string
	Content string
}

func TestReadDirTopFiles(t *testing.T) {

	var tests = []struct {
		scenario string
		files    []File
		filtered []filelist.FileInfo
		top      int
		threads  int
		pattern  string
	}{
		{
			"5 files, top=10, threads=1, pattern=''",
			[]File{{"file1", "LongContent"}, {"file2", "file2"}, {"file3", ""}, {"file4", "HH"}, {"file5", "Content"}},
			[]filelist.FileInfo{{"file1", 11}, {"file5", 7}, {"file2", 5}, {"file4", 2}, {"file3", 0}},
			10,
			1,
			"",
		},
		{
			"5 files, top=10, threads=2, pattern=''",
			[]File{{"file1", "LongContent"}, {"file2", "file2"}, {"file3", ""}, {"file4", "HH"}, {"file5", "Content"}},
			[]filelist.FileInfo{{"file1", 11}, {"file5", 7}, {"file2", 5}, {"file4", 2}, {"file3", 0}},
			10,
			2,
			"",
		},
		{
			"5 files, top=2, threads=1, pattern=''",
			[]File{{"file1", "LongContent"}, {"file2", "file2"}, {"file3", ""}, {"file4", "HH"}, {"file5", "Content"}},
			[]filelist.FileInfo{{"file1", 11}, {"file5", 7}},
			2,
			1,
			"",
		},
		{
			"5 files, top=2, threads=2, pattern=''",
			[]File{{"file1", "LongContent"}, {"file2", "file2"}, {"file3", ""}, {"file4", "HH"}, {"file5", "Content"}},
			[]filelist.FileInfo{{"file1", 11}, {"file5", 7}},
			2,
			2,
			"",
		},
	}

	dir := "TEST"
	defer os.RemoveAll(dir)

	for _, test := range tests {

		err := os.Mkdir(dir, 0777)

		if err != nil {
			t.Error("Cannot create temporary directory")
		}

		t.Run(test.scenario, func(t *testing.T) {

			for _, file := range test.files {
				filePath := filepath.Join(dir, file.Name)
				if err := os.WriteFile(filePath, []byte(file.Content), 0777); err != nil {
					t.Error("Cannot create file")
				}
			}

			res, err := ReadDirTopFiles(dir, test.pattern, test.top, test.threads)
			if err != nil {
				t.Error(err)
			}
			if len(test.files) >= test.top && len(res) != test.top {
				t.Errorf("Should return top=%d files, found %d", test.top, len(res))
			}
			sort.Slice(res, func(i int, j int) bool { return res[i].Size > res[j].Size })
			for i, file := range test.filtered {
				if file.Name != res[i].Name {
					t.Errorf("Name of file #%d should be %s, found %s", i, file.Name, res[i].Name)
				}
			}
		})

		os.RemoveAll(dir)

	}

}
