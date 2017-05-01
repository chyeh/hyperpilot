# Introduction

In this implementation of the experiment, there are two stages:

1. Find an answer by binary search.
2. Improve the answer by trying to reduce both CPU usage and memory usage.

# Run

There are two ways to run the experiment of my implemenation:

1. Use `go test`:

```shell
$ go test -check.f IntegrationTestMainSuite
```

1. Build and run the executable:

```shell
$ go build -o hyperpilot
$ ./hyperpilot
```
