package toold

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"

	// _ "image/png"
	// "image/jpeg"
	_ "image/jpeg"
	// "image/png"
	// _ "image/png"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
	"golang.org/x/image/bmp"
)

//ImageBounds ImageBounds
type ImageBounds struct {
	Width  int
	Height int
}

//GetImageExtFromFile GetImageExtFromFile
func GetImageExtFromFile(file *multipart.FileHeader) string {
	ext := filepath.Ext(file.Filename)
	return GetImageExt(ext)
}

/*
GetImageExtFromPath 获取图片ext
*/
func GetImageExtFromPath(path string) string {
	ext := filepath.Ext(path)
	return GetImageExt(ext)
}

//GetImageExt GetImageExt
func GetImageExt(ext string) string {
	extL := strings.ToLower(ext)
	switch extL {
	case ".png", ".jpg", ".bmp", ".jpeg":
		break
	case "png", "jpg", "bmp", "jpeg":
		ext = fmt.Sprintf(".%v", ext)
		break
	default:
		return ""
	}
	return ext
}

//GetImageExt GetImageExt
func GetGifExt(ext string) string {
	extL := strings.ToLower(ext)
	switch extL {
	case ".gif":
		break
	case "gif":
		ext = fmt.Sprintf(".%v", ext)
		break
	default:
		return ""
	}
	return ext
}

//GetImageInfo GetImageInfo
func GetImageInfo(file *multipart.FileHeader) (img image.Image, bounds *ImageBounds, ext string, err error) {
	f, err := file.Open()
	if err != nil {
		return
	}

	defer f.Close()
	img, ext, err = image.Decode(f)
	width := img.Bounds().Max.X - img.Bounds().Min.X
	height := img.Bounds().Max.Y - img.Bounds().Min.Y
	bounds = &ImageBounds{
		Width:  width,
		Height: height,
	}
	return
}

//GetImageInfoFromByte GetImageInfoFromByte
func GetImageInfoFromByte(imgBody []byte) (img image.Image, bounds *ImageBounds, ext string, err error) {
	img, ext, err = image.Decode(ConversionReaderFromByte(imgBody))
	width := img.Bounds().Max.X - img.Bounds().Min.X
	height := img.Bounds().Max.Y - img.Bounds().Min.Y
	bounds = &ImageBounds{
		Width:  width,
		Height: height,
	}
	return
}

//ImageRaw ImageRaw
func ImageRaw(imgBody []byte, x0, y0, x1, y1 int, quality int) error {
	origin, ext, err := image.Decode(ConversionReaderFromByte(imgBody))
	if err != nil {
		return err
	}
	canvas := origin
	var out bytes.Buffer
	switch ext {
	case "jpeg":
		img := origin.(*image.YCbCr)
		subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.YCbCr)
		return jpeg.Encode(&out, subImg, &jpeg.Options{Quality: quality})
	case "png":
		switch canvas.(type) {
		case *image.NRGBA:
			img := canvas.(*image.NRGBA)
			subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.NRGBA)
			return png.Encode(&out, subImg)
		case *image.RGBA:
			img := canvas.(*image.RGBA)
			subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.RGBA)
			return png.Encode(&out, subImg)
		}
	case "gif":
		img := origin.(*image.Paletted)
		subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.Paletted)
		return gif.Encode(&out, subImg, &gif.Options{})
	case "bmp":
		img := origin.(*image.RGBA)
		subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.RGBA)
		return bmp.Encode(&out, subImg)
	default:
		return errors.New("ERROR FORMAT")
	}
	return nil
}

//ImageBodyToJpeg ImageBodyToJpeg
func ImageBodyToJpeg(imgBody []byte) ([]byte, error) {
	img, _, err := image.Decode(ConversionReaderFromByte(imgBody))
	if err != nil {
		return nil, err
	}
	var out bytes.Buffer
	err = jpeg.Encode(&out, img, &jpeg.Options{Quality: 100})
	if err != nil {
		return nil, err
	}
	return ConversionBodyFromReader(&out)
}

//ImageFormatType ImageFormatType
type ImageFormatType string

//
const (
	ImageFormatTypeNone ImageFormatType = ""
	ImageFormatTypePNG                  = "png"
	ImageFormatTypeJPEG                 = "jpeg"
	ImageFormatTypeGif                  = "gif"
	ImageFormatTypeBMP                  = "bmp"
)

