EnvINI v0.1.0
=============

A small and simple Go configuration library. It is also a command line tool that reviews an INI file and outputs it's contents as JSON to better understand what your config struct should expect

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
* [ ] Support alternate data layouts
	* [ ] Support alternate key value delimiters such as `:` instead of just `=`
	* [ ] Support alternate record delimiters such as `\t` or `|` instead of just `\n` or `\r\n`
	* [ ] SUpport alternate section delimiters
* [ ] Support accessing comments in code
* [ ] Support `array`/`slice` value types
* [ ] Support `interface{}` value types
* [ ] Support `rune` value types
* [ ] Support nested sections in INI files
* [ ] Support setting of acceptable values for `bool` fields


Conversions
-----------

EnvINI will try to convert the values in a UTF-8/ASCII encoded INI text file or UTF-8/ASCII encoded `[]byte` to matches the fields of the target config struct. If there is an error from say trying to assign `256` to a `uint8` or `-1` to a `uint` the error will be returned. For `bool` fields in a target struct EnvINI will convert the following INI values as shown:

| INI Value          | Go Value |
| :-------:          | :------: |
| true, t, yes, y, 1 | `true`   |
| false, f, no, n, 0 | `false`  |

This will be adjustable in the future.


Articles & Reference
--------------------

* [INI file - Wikipedia][]



[INI file - Wikipedia]: https://en.wikipedia.org/wiki/INI_file

