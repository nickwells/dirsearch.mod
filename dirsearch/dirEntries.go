package dirsearch

import (
	"os"
	"path/filepath"

	"github.com/nickwells/check.mod/v2/check"
)

// getDirInfo reads the FileInfo for every entry in the directory. It uses
// the os.Readdir func internally and so the FileInfo details do not follow
// symlinks
func getDirInfo(name string) ([]os.FileInfo, error) {
	// nolint: gosec
	d, err := os.Open(name)
	if err != nil {
		return []os.FileInfo{}, err
	}

	defer d.Close() // nolint: errcheck

	return d.Readdir(0)
}

// passesChecks returns true if all the checks return a nil error, false
// otherwise
func passesChecks(fi os.FileInfo, checks []check.ValCk[os.FileInfo]) bool {
	for _, c := range checks {
		if c(fi) != nil {
			return false
		}
	}
	return true
}

// dirPassesChecks checks that the directory passes all the checks given - on
// directory depth and the FileInfo checks. It also tests that the names are
// neither of ".", or ".."; strictly this is unnecessary as os.Readdir
// already excludes these entries but this is a bug (which would break the Go
// compatibility promise if it was fixed - sigh)
func dirPassesChecks(fi os.FileInfo, depth, maxDepth int,
	checks []check.ValCk[os.FileInfo],
) bool {
	subDirName := fi.Name()
	if subDirName == "." {
		return false
	}
	if subDirName == ".." {
		return false
	}
	if maxDepth >= 0 && (depth+1) > maxDepth {
		return false
	}
	return passesChecks(fi, checks)
}

// find ...
func find(dirName string, depth, maxDepth int,
	dirChecks, checks []check.ValCk[os.FileInfo],
) (map[string]os.FileInfo, []error) {
	allInfo, err := getDirInfo(dirName)
	if err != nil {
		return map[string]os.FileInfo{}, []error{err}
	}

	info := make(map[string]os.FileInfo)
	errors := make([]error, 0)

	for _, fi := range allInfo {
		if passesChecks(fi, checks) {
			info[filepath.Join(dirName, fi.Name())] = fi
		}

		if fi.IsDir() {
			if dirPassesChecks(fi, depth, maxDepth, dirChecks) {
				subDirInfo, subDirErrors := find(
					filepath.Join(dirName, fi.Name()),
					depth+1, maxDepth,
					dirChecks, checks)
				errors = append(errors, subDirErrors...)
				for k, v := range subDirInfo {
					info[k] = v
				}
			}
		}
	}
	return info, errors
}

// FindRecurse reads the directory and returns the FileInfo details for each
// entry for which all the checks return true. Any errors are also
// returned. The details are returned in a map of file names (including the
// directory prefixes) to FileInfo. This differs from Find(...) in that any
// entries in the directory which are themselves directories (except for "."
// and "..") are descended into and the search continues in that
// sub-directory
func FindRecurse(dirName string, checks ...check.ValCk[os.FileInfo],
) (map[string]os.FileInfo, []error) {
	return find(dirName,
		0, -1,
		[]check.ValCk[os.FileInfo]{}, checks)
}

// FindRecursePrune reads the directory and returns the FileInfo details for
// each entry for which all the checks return true. Any errors are also
// returned. The details are returned in a map of file names (including the
// directory prefixes) to FileInfo. This differs from FindRecurse(...) in
// that the entries in the directory which are themselves directories are
// checked to see that they are not too deep and that they pass the dirChecks
// and skipped if not. Pass a maxDepth < 0 to ignore the depth of descent
// down the directory tree.
func FindRecursePrune(dirName string, maxDepth int,
	dirChecks []check.ValCk[os.FileInfo], checks ...check.ValCk[os.FileInfo],
) (map[string]os.FileInfo, []error) {
	return find(dirName,
		0, maxDepth,
		dirChecks, checks)
}

// Find reads the directory and returns the FileInfo details for each entry
// for which all the checks return true. Any errors are also returned. The
// details are returned in a map of file names (including the directory
// prefix) to FileInfo
func Find(dirName string,
	checks ...check.ValCk[os.FileInfo],
) (map[string]os.FileInfo, []error) {
	return find(dirName, 0, 0, []check.ValCk[os.FileInfo]{}, checks)
}

// Count returns the number of entries in the directory that match
// the supplied checks (if any) and any errors detected
func Count(name string, checks ...check.ValCk[os.FileInfo]) (int, []error) {
	info, err := Find(name, checks...)
	return len(info), err
}

// CountRecurse returns the number of entries in the directory and any
// sub-directories that match the supplied checks (if any) and any errors
// detected
func CountRecurse(name string, checks ...check.ValCk[os.FileInfo],
) (int, []error) {
	info, err := FindRecurse(name, checks...)
	return len(info), err
}

// CountRecursePrune returns the number of entries in the directory and any
// sub-directories that match the supplied checks (if any) and any errors
// detected
func CountRecursePrune(name string, maxDepth int,
	dirChecks []check.ValCk[os.FileInfo], checks ...check.ValCk[os.FileInfo],
) (int, []error) {
	info, err := FindRecursePrune(name, maxDepth, dirChecks, checks...)
	return len(info), err
}
