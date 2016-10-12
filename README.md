# 📃 goprove

Inspect your project for the best practices listed in the [Go CheckList](https://github.com/matttproud/gochecklist).

### Get Started

    $ go get github.com/karolgorecki/goprove/cmd/goprove
    $ goprove .

### Usage

```
Usage:
  SIMPLE:
  	goprove <directory>
  WITH OUTPUT:
  	goprove -output=<output: json or text> <directory>
  WITH EXCLUDE:
  	goprove -exclude=<tasks: separated by comma> <directory>
Available tasks for exclude:
  projectBuilds, isFormatted, hasLicense, isLinted, isVetted, hasReadme,
  testPassing, isDirMatch, hasContributing, hasBenches, hasBlackboxTests
```
### Contributing
Contributions are most welcome.
Instructions are documented in [CONTRIBUTING.md](CONTRIBUTING.md).

### License
MIT
