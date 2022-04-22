A Simple RESTful Config Server, supporting the usage elements of Spring Boot and Node "config" in Go
=====================================


![Language Badge](https://img.shields.io/github/languages/top/alt-golang/config) <br/>
[release notes](https://github.com/alt-golang/config/blob/main/HISTORY.md)

<a name="intro">Introduction</a>
--------------------------------
An simple config publishing server, supporting the usage elements of Spring Boot and Node "config", including:
- json, yaml and Java property files,
- cascading value over-rides using, GO_ENV, GO_APP_INSTANCE and GO_PROFILES_ACTIVE
- placeholder resolution (or variable expansion),



<a name="usage">Usage</a>
-------------------------

To use the module, go get the module as so:

```shell
go get github.com/alt-golang/config-server
go run .
```

### Internal Config Location
The Config Server is configured internally from the config/internal directory relative to the working directory

### Server Config Location
The Config Server defaults to publishing config from the config directory relative to the working directory, but can 
be configured with using the internal setting as below:

```yaml
config: 
  dir: "otherLocation" 
```

### File Publishing and Precedence

The module follows the file loading and precedence rules of the popular Node
[config](https://www.npmjs.com/package/config) defaults, with additional rules in the style of Spring Boot.

Files are loaded and over-ridden from the `config` folder in the following order:
- default.( json | yml | yaml | props | properties )
- application.( json | yml | yaml | props | properties )
- {GO_ENV}.( json | yml | yaml | props | properties )
- {GO_ENV}-{GO_APP_INSTANCE}.( json | yml | yaml | props | properties )
- {GO_ENV}-{GO_APP_INSTANCE}.( json | yml | yaml | props | properties )
- {GO_ENV}-{GO_APP_INSTANCE}.( json | yml | yaml | props | properties )
- application-{GO_PROFILES_ACTIVE[0]}.( json | yml | yaml | props | properties )
- application-{GO_PROFILES_ACTIVE[1]}.( json | yml | yaml | props | properties )


Environment variables and command line arguments, are not published by the service, as it is intended only as a networked
config discovery service.

### API 

The API URI format follows the form : `/:GO_ENV/:GO_APP_INSTANCE/*GO_PROFILES_ACTIVE?path=my.key` where trailing `/` 
act as additional comma-separated profiles.

if the path query param is ommitted, the full config set is returned by default.

<a name="license">License</a>
-----------------------------

May be freely distributed under the [MIT license](https://raw.githubusercontent.com/alt-golang/config/main/LICENSE).

Copyright (c) 2022 Craig Parravicini    
