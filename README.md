# Copy Library

[![Yum](https://img.shields.io/badge/-Buy%20me%20a%20cookie-blue?labelColor=grey&logo=cookiecutter&style=for-the-badge)](https://www.buymeacoffee.com/mjwhitta)

[![Go Report Card](https://goreportcard.com/badge/github.com/mjwhitta/cplib?style=for-the-badge)](https://goreportcard.com/report/github.com/mjwhitta/cplib)
![License](https://img.shields.io/github/license/mjwhitta/cplib?style=for-the-badge)

## What is this?

This tool will generate Go source for you. It will read the exports of
a shared library (DLL or so) or it will read the imports of an
executable. The imports can be filtered down by the library they are
expected to be found in.

## How to install

Open a terminal and run the following:

```
$ go install github.com/mjwhitta/cplib/cmd/cplib@latest
```

Or compile from source:

```
$ git clone https://github.com/mjwhitta/cplib.git
$ cd cplib
$ git submodule update --init
$ make
```

## How to use

**NOTE:** The generated source is meant to be used with something like
[goDLL] so you may need to adjust the build tags.

List exports from a DLL:

```
$ cplib -e "c:/program files/windows defender/mpclient.dll"
```

Generate source for all exports from a DLL:

```
$ cplib -e -g "c:/program files/windows defender/mpclient.dll"
```

List exports from a shared-object:

```
$ cplib -e /usr/lib/x86_64-linux-gnu/libyaml.so
```

Generate source for all exports from a DLL:

```
$ cplib -e -g /usr/lib/x86_64-linux-gnu/libyaml.so
```

List specific imports of an executable:

```
$ cplib -f mpclient.dll -i "c:/program files/windows defender/mpcmdrun.exe"
```

Generate source for specific imports of an executable:

```
$ cplib -f mpclient.dll -g -i "c:/program files/windows defender/mpcmdrun.exe"
```

## Links

- [Source](https://github.com/mjwhitta/cplib)

## TODO

- Add support for macOS

[goDLL]: https://github.com/mjwhitta/godll
