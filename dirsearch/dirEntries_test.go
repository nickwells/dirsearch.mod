package dirsearch_test

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/dirsearch.mod/v2/dirsearch"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestCount(t *testing.T) {
	goodDir := path.Join("testdata", "IsADirectory")
	badDir := path.Join("testdata", "NoSuchDir")
	fileName := path.Join("testdata", "IsAFile")

	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		checks               []check.ValCk[os.FileInfo]
		dirChecks            []check.ValCk[os.FileInfo]
		maxDepth             int
		countExp             int
		countExpRecurse      int
		countExpRecursePrune int
		dirName              string
	}{
		{
			ID:      testhelper.MkID("Bad directory: " + badDir),
			ExpErr:  testhelper.MkExpErr(badDir, "no such file or directory"),
			dirName: badDir,
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
			dirName:              goodDir,
		},
		{
			ID:     testhelper.MkID("All files"),
			checks: []check.ValCk[os.FileInfo]{check.FileInfoIsRegular},

			maxDepth:             1,
			countExp:             3,
			countExpRecurse:      6,
			countExpRecursePrune: 6,
			dirName:              goodDir,
		},
		{
			ID: testhelper.MkID("All files - no hidden files (leading '.')"),
			checks: []check.ValCk[os.FileInfo]{
				check.FileInfoIsRegular,
				check.Not(
					check.FileInfoName(check.StringHasPrefix[string](".")), ""),
			},
			dirChecks: []check.ValCk[os.FileInfo]{
				check.Not(
					check.FileInfoName(check.StringHasPrefix[string](".")), ""),
			},
			maxDepth:             -1,
			countExp:             2,
			countExpRecurse:      5,
			countExpRecursePrune: 4,
			dirName:              goodDir,
		},
	}

	for _, tc := range testCases {
		id := fmt.Sprintf("%s - Count(%q, ...)",
			tc.IDStr(), tc.dirName)
		n, errs := dirsearch.Count(tc.dirName, tc.checks...)
		testhelper.CheckExpErrWithID(t, id, errFromErrs(errs), tc)
		testhelper.DiffInt(t, id, "count", n, tc.countExp)

		id = fmt.Sprintf("%s - CountRecurse(%q, ...)",
			tc.IDStr(), tc.dirName)
		n, errs = dirsearch.CountRecurse(tc.dirName, tc.checks...)
		testhelper.CheckExpErrWithID(t, id, errFromErrs(errs), tc)
		testhelper.DiffInt(t, id, "count", n, tc.countExpRecurse)

		id = fmt.Sprintf("%s - CountRecursePrune(%q, ...)",
			tc.IDStr(), tc.dirName)
		n, errs = dirsearch.CountRecursePrune(tc.dirName,
			tc.maxDepth, tc.dirChecks,
			tc.checks...)
		testhelper.CheckExpErrWithID(t, id, errFromErrs(errs), tc)
		testhelper.DiffInt(t, id, "count", n, tc.countExpRecursePrune)
	}
}

// errFromErrs returns the first error from the slice of errors if the slice
// is non-empty and nil otherwise
func errFromErrs(errs []error) error {
	if len(errs) == 0 {
		return nil
	}

	return errs[0]
}
