# pxgo
pxgo is a command line client for ftsell's [pixelflut](https://github.com/ftsell/pixelflut) written in Go.
It supports three different modes:
* Fill the canvas with just one color (specified by its hexcode)
```
pxgo --color 0328ff
```
* Draw a PNG or JPG
```
pxgo --file image.jpg
```
* Draw a rainbow
```
pxgo --rainbow
```

## Installation
Assuming Go is already installed on your machine, run the following commands:
```
go get github.com/fredericobormann/pxgo
go install github.com/fredericobormann/pxgo
```
