# mrlog

mrlog (Machine Readable Log) is a utility used to create log messages that will be parsed by its sister program, [mrreport](https://github.com/cf-platform-eng/mrreport).

## Commands

### Sections

Sections put bookends in the logs around a notable section.

The start and end of the section can be logged separately:

```bash
mrlog section-start --name="run-test"
test_runner execute
mrlog section-end --name="run-test" --result $? 
```

or combined:

```bash
mrlog section --name="run-test" \
      --on-failure="The test failed" \
      --on-success="The test passed" \
      -- test_runner execute
```

#### Examples

```bash
$ mrlog section --name="show-date" --on-success="successfully got the date" --on-failure="failed to get the date" -- date
section-start: 'show-date' MRL:{"type":"section-start","name":"show-date","time":"2021-02-22T13:21:40.132922-06:00"}
Mon Feb 22 13:21:40 CST 2021
section-end: 'show-date' result: 0 message: 'successfully got the date' MRL:{"type":"section-end","name":"show-date","time":"2021-02-22T13:21:40.137741-06:00","message":"successfully got the date"}
```

### Dependency

mrlog has a built-in way of logging dependencies, useful in recording exact versions of other tools involved.

#### Examples

```bash
$ mrlog dependency --type binary --name kubectl --version $(kubectl version  --client -o json | jq -r .clientVersion.gitVersion)
binary dependency: 'kubectl' version 'v1.20.2' MRL:{"type":"binary dependency","version":"v1.20.2","name":"kubectl","metadata":"","time":"2021-02-22T13:30:34.213109-06:00"}
```

## Developing

Utilize the Makefile for testing and building.

`make test` will execute the unit tests

`make build` will build the `mrlog` binary
