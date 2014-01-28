// (c) 2014 John R. Lenton. See LICENSE.

// xdg implements helpers for you to use the XDG base directory spec in your
// apps.
package xdg

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type XDGDir interface {
	// Home() gets the path to the XDG directory in the user's home, as
	// specified (or not) by the user's environment.
	Home() string
	// Dirs() gets the list of paths to the given XDG directories in the
	// user's home, as specified (or not) by the user's environment.
	Dirs() []string
	// FirstExisting returns the path to the first of the given resource
	// that exists in the XDG directories. If none exist, error is set
	// appropriately.
	FirstExisting(...string) (string, error)
	// EnsureFirst returns the path to the given resource in the user's
	// xdg directory; if it doesn't exist it is created before returning.
	EnsureFirst(...string) (string, error)
}

type xdgd struct {
	homeEnv  string
	homeDflt string
	dirsEnv  string
	dirsDflt string
}

var (
	Data XDGDir = &xdgd{"XDG_DATA_HOME", ".local/share",
		"XDG_DATA_DIRS", "/usr/local/share:/usr/share"}
	Config XDGDir = &xdgd{"XDG_CONFIG_HOME", ".config",
		"XDG_CONFIG_DIRS", "/etc/xdg"}
	Cache XDGDir = &xdgd{"XDG_CACHE_HOME", ".cache", "", ""}
)

func (x *xdgd) Home() string {
	dir := os.Getenv(x.homeEnv)
	if dir != "" {
		return dir
	}
	home := os.Getenv("HOME")
	if home == "" {
		user, err := user.Current()
		if err != nil {
			// this bit will stay untested, I suspect
			panic("Unable to determine $HOME")
		}
		home = user.HomeDir
	}
	return filepath.Join(home, x.homeDflt)
}

func (x *xdgd) Dirs() []string {
	dirs := []string{x.Home()}
	if x.dirsEnv != "" {
		xtra := os.Getenv(x.dirsEnv)
		if xtra == "" {
			xtra = x.dirsDflt
		}
		for _, path := range strings.Split(xtra, ":") {
			if path != "" {
				dirs = append(dirs, path)
			}
		}
	}
	return dirs
}

func (x *xdgd) FirstExisting(resources ...string) (string, error) {
	for _, path := range x.Dirs() {
		name := filepath.Join(path, filepath.Join(resources...))
		_, err := os.Stat(name)
		if err == nil {
			return name, nil
		}
	}
	return "", os.ErrNotExist
}

func (x *xdgd) EnsureFirst(resources ...string) (string, error) {
	filename := ""
	if len(resources) > 1 {
		l := len(resources) - 1
		resources, filename = resources[:l], resources[l]
	}
	resource := filepath.Join(x.Home(), filepath.Join(resources...))
	err := os.MkdirAll(resource, 0700)
	if err != nil {
		return "", err
	}
	if filename == "" {
		return resource, nil
	} else {
		filename = filepath.Join(resource, filename)
		_, err = os.OpenFile(filename, os.O_CREATE, 0600)
		if err != nil {
			return "", err
		}
		return filename, nil
	}
}
