# smcl
Generates a changelog based on semantic commit messages.

## Usage

```
$ smcl --help
NAME:
   smcl - Changelog and version generation based on semantic commit messages.

   Specification about the structure is here:
   https://github.com/figome/figo-rfc/blob/master/docs/COMMIT_MESSAGE.md

USAGE:
   smcl [global options] command [command options] [arguments...]

VERSION:
   0.4.0

COMMANDS:
     next-version, n  Tells you the next version you want to release. By default it uses a VERSION file to fetch the history since the last release. the file location may be overridden via --file
     generate, g      generates a changelog and the next version based on semantic commits and writes them to files
     changelog, c     generates the changelog and writes it to stdout. By default it uses a VERSION file to fetch the history since the last release. This can be overridden by defining a --version and --revision
     help, h          Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --dir value    set the current working directory
   --debug        show debug logs
   --help, -h     show help
   --version, -v  print the version
```
#### Example generated CHANGELOG.md
```
## 0.2.0 (2017-11-14)

#### 

* * WIP: generate tests (8835f264) 

#### Documentation

* display usage (d2bd0903) 

#### Feature

* next command (e6beb561) 
* parsing commit body (6c0eed43) 

#### Bug Fixes

* changelog test depended on time oo (ac4b9d06) 
* parsing multi-line commits (9c1b0bea) 

#### Code Refactoring

* removing foo.json (eb481d20) 
* changelog cmd (c80c64e5) 
* getting rid of type constraints (78361cc1) 
* cmd now in repo root (94e034d4) 

#### Tests

* config (ddf8535f) 
* cli commands (dcd068a0) 
```
 