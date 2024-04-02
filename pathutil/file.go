package pathutil

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/errors"
	"github.com/pubgo/funk/recovery"
)

// IsNotExist check if the file exists
func IsNotExist(src string) bool {
	_, err := os.Stat(src)
	return os.IsNotExist(err)
}

// CheckPermission check if the file has permission
func CheckPermission(src string) bool {
	_, err := os.Stat(src)
	return os.IsPermission(err)
}

// IsNotExistMkDir create a directory if it does not exist
func IsNotExistMkDir(src string) (err error) {
	defer recovery.Err(&err)
	if IsNotExist(src) {
		assert.Must(MkDir(src), "MkDir Error")
	}

	return
}

// MkDir create a directory
func MkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	return errors.WrapKV(err, "src", src)
}

// IsExist determines whether the file spcified by the given path is exists.
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// IsDir determines whether the specified path is a directory.
func IsDir(path string) bool {
	fio, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return false
	}

	if nil != err {
		return false
	}

	return fio.IsDir()
}

// CopyFile copies the source file to the dest file.
func CopyFile(source string, dest string) (err error) {
	defer recovery.Err(&err)

	sourcefile := assert.Must1(os.Open(source))
	defer sourcefile.Close()

	destfile := assert.Must1(os.Create(dest))
	defer destfile.Close()

	_ = assert.Must1(io.Copy(destfile, sourcefile))

	sourceinfo := assert.Must1(os.Stat(source))
	assert.Must(os.Chmod(dest, sourceinfo.Mode()))

	return
}

// CopyDir copies the source directory to the dest directory.
func CopyDir(source string, dest string) (err error) {
	defer recovery.Err(&err)

	sourceinfo := assert.Must1(os.Stat(source))

	// create dest dir
	assert.Must(os.MkdirAll(dest, sourceinfo.Mode()))

	directory := assert.Must1(os.Open(source))
	defer directory.Close()

	objects := assert.Must1(directory.Readdir(-1))
	for _, obj := range objects {
		srcFilePath := filepath.Join(source, obj.Name())
		destFilePath := filepath.Join(dest, obj.Name())

		if obj.IsDir() {
			assert.Must(CopyDir(srcFilePath, destFilePath))
		} else {
			assert.Must(CopyFile(srcFilePath, destFilePath))
		}
	}

	return
}

// GrepFile like command grep -E
// for example: GrepFile(`^hello`, "hello.txt")
// \n is striped while read
func GrepFile(patten string, filename string) (lines []string, err error) {
	re, err := regexp.Compile(patten)
	if err != nil {
		return
	}

	fd, err := os.Open(filename)
	if err != nil {
		return
	}
	lines = make([]string, 0)
	reader := bufio.NewReader(fd)
	prefix := ""
	var isLongLine bool
	for {
		byteLine, isPrefix, er := reader.ReadLine()
		if er != nil && er != io.EOF {
			return nil, er
		}
		if er == io.EOF {
			break
		}
		line := string(byteLine)
		if isPrefix {
			prefix += line
			continue
		} else {
			isLongLine = true
		}

		line = prefix + line
		if isLongLine {
			prefix = ""
		}
		if re.MatchString(line) {
			lines = append(lines, line)
		}
	}
	return lines, nil
}

// CheckFileIsExist 检查目录是否存在
func CheckFileIsExist(filename string) bool {
	exist := true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// BuildDir 创建目录
func BuildDir(absDir string) error {
	return os.MkdirAll(path.Dir(absDir), os.ModePerm) // 生成多级目录
}

// DeleteFile 删除文件或文件夹
func DeleteFile(absDir string) error {
	return os.RemoveAll(absDir)
}

// GetPathDirs 获取目录所有文件夹
func GetPathDirs(absDir string) (re []string) {
	if CheckFileIsExist(absDir) {
		files, _ := os.ReadDir(absDir)
		for _, f := range files {
			if f.IsDir() {
				re = append(re, f.Name())
			}
		}
	}
	return
}

// GetPathFiles 获取目录所有文件
func GetPathFiles(absDir string) (re []string) {
	if CheckFileIsExist(absDir) {
		files, _ := os.ReadDir(absDir)
		for _, f := range files {
			if !f.IsDir() {
				re = append(re, f.Name())
			}
		}
	}
	return
}

// GetModelPath 获取目录地址
func GetModelPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path := filepath.Dir(file)
	path, _ = filepath.Abs(path)

	return path
}

// GetCurrentDirectory 获取程序运行路径
func GetCurrentDirectory() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return strings.Replace(dir, "\\", "/", -1)
}

func DirName(argv ...string) string {
	file := ""
	if len(argv) > 0 && argv[0] != "" {
		file = argv[0]
	} else {
		file, _ = exec.LookPath(os.Args[0])
	}
	path, _ := filepath.Abs(file)
	directory := filepath.Dir(path)
	return strings.Replace(directory, "\\", "/", -1)
}

func GetProPath() string {
	return DirName("root")
}

// List list file
func List(dirPth, suffix string) (files []string, err error) {
	files = make([]string, 0, 10)
	dir, err := os.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix)
	for _, fi := range dir {
		if fi.IsDir() {
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			files = append(files, dirPth+PthSep+fi.Name())
		}
	}

	return files, nil
}

// ListDir list dir
func ListDir(dirPth, suffix string) (files []string, err error) {
	files = make([]string, 0, 10)
	dir, err := os.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix)
	for _, fi := range dir {
		if !fi.IsDir() {
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			files = append(files, dirPth+PthSep+fi.Name())
		}
	}

	return files, nil
}
