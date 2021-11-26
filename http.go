package toold

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

//HTTPDataValues HTTPDataValues
var HTTPDataValues = func(data url.Values) io.Reader {
	return ioutil.NopCloser(strings.NewReader(data.Encode()))
}

// HTTPRequestData 请求返回的数据结构
type HTTPRequestData struct {
	HTTPCode int
	Data     []byte
	Header   http.Header
}

//HTTPDataString HTTPDataString
func HTTPDataString(val string) io.Reader {
	return ioutil.NopCloser(strings.NewReader(val))
}

/*
HTTPUploadFile 表单上传方式
*/
func HTTPUploadFileCode(formKeyName string, filename, url string) ([]byte, int, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile(formKeyName, filename)
	if err != nil {
		fmt.Println("写入表单文件失败")
		return nil, 20001, err
	}
	//打开文件句柄操作
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("打开文件失败")
		return nil, 20001, err
	}
	defer fh.Close()
	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return nil, 20001, err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(url, contentType, bodyBuf)
	if err != nil {
		return nil, 20001, err
	}
	defer resp.Body.Close()
	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	return respbody, resp.StatusCode, nil
}

/*
HTTPUploadFile 表单上传方式
*/
func HTTPUploadFile(formKeyName string, filename, url string) ([]byte, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile(formKeyName, filename)
	if err != nil {
		fmt.Println("写入表单文件失败")
		return nil, err
	}
	//打开文件句柄操作
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("打开文件失败")
		return nil, err
	}
	defer fh.Close()
	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return nil, err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(url, contentType, bodyBuf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return respbody, nil
}

/*
HTTPUploadFileHeader 表单上传方式 和其他值
*/
func HTTPUploadFileAndValueHeaderCode(formKeyName, filename string, file *multipart.FileHeader, vals map[string]interface{}, url string) ([]byte, int, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	for k, val := range vals {
		bodyWriter.WriteField(k, fmt.Sprintf("%v", val))
	}
	if file != nil {
		//关键的一步操作
		fileWriter, err := bodyWriter.CreateFormFile(formKeyName, filename)
		if err != nil {
			fmt.Println("写入表单文件失败")
			return nil, 10001, err
		}
		//打开文件句柄操作
		fh, err := file.Open()
		if err != nil {
			fmt.Println("打开文件失败")
			return nil, 10002, err
		}
		defer fh.Close()
		//iocopy
		_, err = io.Copy(fileWriter, fh)
		if err != nil {
			return nil, 10003, err
		}
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(url, contentType, bodyBuf)
	if err != nil {
		return nil, 10004, err
	}
	defer resp.Body.Close()
	respbody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, resp.StatusCode, err
	}
	return respbody, resp.StatusCode, nil
}

/*
HTTPUploadFileHeader 表单上传方式
*/
func HTTPUploadFileHeaderCode(formKeyName, filename string, file *multipart.FileHeader, url string) ([]byte, int, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile(formKeyName, filename)
	if err != nil {
		fmt.Println("写入表单文件失败")
		return nil, 10001, err
	}
	//打开文件句柄操作
	fh, err := file.Open()
	if err != nil {
		fmt.Println("打开文件失败")
		return nil, 10002, err
	}
	defer fh.Close()
	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return nil, 10003, err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	// http.DefaultClien
	resp, err := http.Post(url, contentType, bodyBuf)
	if err != nil {
		return nil, 10004, err
	}
	defer resp.Body.Close()
	respbody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, resp.StatusCode, err
	}
	return respbody, resp.StatusCode, nil
}

/*
HTTPUploadFileHeader 表单上传方式
*/
func HTTPUploadFileHeader(formKeyName, filename string, file *multipart.FileHeader, url string) ([]byte, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile(formKeyName, filename)
	if err != nil {
		fmt.Println("写入表单文件失败")
		return nil, err
	}
	//打开文件句柄操作
	fh, err := file.Open()
	if err != nil {
		fmt.Println("打开文件失败")
		return nil, err
	}
	defer fh.Close()
	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return nil, err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(url, contentType, bodyBuf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respbody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}
	return respbody, nil
}

