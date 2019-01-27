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
	"errors"

	"github.com/MephistoMMM/grafter/util"
	yaml "gopkg.in/yaml.v2"
)

var (
	// ErrYamlFileNotExist ...
	ErrYamlFileNotExist = errors.New("Yaml File Not Exist")
)

// YamlLoadFromFile unmarshal model from a yaml file
func YamlLoadFromFile(path string, model interface{}) error {
	if util.IsNotExist(path) {
		return ErrYamlFileNotExist
	}

	d, err := util.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(d, model)
	if err != nil {
		return err
	}

	return nil
}
