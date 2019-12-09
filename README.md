[![Build Status](https://github.com/axetroy/denox/workflows/ci/badge.svg)](https://github.com/axetroy/denox/actions)
[![Coverage Status](https://coveralls.io/repos/github/axetroy/denox/badge.svg?branch=master)](https://coveralls.io/github/axetroy/denox?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/axetroy/denox)](https://goreportcard.com/report/github.com/axetroy/denox)
![Latest Version](https://img.shields.io/github/v/release/axetroy/denox.svg)
![License](https://img.shields.io/github/license/axetroy/denox.svg)
![Repo Size](https://img.shields.io/github/repo-size/axetroy/denox.svg)

### Execute Deno script even if you don't have Deno installed

> Why? It looks the same as Deno's command line, so why do I need such a tool?
> There are scenarios where I need to run the same script with different versions of Deno
> In such scenarios, Deno's version manager may not be the best option

### Features

- [x] Cross platform support
- [x] Install Deno automatically
- [x] Support any version of Deno with environment variable `DENO_VERSION`
- [x] Consistent with Deno

### Usage

```bash
# run script with latest version of Deno
$ denox https://deno.land/std/examples/welcome.ts
# run script with specific version of Deno
$ DENO_VERSION=v0.26.0 denox https://deno.land/std/examples/welcome.ts
```

### Installation


Download the executable file for your platform at [release page](https://github.com/axetroy/denox/releases)

Then set the environment variable.

eg, the executable file is in the `~/bin` directory.

```bash
# ~/.bash_profile
export PATH="$PATH:~/bin"
```

finally, try it out.

```bash
$ denox https://deno.land/x/std/examples/welcome.ts
```

### Build from source code

```bash
$ git clone github.com/axetroy/denox $GOPATH/src/github.com/axetroy/denox
$ cd $GOPATH/src/github.com/axetroy/denox
$ make build
```

### Test

```bash
$ make test
```

### License

The [MIT License](LICENSE)