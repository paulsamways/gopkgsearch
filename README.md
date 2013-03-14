Web application for searching the packages in $GOROOT/src/pkg and $GOPATH.

## Prerequisites

Gopkgsearch requires that the Go source code is available on the machine and that the GOROOT environment variable has been set. 

See http://golang.org/doc/install/source#environment

## Installation

Gopkgsearch can be fetched using 'go get':

  `go get github.com/PaulSamways/gopkgsearch`

## Usage

### Starting the server

When GoPkgSearch is launched without any parameters, the Go stdlib source files will indexed and the web application started.

Indexing of the source files in the GOPATH directories can be enabled by using the '-useGoPath' option.

  `./gopkgsearch [-usegopath=true]`

### Using the client

After the GoPkgSearch server has finished indexing the source files and the web server has been started, open up a browser and navigate to the web application, by default http://localhost:8000.

## Screenshot

![Screenshot of Gopkgsearch](http://paulsamways.github.com/gopkgsearch/images/gopkgsearch.gif)
