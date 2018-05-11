# Bit

A simple and safe way to learn git. The primary goal of this project is make a tool to help people learn how to use git. The secondary goal is to make a fun and opinionated high level tool.

# Development

Set up Golang and Glide (https://github.com/Masterminds/glide).

```shell
go get github.com/ArcherWheeler/bit
cd $GOPATH/src/github.com/ArcherWheeler/bit/bit
glide install
export PATH=$PATH:$GOPATH/bin
go install
```

Bit should now work from the command line.

## Tests

In order to run the tests locally you need to set up docker (https://docs.docker.com/install/).

Go's built in testing tools are designed to compile and run go code in a very specific manner. However, our setup compiles and installs bit as a full command line tool and then tests that. Because of this, we want a space to set up git from scratch and build temporary repos without messing real things up. We could do this carefully in a script, but docker is a nice solution for this.

Once you have docker running you can build and run the tests with `./run-tests` from the root of the repo.

# More coming soon?
