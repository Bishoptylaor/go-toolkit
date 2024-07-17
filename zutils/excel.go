package zutils

import (
	"bytes"
	"context"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"io"
	"net/http"
)

/*
 *  ┏┓      ┏┓
 *┏━┛┻━━━━━━┛┻┓
 *┃　　　━　　  ┃
 *┃   ┳┛ ┗┳   ┃
 *┃           ┃
 *┃     ┻     ┃
 *┗━━━┓     ┏━┛
 *　　 ┃　　　┃神兽保佑
 *　　 ┃　　　┃代码无BUG！
 *　　 ┃　　　┗━━━┓
 *　　 ┃         ┣┓
 *　　 ┃         ┏┛
 *　　 ┗━┓┓┏━━┳┓┏┛
 *　　   ┃┫┫  ┃┫┫
 *      ┗┻┛　 ┗┻┛
 @Time    : 2024/7/13 -- 14:17
 @Author  : bishop ❤️ MONEY
 @Description: excel
*/

type excelTool struct{}

var ExcelTool excelTool

// GetIOReaderFromUrl 通过url获取io Reader数据流
func (et excelTool) getIOReaderFromUrl(url string) (reader io.ReadCloser, err error) {
	fun := "GetIOReaderFromUrl"
	load := func() (io.ReadCloser, error) {
		httpRes, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		// defer httpRes.Body.Close()

		if httpRes.StatusCode != 200 {
			httpRes.Body.Close()
			return nil, fmt.Errorf("%s Http.Get not 200 error", fun)
		}
		return httpRes.Body, nil
	}
	for i := 1; i <= 3; i++ {
		if reader, err = load(); err == nil && reader != nil {
			break
		}
	}
	return
}

// GetExcelFileFromUrl 通过url获取并打开Excel文件
func (et excelTool) GetExcelFileFromUrl(url string) (f *excelize.File, err error) {
	fun := "GetExcelFileFromUrl"
	if url == "" {
		return nil, fmt.Errorf("%s url为空", fun)
	}
	ioReader, err := et.getIOReaderFromUrl(url)
	if err != nil {
		return nil, err
	}
	defer ioReader.Close()
	f, err = excelize.OpenReader(ioReader)
	if err != nil {
		return nil, fmt.Errorf("%s open ExcelFile error,err:%s", fun, err)
	}
	return f, nil
}

// CreateNewExcel 创建一个新的Excel文件
func (et excelTool) CreateNewExcel() (f *excelize.File) {
	return excelize.NewFile()
}

// UploadExcelFile 上传Excel文件到CDN中,并获取URL
func (et excelTool) UploadExcelFile(ctx context.Context, f *excelize.File,
	uploader func(ctx context.Context, buf *bytes.Buffer, args ...any) (string, error),
	urlFmt func(ctx context.Context, srcUrl string) (dstUrl string),
	args ...any) (url string, err error) {

	writeBuf, err := f.WriteToBuffer()
	if err != nil {
		return "", err
	}

	url, err = uploader(ctx, writeBuf, args)
	if err != nil {
		return "", err
	}
	return urlFmt(ctx, url), nil
}
