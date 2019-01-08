package dirsearch_test

import (
	"fmt"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/dirsearch.mod/dirsearch"
)

func ExampleFind() {
	info, errs := dirsearch.Find("testdata/examples/dir1")
	if len(errs) != 0 {
		fmt.Println("Unexpected errors")
		for _, err := range errs {
			fmt.Println("\t", err)
		}
		return
	}
	for k, v := range info {
		sizeStr := "Non-Empty"
		if v.Size() == 0 {
			sizeStr = "Empty"
		}
		fmt.Println("file:", k, "=", sizeStr)
	}
	// Unordered output:
	// file: testdata/examples/dir1/non-empty-file1 = Non-Empty
	// file: testdata/examples/dir1/.hidden-subdir1 = Non-Empty
	// file: testdata/examples/dir1/empty-file1 = Empty
	// file: testdata/examples/dir1/.non-empty-hidden-file1 = Non-Empty
	// file: testdata/examples/dir1/subdir1 = Non-Empty
}

// ExampleFind_withChecks demonstrates use of the Find function with checks
// supplied to return only non-empty regular files
func ExampleFind_withChecks() {
	info, errs := dirsearch.Find("testdata/examples/dir1",
		check.FileInfoSize(check.Int64GT(0)),
		check.FileInfoIsRegular)
	if len(errs) != 0 {
		fmt.Println("Unexpected errors")
		for _, err := range errs {
			fmt.Println("\t", err)
		}
		return
	}
	for k, v := range info {
		sizeStr := "Non-Empty"
		if v.Size() == 0 {
			sizeStr = "Empty"
		}
		fmt.Println("file:", k, "=", sizeStr)
	}
	// Unordered output:
	// file: testdata/examples/dir1/non-empty-file1 = Non-Empty
	// file: testdata/examples/dir1/.non-empty-hidden-file1 = Non-Empty
}

// ExampleFindRecurse_withChecks demonstrates use of the FindRecurse function
// with checks supplied to return only non-empty regular files
func ExampleFindRecurse_withChecks() {
	info, errs := dirsearch.FindRecurse("testdata/examples/dir1",
		check.FileInfoSize(check.Int64GT(0)),
		check.FileInfoIsRegular)
	if len(errs) != 0 {
		fmt.Println("Unexpected errors")
		for _, err := range errs {
			fmt.Println("\t", err)
		}
		return
	}
	for k, v := range info {
		sizeStr := "Non-Empty"
		if v.Size() == 0 {
			sizeStr = "Empty"
		}
		fmt.Println("file:", k, "=", sizeStr)
	}
	// Unordered output:
	// file: testdata/examples/dir1/non-empty-file1 = Non-Empty
	// file: testdata/examples/dir1/subdir1/non-empty-file = Non-Empty
	// file: testdata/examples/dir1/subdir1/subsubdir1/non-empty-file1 = Non-Empty
	// file: testdata/examples/dir1/.hidden-subdir1/non-empty-file = Non-Empty
	// file: testdata/examples/dir1/.non-empty-hidden-file1 = Non-Empty
}

// ExampleFindRecursePrune_withChecks demonstrates use of the
// FindRecursePrune function with checks supplied to return only non-empty
// regular files and a slice of directory checks provided to prevent descent
// into hidden directories (those with a name starting with a '.')
func ExampleFindRecursePrune_withChecks() {
	info, errs := dirsearch.FindRecursePrune("testdata/examples/dir1",
		-1, []check.FileInfo{
			check.FileInfoName(
				check.StringNot(
					check.StringHasPrefix("."),
					"no leading '.'")),
		},
		check.FileInfoSize(check.Int64GT(0)),
		check.FileInfoIsRegular)
	if len(errs) != 0 {
		fmt.Println("Unexpected errors")
		for _, err := range errs {
			fmt.Println("\t", err)
		}
		return
	}
	for k, v := range info {
		sizeStr := "Non-Empty"
		if v.Size() == 0 {
			sizeStr = "Empty"
		}
		fmt.Println("file:", k, "=", sizeStr)
	}
	// Unordered output:
	// file: testdata/examples/dir1/non-empty-file1 = Non-Empty
	// file: testdata/examples/dir1/.non-empty-hidden-file1 = Non-Empty
	// file: testdata/examples/dir1/subdir1/non-empty-file = Non-Empty
	// file: testdata/examples/dir1/subdir1/subsubdir1/non-empty-file1 = Non-Empty
}