//ImageToFormat ImageToFormat
func ImageToFormat(subImg image.Image, oldExt string, format ImageFormatType) ([]byte, ImageFormatType, error) {
	var out bytes.Buffer
	if format == ImageFormatTypeJPEG {
		jpeg.Encode(&out, subImg, &jpeg.Options{Quality: 100})
	} else if format == ImageFormatTypePNG {
		png.Encode(&out, subImg)
	} else if format == ImageFormatTypeGif {
		gif.Encode(&out, subImg, &gif.Options{})
	} else if format == ImageFormatTypeBMP {
		bmp.Encode(&out, subImg)
	} else {
		switch oldExt {
		case "jpeg":
			jpeg.Encode(&out, subImg, &jpeg.Options{Quality: 100})
			break
		case "png":
			switch subImg.(type) {
			case *image.NRGBA:
				png.Encode(&out, subImg)
				break
			case *image.RGBA:
				png.Encode(&out, subImg)
				break
			}
			break
		case "gif":
			gif.Encode(&out, subImg, &gif.Options{})
			break
		case "bmp":
			bmp.Encode(&out, subImg)
			break
		default:
			break
		}
	}
	imgBodyCp, err := ConversionBodyFromReader(&out)
	return imgBodyCp, format, err
}

//ImageResizeAutoToFormat ImageResizeAutoToFormat
func ImageResizeAutoToFormat(imgBody []byte, weight, height int, format ImageFormatType) ([]byte, ImageFormatType, error) {
	origin, ext, err := image.Decode(ConversionReaderFromByte(imgBody))
	if err != nil {
		return nil, "", err
	}
	w := origin.Bounds().Dx()
	h := origin.Bounds().Dy()

	kuanW := w
	kuanH := 0
	if w > weight && h > height {
		kuanW = weight
		kuanH = height
	} else if w > weight && h <= height {
		kuanW = weight
		kuanH = 0
	} else if w <= weight && h <= height {
		return ImageToFormat(origin, ext, format)
	} else if w <= weight && h >= height {
		kuanH = height
		kuanW = 0
	}
	subImg := resize.Resize(uint(kuanW), uint(kuanH), origin, resize.Lanczos3)
	return ImageToFormat(subImg, ext, format)
}

func ImageResizeAutoFromImage(origin image.Image, weight, height int) ([]byte, error) {
	w := origin.Bounds().Dx()
	h := origin.Bounds().Dy()
	var out bytes.Buffer
	kuanW := w
	kuanH := 0
	if w > weight && h > height {
		kuanW = weight
		kuanH = height
	} else if w > weight && h <= height {
		kuanW = weight
		kuanH = 0
	} else if w <= weight && h <= height {
		jpeg.Encode(&out, origin, &jpeg.Options{Quality: 100})
		return ConversionBodyFromReader(&out)
	} else if w <= weight && h >= height {
		kuanH = height
		kuanW = 0
	}
	subImg := resize.Resize(uint(kuanW), uint(kuanH), origin, resize.Lanczos3)
	jpeg.Encode(&out, subImg, &jpeg.Options{Quality: 100})
	return ConversionBodyFromReader(&out)
}

func ImageImgResizeAutoFromImage(origin image.Image, weight, height int) (image.Image, error) {
	w := origin.Bounds().Dx()
	h := origin.Bounds().Dy()
	var out bytes.Buffer
	kuanW := w
	kuanH := 0
	if w > weight && h > height {
		kuanW = weight
		kuanH = height
	} else if w > weight && h <= height {
		kuanW = weight
		kuanH = 0
	} else if w <= weight && h <= height {
		jpeg.Encode(&out, origin, &jpeg.Options{Quality: 100})
		return origin, nil
	} else if w <= weight && h >= height {
		kuanH = height
		kuanW = 0
	}
	subImg := resize.Resize(uint(kuanW), uint(kuanH), origin, resize.Lanczos3)
	return subImg, nil
}

