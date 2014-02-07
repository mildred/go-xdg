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
	homeDefault string
	dirsEnv  string
	dirsDefault string
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
	return filepath.Join(home, x.homeDefault)
}

// Dirs returns the preference-ordered set of base directories to search for
// files of the given type, starting with the user-specific one, as specified
// (or not) by the user's environment.
func (x *XDGDir) Dirs() []string {
	dirs := []string{x.Home()}
	if x.dirsEnv != "" {
		xtra := os.Getenv(x.dirsEnv)
		if xtra == "" {
			xtra = x.dirsDefault
		}
		for _, path := range filepath.SplitList(xtra) {
			if path != "" {
				dirs = append(dirs, path)
			}
		}
	}
	return dirs
}

// Find attempts to find the path suffix in all of the known XDG directories.
// If not found, an error is returned.
func (x *XDGDir) Find(suffix string) (absPath string, err error) {
	var firstError error = nil
	for _, path := range x.Dirs() {
		name := filepath.Join(path, suffix)
		_, err = os.Stat(name)
		if err == nil {
			return name, nil
		} else if firstError == nil {
			firstError = err
		}
	}
	return "", firstError
}

// Ensure takes the path suffix given, and ensures that a matching file exists
// in the home XDG directory. If it doesn't exist it is created. If it can't
// be created, or exists but is unreadable, an error is returned.
func (x *XDGDir) Ensure(suffix string) (absPath string, err error) {
	absPath = filepath.Join(x.Home(), suffix)
	err = os.MkdirAll(filepath.Dir(absPath), 0700)
	if err == nil {
		f, err := os.OpenFile(absPath, os.O_CREATE, 0600)
		if err == nil {
			f.Close()
		}
	}
	return
}
