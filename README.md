# pullreq

pullreq is just like [pullr](https://github.com/mojotech/pullr) but
written with [Go](http://golang.org)

# Installation
You need [go](http://golang.org), and make sure it's installed
and setup [correctly](https://golang.org/doc/install#testing).

Then do:

```shell
go get -u github.com/julien/pullreq
```

# Usage

```shell
Usage of pullreq:
  -d string
        pull request description
  -f string
        source branch (optional) (default = current branch)
  -from-remote string
        source remote (optional) (default "origin")
  -i string
        target branch (optional) (default "master")
  -into-remote string
        target remote (optional) (default "origin")
  -t string
        pull request title
```

# Examples

```shell
# Open a PR from the current branch to "origin:master"
pullreq -t "My feature" -d "Add easter-egg into website"
```

```shell
# Open a PR from the "origin:my_feature" to "some_remote:dev"
pullreq -t "My feature" -d "Add easter-egg into website" -f my_feature -i dev -into-remote some_remote
```

etc...
