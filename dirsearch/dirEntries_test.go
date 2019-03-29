package dirsearch_test

import (
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
		testhelper.ID
		testhelper.ExpErr
		checks               []check.FileInfo
		dirChecks            []check.FileInfo
		maxDepth             int
		countExp             int
		countExpRecurse      int
		countExpRecursePrune int
		dirName              string
	}{
		{
			ID: testhelper.MkID("bad directory: " + badDirName),
			ExpErr: testhelper.MkExpErr(badDirName,
				"no such file or directory"),
			dirName: badDirName,
		},
		{
			ID:      testhelper.MkID("Not a directory: " + fileName),
			ExpErr:  testhelper.MkExpErr("not a directory"),
			dirName: fileName,
		},
		{
			ID:                   testhelper.MkID("all entries"),
			maxDepth:             0,
			countExp:             5,
			countExpRecurse:      8,
			countExpRecursePrune: 5,
			dirName:              goodDirName,
		},
		{
			ID: testhelper.MkID("all files"),
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
			ID: testhelper.MkID("all files - ignore hidden files (leading '.')"),
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

	for _, tc := range testCases {
		n, errs := dirsearch.Count(tc.dirName, tc.checks...)
		var err error
		if len(errs) > 0 {
			err = errs[0]
		}
		testhelper.CheckExpErr(t, err, tc)

		if n != tc.countExp {
			t.Log(tc.IDStr())
			t.Logf("\t: Count() in dir: %s\n", tc.dirName)
			t.Errorf("\t: expected count: %d got: %d\n", tc.countExp, n)
		}

		n, errs = dirsearch.CountRecurse(tc.dirName, tc.checks...)
		err = nil
		if len(errs) > 0 {
			err = errs[0]
		}
		testhelper.CheckExpErr(t, err, tc)

		if n != tc.countExpRecurse {
			t.Log(tc.IDStr())
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
		testhelper.CheckExpErr(t, err, tc)

		if n != tc.countExpRecursePrune {
			t.Log(tc.IDStr())
			t.Logf("\t: CountRecursePrune() in dir: %s\n", tc.dirName)
			t.Errorf("\t: expected count: %d got: %d\n",
				tc.countExpRecursePrune, n)
		}
	}
}
