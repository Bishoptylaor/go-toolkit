package zfile

import (
	"errors"
	"fmt"
	"github.com/Bishoptylaor/go-toolbox/zcontainer"
	"os"
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
 @Time    : 2024/7/17 -- 14:44
 @Author  : bishop ❤️ MONEY
 @Description: 常用文件操作封装 dir
*/

// DirExists 判断文件路径是否存在
func DirExists(p string) bool {
	if fi, err := os.Stat(p); err == nil {

		return fi.IsDir()
	}
	return false
}

/*
IsDir
判断给定的filename是否是一个目录
*/
func IsDir(filename string) (bool, error) {

	if len(filename) <= 0 {
		return false, ErrEmptyArguments
	}

	stat, err := os.Stat(filename)
	if err != nil {
		return false, err
	}

	if !stat.IsDir() {
		return false, nil
	}

	return true, nil
}

// DirEmpty 判断文件路径是否为空文件夹
func DirEmpty(p string) bool {
	fs, err := os.ReadDir(p)
	if err != nil {
		return false
	}
	return len(fs) == 0
}

// DirSize 递归获取一个文件夹下所有文件的大小
func DirSize(filename string) (int64, error) {
	if len(filename) <= 0 {
		return 0, ErrEmptyArguments
	}
	var (
		size = int64(0)
	)
	err := EachFile(filename, func(path string, fio os.FileInfo) error {
		if !fio.IsDir() {
			fsz, e := FileSize(path)
			if e != nil {
				return e
			}
			size += fsz
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return size, nil
}

/*
DirAppendFileName
给路径dir 之后加上 name，该函数主要 fmt 处理 path 尾部的 / 和 name 最前面的 / 的问题
*/
func DirAppendFileName(dir, name string) string {
	if dir[len(dir)-1] == '/' {
		if name[0] == '/' {
			return dir + name[1:]
		}
		return dir + name
	}
	if name[0] == '/' {
		return dir + name
	}
	return dir + "/" + name
}

type FileInfoCallback func(pathname string, fi os.FileInfo) error

/*
EachFile
递归的遍历某个文件/文件夹中的所有对象，
*/
func EachFile(filepath string, fic FileInfoCallback) error {
	if filepath == `` {
		return ErrEmptyArguments
	}
	fi, err := os.Stat(filepath)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return errors.New("filepath is not dir")
	}
	var stack zcontainer.Stack
	stack.Push(filepath)
	for {
		if stack.Len() <= 0 {
			break
		}
		se := stack.Pop()
		if se == nil {
			break
		}
		if fpath, ok := se.Value.(string); ok {
			fp, err := os.Open(fpath)
			if err != nil {
				fmt.Printf("try to iterate path %s but failed, error is: %#v\n", fpath, err)
				continue
			}
			dirs, err := fp.Readdir(-1)
			if err != nil {
				fmt.Printf("try to read %s sub items but failed, error is: %#v\n", fpath, err)
				fp.Close()
				continue
			}
			fp.Close()
			for _, sfi := range dirs {
				pathname := DirAppendFileName(fpath, sfi.Name())
				if sfi.IsDir() {
					stack.Push(pathname)
				}
				err := fic(pathname, sfi)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
