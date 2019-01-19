package msync

import "testing"

func in(s string, sa []string) bool {
	for _, item := range sa {
		if s == item {
			return true
		}
	}
	return false
}

func TestWalk(t *testing.T) {
	testingFiles := []string{
		"../test/sync/dir/file_dir.flac",
		"../test/sync/full.mp3",
		"../test/sync/album/1.mp3",
		"../test/sync/album/2.mp3",
		"../test/sync/album/3.mp3"}

	files, err := getFiles("../test")
	if err != nil {
		t.Error(err)
	}

	var valid int
	for _, file := range files {
		if !in(file, testingFiles) {
			t.Errorf("%v shouldn't be here", file)
		} else {
			valid++
		}
	}

	if valid != len(testingFiles) {
		t.Errorf("Not all files are validated")
	}

}
