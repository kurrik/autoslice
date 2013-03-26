// Copyright 2013 Arne Roomann-Kurrik
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"os"
	"path"
)

const (
	THRESH   = 20
	THRESH_A = 250
)

type Region struct {
	image.Rectangle
}

type AutoSlicer struct {
	SrcPath string
	DstPath string
	img     image.Image
	thresh  color.Color
}

func (hs *AutoSlicer) isEdge(pt color.Color) bool {
	var (
		r, g, b, a     = pt.RGBA()
		tr, tg, tb, ta = hs.thresh.RGBA()
	)
	return r <= tr && g <= tg && b <= tb && a >= ta
}

func (hs *AutoSlicer) isInRegion(x int, y int, regions []Region) bool {
	var pt = image.Pt(x, y)
	for _, r := range regions {
		if pt.In(r.Rectangle) {
			return true
		}
	}
	return false
}

func (hs *AutoSlicer) findRun(img image.Image, x int, y int, dx int, dy int) int {
	var count = -1
	for hs.isEdge(img.At(x, y)) {
		x += dx
		y += dy
		count += 1
	}
	return count
}

func (hs *AutoSlicer) mergeRegions(regions []Region, region *Region) []Region {
	if region.Empty() {
		// New region is empty
		return regions
	}
	if len(regions) == 0 {
		// No existing regions
		regions = append(regions, *region)
		return regions
	}
	for i := len(regions) - 1; i >= 0; i-- {
		r := regions[i]
		if region.In(r.Rectangle) {
			// Subset of existing regions
			return regions
		}
		if region.Overlaps(r.Rectangle) {
			// Neither is valid
			if i < len(regions)-1 {
				regions = append(regions[:i], regions[i+1:]...)
			} else {
				regions = regions[:i]
			}
			return regions
		}
	}
	regions = append(regions, *region)
	return regions
}

func (hs *AutoSlicer) scanImage() (regions []Region, err error) {
	var (
		bounds       image.Rectangle
		x, y         int
		w, h, w2, h2 int
		r            *Region
	)
	bounds = hs.img.Bounds()
	for x = bounds.Min.X; x < bounds.Max.X; x++ {
		for y = bounds.Min.Y; y < bounds.Max.Y; y++ {
			if hs.isEdge(hs.img.At(x, y)) {
				w = hs.findRun(hs.img, x, y, 1, 0)
				h = hs.findRun(hs.img, x, y, 0, 1)
				w2 = hs.findRun(hs.img, x, y+h, 1, 0)
				h2 = hs.findRun(hs.img, x+w, y, 0, 1)
				if w != w2 || h != h2 {
					continue
				}
				r = &Region{image.Rect(x, y, x+w, y+h)}
				regions = hs.mergeRegions(regions, r)
			}
		}
	}
	return
}

func (hs *AutoSlicer) readImage() (err error) {
	var f *os.File
	if f, err = os.Open(hs.SrcPath); err != nil {
		return
	}
	defer f.Close()
	if hs.img, _, err = image.Decode(f); err != nil {
		return
	}
	hs.thresh = hs.img.ColorModel().Convert(color.RGBA{THRESH, THRESH, THRESH, THRESH_A})
	return
}

func (hs *AutoSlicer) writeImage(img image.Image, filename string) (err error) {
	var f *os.File
	if f, err = os.Create(path.Join(hs.DstPath, filename)); err != nil {
		return
	}
	defer f.Close()
	err = png.Encode(f, img)
	return
}

func (hs *AutoSlicer) getDrawable() draw.Image {
	var (
		bounds image.Rectangle
		img    draw.Image
	)
	bounds = hs.img.Bounds()
	img = image.NewRGBA(bounds)
	draw.Draw(img, bounds, hs.img, bounds.Min, draw.Src)
	return img
}

func (hs *AutoSlicer) getSlice(r Region) image.Image {
	var (
		bounds image.Rectangle
		img    draw.Image
	)
	bounds = r.Sub(r.Min)
	bounds.Max = bounds.Max.Sub(image.Pt(1, 1))
	img = image.NewRGBA(bounds)
	draw.Draw(img, bounds, hs.img, r.Min.Add(image.Pt(1, 1)), draw.Src)
	return img
}

func (hs *AutoSlicer) drawRegion(img draw.Image, r Region) {
	var (
		c    color.Color
		x, y int
	)
	c = color.RGBA{0xFF, 0, 0, 0xFF}
	for x = r.Min.X + 1; x < r.Max.X; x++ {
		for y = r.Min.Y + 1; y < r.Max.Y; y++ {
			img.Set(x, y, c)
		}
	}
}

func (hs *AutoSlicer) Slice() (err error) {
	var (
		regions  []Region
		img      draw.Image
		slice    image.Image
		filename string
	)
	if err = hs.readImage(); err != nil {
		return
	}
	if regions, err = hs.scanImage(); err != nil {
		return
	}
	fmt.Println("Exporting regions:")
	img = hs.getDrawable()
	for i, r := range regions {
		hs.drawRegion(img, r)
		slice = hs.getSlice(r)
		filename = fmt.Sprintf("image%04v.png", i+1)
		if err = hs.writeImage(slice, filename); err != nil {
			return
		}
		fmt.Printf("  * %v - %v\n", filename, r)
	}
	fmt.Println("Writing out.png.")
	err = hs.writeImage(img, "out.png")
	return
}
