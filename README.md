# asdf
Generates a changelog based on semantic commit messages.

## Usage

```
$ asdf --help
NAME:
   asdf - Changelog generation based on semantic commit messages

USAGE:
   asdf [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
     generate, g   generates a changelog and the next version based on semantic commits and writes them to file
     changelog, c  generates only the changelog and writes it to stdout
     help, h       Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --branch value  name of the current branch (default: "master") [$RELEASE_BRANCH]
   --token value   github token [$RELEASE_GITHUB_TOKEN]
   --help, -h      show help
   --version, -v   print the version
```
### Commit Message Schema
Commit messages have to follow the angularjs commit message conventions [[link](https://docs.google.com/document/d/1QrDFcIiPjSLDn3EL15IJygNPiHORgU1_OOAqWjiDU5Y/edit)].

#### Examples
- `test(PROJ-1312): write tests. do yourself a favor`
- `breaking(PROJ-1000): break all the things!`
- `(TICKK-123): foobar booman`
- `bug: Y U NO GOAT?`

#### Generated Changelog
Example Changelog file:
```
## 1.23.5 (2017-11-12)

#### Feature

* some commit message (53c7d6c2)

#### Bug Fixes

* fix typo (d9a3a253)
```

### default types
| Key | Label | Change type |
| --- | --- | --- |
| breaking | Breaking Changes | Major |
| feat | Feature | Minor |
| fix | Bug Fixes | Patch |
| perf | Performance Improvements | Patch |
| revert | Reverts | Patch |
| docs | Documentation | Patch |
| refactor | Code Refactoring | Patch |
| test | Tests | Patch |
| chore | Chores | Patch |


### TODO

#### release command?
- do the merges based on current branch
  releases should only happen on
    - release-type branches (see asdf.json -> branch_suffix)
      - -> should be merged into master
      - -> then merge master back to develop
  - master
    - merge back to develop
 