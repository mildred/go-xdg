Go, XDG, go!
===========

This is `go-xdg`, a little library to help you use the `XDG`
[base directory spec](http://standards.freedesktop.org/basedir-spec/basedir-spec-latest.html).

(There are other `XDG` specs, that might get included in time. Patches welcome.)

Sample usage
------------

Let's say you are writing an app called “frobz”. It has a config file
and a sqlite database. You'd do something like this:

    configFile, err := xdg.Config.FirstExisting("frobz/config.txt")
    if err == nil {
        // a config file exists! load it...
    }
    dbFile, err := xdg.Data.EnsureFirst("frobz", "frobz.db")
    // etc


License, etc.
------------

GPLv3, © John R. Lenton, blah blah.
