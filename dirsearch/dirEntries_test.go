package dirsearch_test

import (
	"fmt"
	"path"
	"testing"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/dirsearch.mod/dirsearch"
	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestCount(t *testing.T) {
	goodDirName := path.Join("testdata", "IsADirectory")
	badDirName := path.Join("testdata", "NoSuchDir")
	fileName := path.Join("testdata", "IsAFile")

	testCases := []struct {
		name                 string
		checks               []check.FileInfo
		dirChecks            []check.FileInfo
		maxDepth             int
		errExpected          bool
		countExp             int
		countExpRecurse      int
		countExpRecursePrune int
		dirName              string
	}{
		{
			name:        "bad directory: " + badDirName,
			errExpected: true,
			dirName:     badDirName,
		},
		{
			name:        "Not a directory: " + fileName,
			errExpected: true,
			dirName:     fileName,
		},
		{
			name:                 "all entries",
			maxDepth:             0,
			countExp:             5,
			countExpRecurse:      8,
			countExpRecursePrune: 5,
			dirName:              goodDirName,
		},
		{
			name: "all files",
			checks: []check.FileInfo{
				check.FileInfoIsRegular,
			},
			maxDepth:             1,
			countExp:             3,
			countExpRecurse:      6,
			countExpRecursePrune: 6,
			dirName:              goodDirName,
		},
		{
			name: "all files - ignore hidden files (leading '.')",
			checks: []check.FileInfo{
				check.FileInfoIsRegular,
				check.FileInfoNot(
					check.FileInfoName(
						check.StringHasPrefix(".")),
					""),
			},
			dirChecks: []check.FileInfo{
				check.FileInfoNot(
					check.FileInfoName(
						check.StringHasPrefix(".")),
					""),
			},
			maxDepth:             -1,
			countExp:             2,
			countExpRecurse:      5,
			countExpRecursePrune: 4,
			dirName:              goodDirName,
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s", i, tc.name)

		n, errs := dirsearch.Count(tc.dirName, tc.checks...)
		var err error
		if len(errs) > 0 {
			err = errs[0]
		}
		testhelper.CheckError(t, tcID, err, tc.errExpected, []string{})

		if n != tc.countExp {
			t.Logf("test %d: %s :\n", i, tc.name)
			t.Logf("\t: Count() in dir: %s\n", tc.dirName)
			t.Errorf("\t: expected count: %d got: %d\n", tc.countExp, n)
		}

		n, errs = dirsearch.CountRecurse(tc.dirName, tc.checks...)
		err = nil
		if len(errs) > 0 {
			err = errs[0]
		}
		testhelper.CheckError(t, tcID, err, tc.errExpected, []string{})

		if n != tc.countExpRecurse {
			t.Logf("test %d: %s :\n", i, tc.name)
			t.Logf("\t: CountRecurse() in dir: %s\n", tc.dirName)
			t.Errorf("\t: expected count: %d got: %d\n", tc.countExpRecurse, n)
		}

		n, errs = dirsearch.CountRecursePrune(
			tc.dirName,
			tc.maxDepth,
			tc.dirChecks,
			tc.checks...)
		err = nil
		if len(errs) > 0 {
			err = errs[0]
		}
		testhelper.CheckError(t, tcID, err, tc.errExpected, []string{})

		if n != tc.countExpRecursePrune {
			t.Logf("test %d: %s :\n", i, tc.name)
			t.Logf("\t: CountRecursePrune() in dir: %s\n", tc.dirName)
			t.Errorf("\t: expected count: %d got: %d\n",
				tc.countExpRecursePrune, n)
		}
	}
}
