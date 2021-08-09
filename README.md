EnvINI
======

A small and simple Go configuration library.

* Reads INI file for config
* Checks Environment Variables for system overrides to the config
* Config can be used as defaults for command line parameters


Features
--------

* [x] Parse basic Windows API and POSIX style INI files
* [x] Recursively populate a two level struct with an INI file with normal sections
* [x] Support `bool` value types
* [x] Support `float32`, `float64` number types
* [x] Support `int`, `int8`, `int16`, `int32`, `int64` number types
* [x] Support `string` value types
* [x] Support `struct` value types (for recursion)
* [x] Support `uint`, `uint8`, `uint16`, `uint32`, `uint64` number types
* [x] Supports Windows API and POSIX style comments
* [x] UTF-8 based, which includes ASCII support automatically
* [ ] Create/Update INI Files
	* [ ] Non-destructive Updates
		* [ ] Keeps comments
		* [ ] Keeps line order
* [ ] Parse multi-line text values
* [ ] Support alternate data formats that are INI File like
	* [ ] Support alternate key value delimiters such as `:` instead of just `=`
	* [ ] Support alternate record delimiters such as `\t` or `|` instead of just `\n` or `\r\n`
	* [ ] SUpport alternate section delimiters
* [ ] Support accessing comments in code
* [ ] Support `array`/`slice` value types
* [ ] Support `interface{}` value types
* [ ] Support `rune` value types
* [ ] Support nested sections in INI files



Articles & Reference
--------------------

* [INI file - Wikipedia][]



[INI file - Wikipedia]: https://en.wikipedia.org/wiki/INI_file

