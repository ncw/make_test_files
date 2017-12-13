Fsync bench
===========

Benchmark fsync

Download a [binary from github](https://github.com/ncw/make_test_files/releases/latest)
or build from source (see later).

Usage
=====

Usage: make_test_files [flags] <directory>

This command makes a random directory structure with random files in
<directory>.  The options can be used to control exactly which files
get made.

The file names and sizes will be identical each time the command is
run with the same parameters.  -seed can be used to change what is
created.

Options:
  -files-per-directory int
    	Average number of files per directory (default 10)
  -loop
    	Loop forever
  -max-depth int
    	Maximum depth of directory heirachy (default 10)
  -max-name-length int
    	Maximum size of files to create (default 12)
  -max-size int
    	Maximum size of files to create (default 100)
  -min-name-length int
    	Minimum size of file to create (default 4)
  -min-size int
    	Minimum size of file to create
  -n int
    	Number of files to create (default 1000)
  -seed int
    	Seed for the random number generator (default 1)
  -sync
    	Fsync each file
  -v	Be more verbose
  -z	Fill files with zeroes instead of random data

Build
=====

You'll need go installed, then 

    go get github.com/ncw/make_test_files

and this will build the binary in `$GOPATH/bin`.  You can then modify
the source and submit patches.

License
=======

This is free software under the terms of the MIT license (check the
COPYING file included in this package).

Contact and support
===================

The project website is at:

- https://github.com/ncw/make_test_files

There you can file bug reports, ask for help or contribute patches.

Authors
=======

- Nick Craig-Wood <nick@craig-wood.com>

Contributors
------------

- Your name goes here!
