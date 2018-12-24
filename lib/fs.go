// Copyright © 2018 Mephis Pheies <mephistommm@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
package lib

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

// FileLineParser defines which file to be parsed and how to parse it.
type FileLineParser interface {
	// FilePath return the path of the file to be parsed
	FilePath() string
	// Parse accepts a line as parameter and should parse
	// it to a struct
	Parse(line string) (interface{}, error)
}

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, copy the file contents from src to dst.
func CopyFile(src, dst string) (err error) {
	Logger.Debugf("Copy file %s to %s...", src, dst)
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	err = copyFileContents(src, dst)
	Logger.Debugf("Finish Copying from %s to %s...", src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	buf := make([]byte, os.Getpagesize())
	for {
		n, err := in.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := out.Write(buf[:n]); err != nil {
			return err
		}
	}
	err = out.Sync()
	return
}

func IsDir(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}

func IsNotExist(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}

// WriteFile write data to file, and create its directories if necessary
func WriteFile(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		Logger.Fatal(err)
	}
	return ioutil.WriteFile(path, data, 0664)
}

// ReadFile just keep the same style as WriteFile
func ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

// ScanLines read file line by line and call Parse method of FileLineParser
// for each line.
//
// This function only collect non-nil result returned by Parser().
func ScanLines(parser FileLineParser) ([]interface{}, error) {
	fileHandle, _ := os.Open(parser.FilePath())
	defer fileHandle.Close()

	fileScanner := bufio.NewScanner(fileHandle)
	results := make([]interface{}, 0)
	for fileScanner.Scan() {
		result, err := parser.Parse(fileScanner.Text())
		if err != nil {
			return nil, err
		}

		if result != nil {
			results = append(results, result)
		}
	}

	return results, nil
}

func HomeDir() string {
	usr, err := user.Current()
	if err != nil {
		Logger.Fatal(err)
	}
	return usr.HomeDir
}
