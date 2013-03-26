Autoslice
=========

Slice images based off of features in the image.

![test source image](https://github.com/kurrik/autoslice/raw/master/test/test01.fw.png?raw=true)

Installing
----------

    go get -u github.com/kurrik/autoslice

Running
-------

    $ autoslice -dst=/path/to/dst/folder /path/to/source/image.png

    Exporting regions:
      * image0001.png - (139,139)-(340,340)
      * image0002.png - (379,79)-(530,130)
    Writing out.png.

The produced images look like:

### image0001.png
![test export 01](https://github.com/kurrik/autoslice/raw/master/dst/image0001.png?raw=true)

### image0002.png
![test export 02](https://github.com/kurrik/autoslice/raw/master/dst/image0002.png?raw=true)

### out.png
![test export out](https://github.com/kurrik/autoslice/raw/master/dst/out.png?raw=true)

