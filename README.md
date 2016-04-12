# Docker Machine - Atlantic.Net Driver

[Atlantic.Net](https://atlantic.net) driver for [Docker Machine](https://www.docker.com/products/docker-machine)

## Status

Alpha

## Getting source

* Install the latest version of [Go](https://golang.org)
* Make sure your GOPATH is set
* Add $GOPATH/bin to your PATH
* Use `go get` to download the source
```
$ go get -d github.com/jdextraze/docker-machine-driver-atlanticnet
```

## Dependencies

* Install [govendor](https://github.com/kardianos/govendor)
```
go get -u github.com/kardianos/govendor
```
* Go into source folder
```
cd $GOPATH/src/github.com/jdextraze/docker-machine-driver-atlanticnet
```
* Sync vendor folder
```
$ govendor sync
```

## Installation from source

* Get source (See above)
* Sync dependencies
* Build && install the driver
```
$ go install github.com/jdextraze/docker-machine-driver-atlanticnet
```

## License

Released under the MIT license, see [LICENSE](https://github.com/jdextraze/go-atlanticnet/blob/master/LICENSE).