# Dotty
Converts images to braille text

## Features
* Use local files or HTTP(S) URLs
* Fine-tune the brightness threshold value to get the best result with `-c`
* Invert the image with `-i`
* Allows choice of `-width` and `-height` of output text

## Installation
`go get github.com/bspammer/dotty && go install github.com/bspammer/dotty`

## Examples

### Baboon

![](https://homepages.cae.wisc.edu/~ece533/images/baboon.png)
![](https://i.imgur.com/uBwvAHf.png)

`dotty -width 60 https://homepages.cae.wisc.edu/~ece533/images/baboon.png`

### Lenna

![](https://upload.wikimedia.org/wikipedia/en/7/7d/Lenna_%28test_image%29.png)
![](https://i.imgur.com/WHTOCiz.png)

`dotty -width 60 https://upload.wikimedia.org/wikipedia/en/7/7d/Lenna_%28test_image%29.png`
