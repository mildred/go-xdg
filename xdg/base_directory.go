// (c) 2014 John R. Lenton. See LICENSE.

package xdg

import (
	"os"
	"os/user"
	"path/filepath"
)

// An XDGDir holds configuration for and can be used to access the
// XDG-specified base directories relative to which user-specific files of a
// given type should be stored.
type XDGDir struct {
	homeEnv  string
	homeDflt string
	dirsEnv  string
	dirsDflt string
}

var (
	Data   *XDGDir // for data files.
	Config *XDGDir // for configuration files.
	Cache  *XDGDir // for non-essential data files.
)

func init() {
	// do this here to make the docs nicer
	Data = &XDGDir{"XDG_DATA_HOME", ".local/share", "XDG_DATA_DIRS", "/usr/local/share:/usr/share"}
	Config = &XDGDir{"XDG_CONFIG_HOME", ".config", "XDG_CONFIG_DIRS", "/etc/xdg"}
	Cache = &XDGDir{"XDG_CACHE_HOME", ".cache", "", ""}
}

// Home gets the path to the given user-specific XDG directory, as specified
// (or not) by the user's environment.
func (x *XDGDir) Home() string {
	dir := os.Getenv(x.homeEnv)
	if dir != "" {
		return dir
	}
	home := os.Getenv("HOME")
	if home == "" {
		user, err := user.Current()
		if err != nil {
			panic("unable to determine $HOME")
		}
		home = user.HomeDir
	}
	return filepath.Join(home, x.homeDflt)
}

// Dirs returns the preference-ordered set of base directories to search for
// files of the given type, starting with the user-specific one, as specified
// (or not) by the user's environment.
func (x *XDGDir) Dirs() []string {
	dirs := []string{x.Home()}
	if x.dirsEnv != "" {
		xtra := os.Getenv(x.dirsEnv)
		if xtra == "" {
			xtra = x.dirsDflt
		}
		for _, path := range filepath.SplitList(xtra) {
			if path != "" {
				dirs = append(dirs, path)
			}
		}
	}
	return dirs
}

// FirstExisting returns the path to the  first of the given resource that
// exists in the XDG directories. If none exist, error is set appropriately.
//
// A resource is usually an application name, followed by a filename. Nothing
// in the standard nor in this library limits you to that.
func (x *XDGDir) FirstExisting(resourcePath ...string) (fullPath string, err error) {
	var firstError error = nil
	for _, path := range x.Dirs() {
		name := filepath.Join(path, filepath.Join(resourcePath...))
		_, err = os.Stat(name)
		if err == nil {
			return name, nil
		} else if firstError == nil {
			firstError = err
		}
	}
	return "", firstError
}

// EnsureFirst returns the path to the given resource in the user-specific XDG
// directory; if it doesn't exist it is created before returning.
//
// A resource is usually an application name, followed by a filename. Nothing
// in the standard nor in this library limits you to that, although
// EnsureFirst supports that usage better than the alternatives.
//
// If only a single component of the resource path is given, it is assumed to
// be a directory. If many are given, it is assumed that the last one is a
// filename. To create a file in the base directory itself specify a directory
// of "." before the filename; to create a folder tree without a file at the
// end specify a last element of "".
func (x *XDGDir) EnsureFirst(resourcePath ...string) (fullPath string, err error) {
	filename := ""
	if len(resourcePath) > 1 {
		l := len(resourcePath) - 1
		resourcePath, filename = resourcePath[:l], resourcePath[l]
	}
	resource := filepath.Join(x.Home(), filepath.Join(resourcePath...))
	err = os.MkdirAll(resource, 0700)
	if err != nil {
		return "", err
	}
	if filename == "" {
		return resource, nil
	} else {
		filename = filepath.Join(resource, filename)
		f, err := os.OpenFile(filename, os.O_CREATE, 0600)
		if err != nil {
			return "", err
		}
		f.Close()
		return filename, nil
	}
}
