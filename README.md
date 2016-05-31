# 📃 goprove

Inspect your project for the best practices listed in the [Go CheckList](https://github.com/matttproud/gochecklist).

[![Stories in Ready](https://badge.waffle.io/karolgorecki/goprove.png?label=ready&title=Ready)](https://waffle.io/karolgorecki/goprove)
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

### License
MIT
