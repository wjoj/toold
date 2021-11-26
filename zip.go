package toold

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

/*
ZipCompressPaths zip压缩 如果文件不存在会滤过 paths:文件 dest:压缩后的文件zip路径
*/
func ZipCompressPaths(paths []string, dest string) error {
	files := []*os.File{}
	for _, p := range paths {
		f, err := os.Open(p)
		if err != nil {
			continue
		}
		files = append(files, f)
	}
	return ZipCompressFiles(files, dest)
}

/*
ZipCompressFiles zip压缩 files:要压缩的文件 dest:压缩后的文件zip路径
*/
func ZipCompressFiles(files []*os.File, dest string) error {
	d, _ := os.Create(dest)
	defer d.Close()
	w := zip.NewWriter(d)
	defer w.Close()
	for _, file := range files {
		defer file.Close()
		err := compress(file, "", w)
		if err != nil {
			return err
		}
	}
	return nil
}
func compress(file *os.File, prefix string, zw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			defer f.Close()
			err = compress(f, prefix, zw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func Unzip(zipFile string, destDir string, infoFunc func(fileNumber int, progress int)) error {
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer zipReader.Close()
	fileNumber := len(zipReader.File)
	if infoFunc != nil && fileNumber == 0 {
		infoFunc(fileNumber, 0)
		return nil
	}
	for i, f := range zipReader.File {
		if f.Flags == 0 {
			i := bytes.NewReader([]byte(f.Name))
			decoder := transform.NewReader(i, simplifiedchinese.GB18030.NewDecoder())
			content, _ := ioutil.ReadAll(decoder)
			f.Name = string(content)
		}
		fpath := filepath.Join(destDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			inFile, err := f.Open()
			if err != nil {
				return err
			}
			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				inFile.Close()
				return err
			}
			buf := make([]byte, 1024*100)
			for {
				n, err := inFile.Read(buf)
				if err != nil && err != io.EOF {
					return err
				}
				if n == 0 {
					break
				}
				_, err2 := outFile.Write(buf[:n])
				if err2 != nil {
					return err2
				}
				if err == io.EOF {
					break
				}
			}
			outFile.Close()
			inFile.Close()
		}
		if infoFunc != nil {
			infoFunc(len(zipReader.File), i)
		}
	}
	return nil
}
