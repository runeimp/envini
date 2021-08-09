CHANGELOG
=========

2021-08-09 v0.1.0
-----------------

### Initial release

* Can parse a basic INI config file with one level of sections
* Can recursively populate a two level struct with data from the config file using struct tags
* Can use environment variables to update the struct if they are present and specified in struct tags
* Environment variables take precedence over config values

