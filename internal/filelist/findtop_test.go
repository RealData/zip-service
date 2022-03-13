package filelist

import (
	"sort"
	"testing"
)

func TestFindTop(t *testing.T) {

	var tests = []struct {
		scenario string
		files    []FileInfo
		filtered []FileInfo
		top      int
	}{
		{
			"5 files, top=10",
			[]FileInfo{{"file1", 5}, {"file2", 6}, {"file3", 10}, {"file4", 16}, {"file5", 17}},
			[]FileInfo{{"file5", 17}, {"file4", 16}, {"file3", 10}, {"file2", 6}, {"file1", 5}},
			10,
		},
		{
			"5 files, top=5",
			[]FileInfo{{"file1", 5}, {"file2", 6}, {"file3", 10}, {"file4", 16}, {"file5", 17}},
			[]FileInfo{{"file5", 17}, {"file4", 16}, {"file3", 10}, {"file2", 6}, {"file1", 5}},
			10,
		},
		{
			"5 files, top=2",
			[]FileInfo{{"file1", 5}, {"file2", 6}, {"file3", 10}, {"file4", 16}, {"file5", 17}},
			[]FileInfo{{"file5", 17}, {"file4", 16}},
			10,
		},
		{
			"5 files, top=0",
			[]FileInfo{{"file1", 5}, {"file2", 6}, {"file3", 10}, {"file4", 16}, {"file5", 17}},
			[]FileInfo{},
			0,
		},
		{
			"0 files, top=5",
			[]FileInfo{},
			[]FileInfo{},
			5,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			res := findTop(test.files, test.top)
			sort.Slice(res, func(i int, j int) bool { return res[i].Size > res[j].Size })
			if len(test.files) >= test.top && len(res) != test.top {
				t.Errorf("Should return top=%d files, found %d", test.top, len(res))
			}
			for i, file := range test.filtered {
				if file.Name != res[i].Name {
					t.Errorf("Name of file #%d should be %s, found %s", i, file.Name, res[i].Name)
				}
			}
		})
	}

}

func TestParallelFindTop(t *testing.T) {

	var tests = []struct {
		scenario string
		files    []FileInfo
		filtered []FileInfo
		top      int
		threads  int
	}{
		{
			"5 files, top=10",
			[]FileInfo{{"file1", 5}, {"file2", 6}, {"file3", 10}, {"file4", 16}, {"file5", 17}},
			[]FileInfo{{"file5", 17}, {"file4", 16}, {"file3", 10}, {"file2", 6}, {"file1", 5}},
			10,
			2,
		},
		{
			"5 files, top=5",
			[]FileInfo{{"file1", 5}, {"file2", 6}, {"file3", 10}, {"file4", 16}, {"file5", 17}},
			[]FileInfo{{"file5", 17}, {"file4", 16}, {"file3", 10}, {"file2", 6}, {"file1", 5}},
			10,
			2,
		},
		{
			"5 files, top=2",
			[]FileInfo{{"file1", 5}, {"file2", 6}, {"file3", 10}, {"file4", 16}, {"file5", 17}},
			[]FileInfo{{"file5", 17}, {"file4", 16}},
			10,
			2,
		},
		{
			"5 files, top=0",
			[]FileInfo{{"file1", 5}, {"file2", 6}, {"file3", 10}, {"file4", 16}, {"file5", 17}},
			[]FileInfo{},
			0,
			2,
		},
		{
			"0 files, top=5",
			[]FileInfo{},
			[]FileInfo{},
			5,
			2,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			res := ParallelFindTop(test.files, test.top, test.threads)
			sort.Slice(res, func(i int, j int) bool { return res[i].Size > res[j].Size })
			if len(test.files) >= test.top && len(res) != test.top {
				t.Errorf("Should return top=%d files, found %d", test.top, len(res))
			}
			for i, file := range test.filtered {
				if file.Name != res[i].Name {
					t.Errorf("Name of file #%d should be %s, found %s", i, file.Name, res[i].Name)
				}
			}
		})
	}

}