//ImageResizeAuto ImageResizeAuto
func ImageResizeAuto(imgBody []byte, weight, height int) ([]byte, error) {
	origin, ext, err := image.Decode(ConversionReaderFromByte(imgBody))
	if err != nil {
		return nil, err
	}
	w := origin.Bounds().Dx()
	h := origin.Bounds().Dy()

	kuanW := w
	kuanH := 0
	if w > weight && h > height {
		kuanW = weight
		kuanH = height
	} else if w > weight && h <= height {
		kuanW = weight
		kuanH = 0
	} else if w <= weight && h <= height {
		return imgBody, nil
	} else if w <= weight && h >= height {
		kuanH = height
		kuanW = 0
	}

	subImg := resize.Resize(uint(kuanW), uint(kuanH), origin, resize.Lanczos3)

	var out bytes.Buffer
	switch ext {
	case "jpeg":
		jpeg.Encode(&out, subImg, &jpeg.Options{Quality: 100})
		break
	case "png":
		switch subImg.(type) {
		case *image.NRGBA:
			png.Encode(&out, subImg)
			break
		case *image.RGBA:
			png.Encode(&out, subImg)
			break
		}
		break
	case "gif":
		gif.Encode(&out, subImg, &gif.Options{})
		break
	case "bmp":
		bmp.Encode(&out, subImg)
		break
	default:
		return []byte{}, errors.New("ERROR FORMAT")
	}
	return ConversionBodyFromReader(&out)
}

//ImageResize ImageResize
func ImageResize(imgBody []byte, isHeight bool) ([]byte, error) {
	origin, ext, err := image.Decode(ConversionReaderFromByte(imgBody))
	if err != nil {
		return nil, err
	}
	kuan := origin.Bounds().Dx()
	if isHeight {
		kuan = origin.Bounds().Dy()
	}
	subImg := resize.Resize(uint(kuan), 0, origin, resize.Lanczos3)
	var out bytes.Buffer
	switch ext {
	case "jpeg":
		jpeg.Encode(&out, subImg, &jpeg.Options{Quality: 0})
		break
	case "png":
		switch subImg.(type) {
		case *image.NRGBA:
			png.Encode(&out, subImg)
			break
		case *image.RGBA:
			png.Encode(&out, subImg)
			break
		}
		break
	case "gif":
		gif.Encode(&out, subImg, &gif.Options{})
		break
	case "bmp":
		bmp.Encode(&out, subImg)
		break
	default:
		return []byte{}, errors.New("ERROR FORMAT")
	}
	return ConversionBodyFromReader(&out)
}

//GetImageInfo GetImageInfo
func GetImageImageThumbnail(file *multipart.FileHeader, witch int, save string) error {
	f, err := file.Open()
	if err != nil {
		return err
	}
	defer f.Close()
	img, err := imaging.Decode(f)
	if err != nil {
		return fmt.Errorf("解码错误:%v size:%v", err, file.Size)
	}

	imh := imaging.Resize(img, witch, 0, imaging.Lanczos)
	err = imaging.Save(imh, save+".jpg")
	if err != nil {
		return fmt.Errorf("保存错误:%v", err)
	}
	return nil
}

//GetImageInfo GetImageInfo
func GetImageImageThumbnailGIF(file *multipart.FileHeader, witch int, save string) error {
	f, err := file.Open()
	if err != nil {
		return err
	}
	defer f.Close()
	img, err := gif.Decode(f)
	if err != nil {
		return fmt.Errorf("解码错误:%v size:%v", err, file.Size)
	}

	imh := imaging.Resize(img, witch, 0, imaging.Lanczos)
	err = imaging.Save(imh, save+".jpg")
	if err != nil {
		return fmt.Errorf("保存错误:%v", err)
	}
	return nil
}

//ImageThumbnail 缩略图
func ImageThumbnail(body []byte, witch int, save string) error {
	img, err := imaging.Decode(ConversionReaderFromByte(body))
	if err != nil {
		return err
	}
	imh := imaging.Resize(img, witch, 0, imaging.Lanczos)
	err = imaging.Save(imh, save+".jpg")
	return err
}

//ImageThumbnail 缩略图
func ImageThumbnailImage(img image.Image, witch int, save string) error {
	imh := imaging.Resize(img, witch, 0, imaging.Lanczos)
	err = imaging.Save(imh, save+".jpg")
	return err
}

//ImageThumbnail 缩略图
func GetImageThumbnailFile(imgFile string, witch int, save string) error {
	imgData, err := ioutil.ReadFile(imgFile)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(imgData)
	img, err := imaging.Decode(buf)
	if err != nil {
		return err
	}
	imh := imaging.Resize(img, witch, 0, imaging.Lanczos)
	err = imaging.Save(imh, save+".jpg")
	return err
}

//GetRatioImageSize 获取图片比例
func GetRatioImageSize(w, h int, maxW, maxH int) (int, int) {

	return 0, 0
}
