# mrlog

mrlog (Machine Readable Log) is a utility used to create log messages that will be parsed by its sister program, [mrreport](https://github.com/cf-platform-eng/mrreport).

Currently, it can log two specific things:

* Sections - indicating the start and stop of a chunk of work
* Dependencies - logging the versions of dependencies for a given execution

## Developing

Utilize the Makefile for testing and building.

`make test` will execute the unit tests

`make build` will build the `mrlog` binary
