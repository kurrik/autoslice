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


