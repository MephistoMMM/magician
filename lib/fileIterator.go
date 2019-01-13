// Copyright © 2019 Mephis Pheies <mephistommm@gmail.com>
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
	"io/ioutil"
	"os"
	"path/filepath"
)

type defaultFileIterator struct {
	dir string

	index int
	files []os.FileInfo
	dirs  []string

	filter FilterSupport
}

// Init ...
func (dfi *defaultFileIterator) Init(dir string) error {
	return dfi.loadFilesAndDirs(dir, dfi.files[:0], dfi.dirs[:0])
}

// loadFilesAndDirs ...
func (dfi *defaultFileIterator) loadFilesAndDirs(dir string, files []os.FileInfo, dirs []string) error {
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range fs {
		path := filepath.Join(dir, file.Name())
		if !Filter(dfi.filter, path, file) {
			continue
		}

		if file.IsDir() {
			dirs = append(dirs, path)
		} else {
			files = append(files, file)
		}
	}

	dfi.dir = dir
	dfi.dirs = dirs
	dfi.files = files
	dfi.index = 0
	return nil
}

// HasNext ...
func (dfi *defaultFileIterator) HasNext() bool {
	if dfi.index < len(dfi.files) {
		return true
	}

	if len(dfi.dirs) > 0 {
		// read files and directories under the first dir element
		dir := dfi.dirs[0]
		dfi.dirs = dfi.dirs[1:]
		err := dfi.loadFilesAndDirs(dir, dfi.files[:0], dfi.dirs[:])
		if err != nil {
			panic(err)
		}

		return dfi.HasNext()
	}

	return false
}

// Next ...
func (dfi *defaultFileIterator) Next() (path string, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
			return
		}
	}()

	if !dfi.HasNext() {
		return "", nil
	}

	result := filepath.Join(dfi.dir, dfi.files[dfi.index].Name())
	dfi.index++
	return result, nil
}

func defaultFilterChain() FilterSupport {
	var filterChain FilterSupport
	tmp, _ := NewFilterIgnoreDotSupport()
	filterChain = tmp
	tmp, _ = NewFilterIgnoreUnregularSupport()
	filterChain.SetNext(tmp)
	return filterChain
}

func NewFileIteratorWithDefaultFilter(directory string) (FileIterator, error) {

	iterator := &defaultFileIterator{
		index:  0,
		files:  make([]os.FileInfo, 0),
		dirs:   make([]string, 0),
		filter: defaultFilterChain(),
	}

	if err := iterator.Init(directory); err != nil {
		return nil, err
	}

	return iterator, nil
}

// NewFileIterator create a new file iterator.
func NewFileIterator(directory string, filterChain FilterSupport) (FileIterator, error) {
	iterator := &defaultFileIterator{
		index:  0,
		files:  make([]os.FileInfo, 0),
		dirs:   make([]string, 0),
		filter: filterChain,
	}

	if err := iterator.Init(directory); err != nil {
		return nil, err
	}

	return iterator, nil
}
