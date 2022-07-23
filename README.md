# color-filter
This a go package to apply a color filter to a image made in go

### Description:
The first implementation will acept only ppm images, the objective of this proyect is to apply concurrency patters to filter an image by an rgb color filter.

To start it will consise of 3 go routines, one for every RGB value (Red, Green, Blue)

The input image will be readed in n bytes block, specified by the user. Every go routine will read the image that number of bytes at a time.

The output will be 3 distinct images (red, blue, green). 

Additionally there are 3 additional options *-r value* , *-g value* and *-b value* that scale the intensity of the color for the filter that is generated.

Once the jobs is finished it will print a massage show it completed.
```
$ ./color-filter -h
usage: color-filter filter [-h] [-r RED] [-g GREEN] [-b BLUE] -s SIZE -f FILE

Color Filter - Apply color filter to ppm image
  -h, --help                show this help message and exit
  -r RED, --red RED         Red scale
  -g GREEN, --green GREEN   Green scale
  -b BLUE, --blue BLUE      Blue scale
  -s SIZE, --size SIZE      Reading size
  -f FILE, --file FILE      File to process

Usage Example:
$ ./color-filter -s 1024 -f dog.ppm -r 2 -g 0 -b 0.5

Succesfuly filter

$ ls *ppm
dog.ppm
b_dog.ppm
g_dog.ppm
r_dog.ppm

```

### Links:
 - PPM Format implementation: http://netpbm.sourceforge.net/doc/ppm.html