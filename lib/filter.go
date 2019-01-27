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
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// Filter files is designed as `chain of repositories` mode. It includes three parts,
// 1.a abstract Support class implemented by FilterSupport interface and
// BaseSupport struct , 2.some concrete Support, 3.Filter function with FilterSupport
// parameter calling each IsIgnore method of each FilterSupport in chain of repositories.

// FilterSupport is a part of abstract Support.
type FilterSupport interface {
	SetNext(FilterSupport) FilterSupport
	SetNexts([]FilterSupport) FilterSupport
	Next() FilterSupport
	String() string

	IsIgnore(path string, info os.FileInfo) (bool, error)
	Done(path string, info os.FileInfo)
	Fail(path string, info os.FileInfo)
}

// BaseSupport implements almost methods of FilterSupport interface except
// `IsIgnore`. It should be embeded into a concrete FilterSupport.
type BaseSupport struct {
	next FilterSupport
	name string
}

// SetName set n to name
// This method should only be used when concrete FilterSupport inits, so it
// is not a part of FilterSupport interface.
func (bs *BaseSupport) SetName(n string) {
	bs.name = n
}

// SetNext assign another FilterSupport to inner next field. It is used to
// construct a chain of FilterSupports.
func (bs *BaseSupport) SetNext(n FilterSupport) FilterSupport {
	bs.next = n
	return n
}

// SetNexts assign a FilterSupport link list to inner next field.
func (bs *BaseSupport) SetNexts(ns []FilterSupport) FilterSupport {
	if len(ns) < 1 {
		return nil
	}
	innerBs := bs.SetNext(ns[0])
	ns = ns[1:]
	for _, n := range ns {
		innerBs = innerBs.SetNext(n)
	}
	return innerBs
}

// Next return next.
func (bs *BaseSupport) Next() FilterSupport {
	return bs.next
}

// String describe the chain of FilterSupport
func (bs *BaseSupport) String() string {
	if bs.Next() == nil {
		return bs.name
	}
	return bs.name + " | " + bs.Next().String()
}

// Done does nothing but implement FilterSupport interface
func (bs *BaseSupport) Done(path string, info os.FileInfo) {
	Logger.Debugf("%s ignored by %s\n", path, bs.name)
}

// Fail does nothing but implement FilterSupport interface
func (bs *BaseSupport) Fail(path string, info os.FileInfo) {
}

// Filter calls each IsIgnore method of each FilterSupport in chain of repositories.
func Filter(chain FilterSupport, path string, info os.FileInfo) bool {
	if chain == nil {
		return true
	}

	for chain != nil {
		result, err := chain.IsIgnore(path, info)
		if err != nil {
			chain.Fail(path, info)
		}

		if result {
			chain.Done(path, info)
			return false
		}

		chain = chain.Next()
	}

	return true
}

type RegexpMatchSupport struct {
	BaseSupport

	pattern *regexp.Regexp
}

func NewMultiRegexpMatchSupports(exprs []string) ([]FilterSupport, error) {
	iss := make([]FilterSupport, 0, len(exprs))
	for _, exprs := range exprs {
		is, err := NewFilterRegexpMatchSupport(exprs)
		if err != nil {
			return nil, err
		}

		iss = append(iss, is)
	}
	return iss, nil
}

func NewFilterRegexpMatchSupport(expr string) (FilterSupport, error) {
	pattern, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}

	is := &RegexpMatchSupport{
		pattern: pattern,
	}
	is.SetName(fmt.Sprintf("RegexpMatchSupport[%s]", expr))
	return is, nil
}

// IsIgnore ...
func (rms *RegexpMatchSupport) IsIgnore(path string, info os.FileInfo) (bool, error) {
	if rms.pattern.MatchString(path) {
		return false, nil
	}

	return true, nil
}

type IgnoreRegexpMatchSupport struct {
	BaseSupport

	pattern *regexp.Regexp
}

func NewMultiIgnoreRegexpMatchSupports(exprs []string) ([]FilterSupport, error) {
	iss := make([]FilterSupport, 0, len(exprs))
	for _, exprs := range exprs {
		is, err := NewFilterIgnoreRegexpMatchSupport(exprs)
		if err != nil {
			return nil, err
		}

		iss = append(iss, is)
	}
	return iss, nil
}

func NewFilterIgnoreRegexpMatchSupport(expr string) (FilterSupport, error) {
	pattern, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}

	is := &IgnoreRegexpMatchSupport{
		pattern: pattern,
	}
	is.SetName("IgnoreRegexpMatchSupport")
	return is, nil
}

// IsIgnore ...
func (irms *IgnoreRegexpMatchSupport) IsIgnore(path string, info os.FileInfo) (bool, error) {
	if irms.pattern.MatchString(path) {
		return true, nil
	}
	return false, nil
}

// IgnoreSpecialModeSupport ignore special type fies.
type IgnoreSpecialModeSupport struct {
	BaseSupport

	modeMask os.FileMode
}

func NewFilterIgnoreUnregularSupport() (FilterSupport, error) {
	is := &IgnoreSpecialModeSupport{
		modeMask: os.ModeSymlink | os.ModeNamedPipe | os.ModeSocket | os.ModeDevice | os.ModeIrregular,
	}
	is.SetName("IgnoreUnregularSupport")
	return is, nil
}

// IsIgnore ...
func (isms *IgnoreSpecialModeSupport) IsIgnore(path string, info os.FileInfo) (bool, error) {
	if isms.modeMask&info.Mode() != 0 {
		return true, nil
	}

	return false, nil
}

// IgnoreDotSupport ignore all dot fies.
type IgnoreDotSupport struct {
	BaseSupport
}

func NewFilterIgnoreDotSupport() (FilterSupport, error) {
	is := &IgnoreDotSupport{}
	is.SetName("IgnoreDotSupport")
	return is, nil
}

// IsIgnore ...
func (ids *IgnoreDotSupport) IsIgnore(path string, info os.FileInfo) (bool, error) {
	if filepath.Base(path)[0] == '.' {
		return true, nil
	}

	return false, nil
}