func HTTPUploadForm(urls string, timeOut time.Duration, formKeyName, filename string, heard map[string]string, file *multipart.FileHeader) (*HTTPRequestData, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile(formKeyName, filename)
	if err != nil {
		fmt.Println("写入表单文件失败")
		return nil, err
	}
	//打开文件句柄操作
	fh, err := file.Open()
	if err != nil {
		fmt.Println("打开文件失败")
		return nil, err
	}
	defer fh.Close()
	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return nil, err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	heard["Content-Type"] = contentType
	heard["Connection"] = "close"
	return HTTPSetHeaderRequest("POST", timeOut, heard, urls, bodyBuf)
}

/*
HTTPSetHeaderRequest 头设置请求
*/
func HTTPSetHeaderRequest(method string, timeOut time.Duration, heard map[string]string, urls string, body io.Reader) (*HTTPRequestData, error) {
	var reqest *http.Request
	var err error
	client := &http.Client{
		Timeout: timeOut * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	reqest, err = http.NewRequest(method, urls, body)
	for keys, val := range heard {
		reqest.Header.Set(keys, val)
	}
	if err != nil {
		return nil, err
	}
	response, err := client.Do(reqest)
	if err != nil {
		return nil, err
	}

	status := response.StatusCode

	byt, err := ioutil.ReadAll(response.Body)
	defer func() {
		response.Body.Close()
		client.CloseIdleConnections()
	}()

	if err != nil {
		return nil, err
	}

	return &HTTPRequestData{
		HTTPCode: status,
		Data:     byt,
		Header:   response.Header,
	}, err
}

/*
HTTPSetHeadersRequest 头设置请求
*/
func HTTPSetHeadersRequest(method string, timeOut time.Duration, head func(reqh *http.Request), urls string, body io.Reader) (*HTTPRequestData, error) {
	var reqest *http.Request
	var err error
	client := &http.Client{
		Timeout: timeOut * time.Second,
	}
	reqest, err = http.NewRequest(method, urls, body)
	if reqest != nil && head != nil {
		head(reqest)
	}
	if err != nil {
		return nil, err
	}

	response, err := client.Do(reqest)
	if err != nil {
		return nil, err
	}

	status := response.StatusCode

	byt, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return nil, err
	}

	return &HTTPRequestData{
		HTTPCode: status,
		Data:     byt,
		Header:   response.Header,
	}, err
}

//HTTPMethod HTTPMethod
func HTTPMethod(urls string, method string, timeOut time.Duration, data url.Values) (*HTTPRequestData, error) {
	return HTTPSetHeaderRequest(method, timeOut, map[string]string{}, urls, HTTPDataValues(data))
}

func HTTPMethodStr(urls string, method string, timeOut time.Duration, data string) (*HTTPRequestData, error) {
	return HTTPSetHeaderRequest(method, timeOut, map[string]string{}, urls, HTTPDataString(data))
}

func HTTPMethodHeaderStr(urls string, method string, timeOut time.Duration, headr map[string]string, data string) (*HTTPRequestData, error) {
	return HTTPSetHeaderRequest(method, timeOut, headr, urls, HTTPDataString(data))
}

//HttpDownloadFile 下载资源文件
func HttpDownloadFile(method string, urls string, timeOut time.Duration, body io.Reader, save, newFileName string, proFunc func(fileSize, progress int64)) error {
	var reqest *http.Request
	var err error
	client := &http.Client{
		Timeout: timeOut * time.Second,
	}
	reqest, err = http.NewRequest(method, urls, body)
	if err != nil {
		return fmt.Errorf("%v %v", err, urls)
	}
	response, err := client.Do(reqest)
	if err != nil {
		return fmt.Errorf("%v %v", err, urls)
	}
	defer func() {
		response.Body.Close()
		client.CloseIdleConnections()
	}()

	status := response.StatusCode
	if status != 200 {
		return fmt.Errorf("%v(%v) %v", response.Status, status, urls)
	}

	size := ConversionToInt64(response.Header.Get("Content-Length"))
	num := size / (1024 * 100)
	if num > 10000 {
		num = 100
	} else if num <= 10000 && num > 100 {
		num = 50
	} else {
		num = 50
	}
	outFile, err := os.OpenFile(fmt.Sprintf("%v/%v", save, newFileName), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer outFile.Close()
	buf := make([]byte, 1024*50)
	if proFunc != nil {
		proFunc(size, 0)
	}
	lng := int64(0)
	for {
		n, err := response.Body.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		lng += int64(n)
		if proFunc != nil {
			proFunc(size, lng)
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
	if proFunc != nil {
		proFunc(size, lng)
	}

	return nil
}
