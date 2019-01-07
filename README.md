# dirsearch
Some useful functions for searching directories and returning useful information.

The information returned is a map of pathnames to `os.FileInfo` and a slice
of errors detected while scanning the directory

There are functions for

  * finding files in a single directory
  * finding files in a directory and all its subdirectories
  * finding files in a directory and those of its subdirectories that pass
    directory-specific tests
  
All of these functions allow you to pass additional checks on the entries to
restrict the results returned. You do this by passing additional arguments
which are `check.FileInfo` funcs. The simplest example of how you might use
the functions is:

```go
	results, errs := dirsearch.Find("/tmp")
```
This will return a map containing all the entries in the `/tmp` directory.

Alternatively you could find just the hidden files in your home directory with:

```go
	results, errs := dirsearch.Find(os.Getenv("HOME"),
		check.FileInfoName(check.StringHasPrefix(".")))
```

See the go documentation for more examples.
