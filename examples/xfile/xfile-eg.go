package main

import (
	"fmt"
	"github.com/Bishoptylaor/go-toolkit/xfile"
)

func main() {
	TestDirExists()
}

func TestDirExists() {
	var testCase = []struct {
		path     string
		expected bool
	}{
		{"/Users/gtec/wxy", true},
	}

	for _, tc := range testCase {
		var result = xfile.DirExists(tc.path)
		if result != tc.expected {
			fmt.Printf("文件存在判断 path:[%s], 期望:[%t], 结果:[%t]\n", tc.path, tc.expected, result)
		}
	}
}
