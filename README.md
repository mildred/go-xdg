Go, XDG, go!
===========

This is `go-xdg`, a little library to help you use the `XDG`
[base directory spec](http://standards.freedesktop.org/basedir-spec/basedir-spec-latest.html).

(There are other `XDG` specs, that might get included in time. Patches welcome.)

Sample usage
------------

Let's say you are writing an app called “frobz”. It has a config file
and a sqlite database. You'd do something like this:

    configFileName, err := xdg.Config.FirstExisting("frobz/config.txt") // XXX see below
    if err == nil {
        // a config file exists! load it...
    }
    dbFileName, err := xdg.Data.EnsureFirst("frobz", "frobz.db")
    // etc


Resources
---------

Both `FirstExisting` and `EnsureFirst` take a `resource` as a list of
components to construct the path they return.

A resource is usually an application name (or a well-known shared resource
pool name, such as `icons`), followed by a filename. However nothing in the
standard nor in this library limits you to that; you may store e.g. your
application's configuration in just `$XDG_CONFIG_HOME/application.conf` (in
which case the "resource" here would be just `application.conf`), or in a
sub-directory of an application-sepcific directory.

While this library (and the spec) isn't currently portable, some work
has been done to make it so, and in case that prospers (and you care
about portability (and you probably should)) it is recommended that
you specify resources as a list, as in the `EnsureFirst` sample call,
and not as a single path component as in the `FirstExisting` sample
call above. Nothing in the library (or the spec...) will stop you if
you don't want to do it that way, however. Note also that
`EnsureFirst` supports the recommended usage a lot better than the
alternatives.

License, etc.
------------

BSD simplified, © John R. Lenton, blah blah.
