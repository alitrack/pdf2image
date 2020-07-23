package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/alitrack/imaging"
	imaging1 "github.com/disintegration/imaging"
	"github.com/gen2brain/go-fitz"
)

const (
	dpi = 72
)

func main() {

	if l := len(os.Args); l != 3 {
		help()
		return
	}

	path := os.Args[1]

	images, err := pdf2Images(path)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	countX := len(images)
	rgba, err := imaging.MergeGrids(images, 1, countX)

	// file, err := os.Create(os.Args[2])
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// err = png.Encode(file, rgba)
	err = imaging1.Save(rgba, os.Args[2])
	if err != nil {
		fmt.Println(err.Error())
	}

}

func help() {
	fmt.Println("Usage: pdf2longimage <input pdf path> <output image path>")
	fmt.Println("convert a multi page PDF document to a long image")
	fmt.Println("Image format support: jpg(or jpeg), png, gif, tif (or tiff) and bmp.")
	fmt.Println("Author: Steven Lee")
	fmt.Println("Website: alitrack.com")
}

func pdf2Images(path string) ([]image.Image, error) {
	doc, err := fitz.New(path)
	if err != nil {
		panic(err)
	}

	var images []image.Image

	for n := 0; n < doc.NumPage(); n++ {
		img, err := doc.ImageDPI(n, dpi)
		if err != nil {
			panic(err)
		}
		images = append(images, img)
	}
	return images, nil
}

func pdf2Images1(path string) ([]image.Image, error) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	doc, err := fitz.New(path)
	if err != nil {
		panic(err)
	}

	var images []image.Image

	mImg := make(map[int]image.Image)

	var wg sync.WaitGroup

	pn := doc.NumPage()
	wg.Add(pn)

	//出现报错，怎么处理？
	for n := 0; n < pn; n++ {

		go func(i int, wg *sync.WaitGroup) {
			defer wg.Done()
			img, err := doc.ImageDPI(i, dpi)
			if err != nil {
				log.Fatalln(err)
			}
			mImg[i] = img
		}(n, &wg)
	}
	wg.Wait()

	for i := 0; i < pn; i++ {
		images = append(images, mImg[i])
	}

	return images, nil
}
