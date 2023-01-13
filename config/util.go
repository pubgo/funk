package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pubgo/funk/assert"
)

const pkgKey = "name"
const defaultKey = "default"

func getPkgId(m map[string]interface{}) string {
	if m == nil {
		return defaultKey
	}

	var val, ok = m[pkgKey]
	if !ok || val == nil {
		return defaultKey
	}

	return fmt.Sprintf("%v", val)
}

// getPathList 递归得到当前目录到跟目录中所有的目录路径
//
//	[./, ../, ../../, ..., /]
func getPathList() (paths []string) {
	var wd = assert.Must1(filepath.Abs(""))
	for {
		if len(wd) == 0 || os.IsPathSeparator(wd[len(wd)-1]) {
			break
		}

		paths = append(paths, wd)
		wd = filepath.Dir(wd)
	}
	return
}

func strMap(strList []string, fn func(str string) string) []string {
	for i := range strList {
		strList[i] = fn(strList[i])
	}
	return strList
}
