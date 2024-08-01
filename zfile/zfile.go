package zfile

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/Bishoptylaor/go-toolbox/zutils"
	"io"
	"os"
	"strings"
	"time"
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
 @Time    : 2024/7/13 -- 14:44
 @Author  : bishop ❤️ MONEY
 @Description: 常用文件操作封装
*/

const IO_BUFFER_SIZE = 4096

var (
	ErrEmptyArguments = errors.New("argument Can't be empty")
)

/*
FileExists
判断文件是否存在，该函数和PHP的file_exists一致，当filename 为 文件/文件夹/符号链接 时均返回 true
*/
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

/*
IsFile
判断给定的filename是否是一个文件
*/
func IsFile(filename string) (bool, error) {

	if len(filename) <= 0 {
		return false, ErrEmptyArguments
	}

	stat, err := os.Stat(filename)
	if err != nil {
		return false, err
	}

	if stat.IsDir() {
		return false, nil
	}

	return true, nil
}

// FileSize 获取文件大小
func FileSize(filename string) (int64, error) {
	if len(filename) <= 0 {
		return 0, ErrEmptyArguments
	}

	stat, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}

	if !stat.IsDir() {
		return stat.Size(), nil
	}
	return 0, os.ErrInvalid
}

// CopyFile 将文件从src Copy 到 dst
func CopyFile(src string, dst string) error {
	sfp, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sfp.Close()

	dfp, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	_, err = io.Copy(dfp, sfp)
	if err != nil {
		_ = dfp.Close()
		return err
	}
	return dfp.Close()
}

/*
GetFileModifyTime
获取文件的最后修改时间
*/
func GetFileModifyTime(filename string) (time.Time, error) {

	fileInfo, err := os.Stat(filename)
	if err != nil {
		return time.Time{}, err
	}

	return fileInfo.ModTime(), nil
}

func SearchFileInDir(filename string, dirs ...string) (fullPath string, err error) {
	for _, dir := range dirs {
		if strings.HasSuffix(dir, "/") {
			fullPath = dir + filename
		} else {
			fullPath = dir + "/" + filename
		}
		if FileExists(fullPath) {
			return
		}
	}
	fullPath = ``
	err = os.ErrNotExist
	return
}

/*
CalFileLineCount
计算一个文件的行数
*/
func CalFileLineCount(pathname string) (cnt uint32, err error) {
	var fp *os.File
	fp, err = os.Open(pathname)
	if err != nil {
		return
	}
	defer fp.Close()

	var buf []byte = make([]byte, 1024)
	var rdSize int
	for {
		rdSize, err = fp.Read(buf)
		if rdSize > 0 {
			for i := 0; i < rdSize; i++ {
				if buf[i] == '\n' {
					cnt++
				}
			}
		} else {
			if err == io.EOF {
				err = nil
				break
			} else {
				return 0, err
			}
		}
	}
	return
}

// LineCallback 行回调函数，当该函数返回 false 时，回调终止，不再继续
type LineCallback func(line string) (doContinue bool)

// EachLine 对文件中每一行内容进行回调,当 lc 返回 false 时回调也会终止
func EachLine(pathname string, lc LineCallback, skipEmpty bool) (err error) {
	var fp *os.File
	fp, err = os.Open(pathname)
	if err != nil {
		return
	}
	defer fp.Close()

	rd := bufio.NewReader(fp)
	for {
		line, err := rd.ReadString('\n')
		if err == nil {
			line = strings.TrimRightFunc(line, func(r rune) bool {
				if r == '\r' || r == '\n' {
					return true
				}
				return false
			})
			if skipEmpty {
				if line == `` {
					continue
				}
			}
			if !lc(line) {
				break
			}
		} else {
			if err == io.EOF {
				err = nil
				break
			} else {
				return err
			}
		}
	}
	return
}

func ReadAll(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

func WriteAll(filePath string, bs []byte) error {
	return os.WriteFile(filePath, bs, 0666)
}

func ReadAllString(filePath string) (string, error) {
	b, err := ReadAll(filePath)
	if err != nil {
		return "", err
	}
	return zutils.Bytes2Str(b), nil
}

// ReadLines 通过 channel 读取每行
func ReadLines(ctx context.Context, fn string) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)

		f, _ := os.Open(fn)
		defer f.Close()
		r := bufio.NewReader(f)

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			line, err := r.ReadString('\n')
			if errors.Is(err, io.EOF) {
				break
			}
			if line == "" {
				continue
			}
			line = line[:len(line)-1]
			if line == "" {
				continue
			}

			ch <- line
		}
	}()
	return ch
}

func PutFile(filePath string, contents []byte, append bool) error {
	var (
		out *os.File
		err error
	)
	if append && FileExists(filePath) {
		out, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	} else {
		out, err = os.Create(filePath)
	}
	if err != nil {
		return err
	}
	defer out.Close()
	return ensureWrite(out, contents)
}

func ensureWrite(out io.Writer, bytes []byte) error {
	var (
		total     = len(bytes)
		wrote int = 0
	)
	for {
		if wrote >= total {
			return nil
		}
		n, err := out.Write(bytes[wrote:])
		if err != nil {
			return err
		}
		if n <= 0 {
			return fmt.Errorf("cann't write bytes to io.Writer, it writes %d bytes", n)
		}
		wrote += n
	}
	return nil
}
