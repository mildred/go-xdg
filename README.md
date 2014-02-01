Go, XDG, go!
===========

This is `go-xdg`, a little library to help you use the `XDG`
[base directory spec](http://standards.freedesktop.org/basedir-spec/basedir-spec-latest.html).

(There are other `XDG` specs, that might get included in time. Patches welcome.)

Sample usage
------------

Let's say you are writing an app called “frobz”. It has a config file
and a sqlite database. You'd do something like this:

    configFileName, err := xdg.Config.Find("frobz/config.txt")
    if err == nil {
        // a config file exists! load it...
    }
    dbFileName, err := xdg.Data.Ensure("frobz/frobz.db")
    // now the file and all its directories exist; it's up to you to
    // determine if it's empty, etc.


Resources
---------

Both `Find` and `Ensure` take a `resource` to construct the path they return.

A resource is usually an application name (or a well-known shared resource
pool name, such as `icons`), followed by a filename. However nothing in the
standard nor in this library limits you to that; you may store e.g. your
application's configuration in just `$XDG_CONFIG_HOME/application.conf` (in
which case the "resource" here would be just `application.conf`), or in a
sub-directory of an application-specific directory.

License, etc.
------------

BSD simplified, © John R. Lenton, blah blah.
