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
package parser

import (
	"path/filepath"
	re "regexp"
)

var (
	reHeader = re.MustCompile(`^(?P<Level>\*+)\s`)
	reLink   = re.MustCompile(`\[\[(?P<Type>\w+):(?P<Path>.+?)\]\]`)
)

type OrgHeader struct {
	Stars string
	Level int
	Text  string
}

// OrgLink is a data struct describing a link with the file and
// headers contain it
type OrgLink struct {
	File    string
	Headers []OrgHeader
	Link    string

	Type string
	Path string
}

// OrgLinkParser implements FileLineParser, is used to parse org file
// to get OrgLink.
type OrgLinkParser struct {
	file       string
	curHeaders []OrgHeader
}

// NewOrgLinkParser create a new OrgLinkParser
func NewOrgLinkParser(file string) *OrgLinkParser {
	abspath, err := filepath.Abs(file)
	if err != nil {
		panic(err)
	}

	return &OrgLinkParser{
		file:       abspath,
		curHeaders: make([]OrgHeader, 0, 4),
	}
}

// cloneHeader ...
func (fp *OrgLinkParser) cloneHeader() []OrgHeader {
	orgHeader := make([]OrgHeader, len(fp.curHeaders))
	copy(orgHeader, fp.curHeaders)
	return orgHeader
}

// pushHeader ...
func (fp *OrgLinkParser) pushHeader(header OrgHeader) {
	fp.curHeaders = append(fp.curHeaders, header)
}

// popHeader ...
func (fp *OrgLinkParser) popHeader() (header OrgHeader) {
	if len(fp.curHeaders) == 0 {
		return OrgHeader{}
	}

	length := len(fp.curHeaders)
	v := fp.curHeaders[length-1]

	fp.curHeaders = fp.curHeaders[:length-1]
	return v
}

// getHeader ...
func (fp *OrgLinkParser) getHeader() (header OrgHeader) {
	if len(fp.curHeaders) == 0 {
		return OrgHeader{}
	}
	return fp.curHeaders[len(fp.curHeaders)-1]
}

// FilePath return the path of file to be parsed
func (fp *OrgLinkParser) FilePath() string {
	return fp.file
}

// Parse parse org file, if line contain a link of file, Parse return a
// OrgLink with file name and headers above link.
func (fp *OrgLinkParser) Parse(line string) (interface{}, error) {
	// first match stars at the beginning of line
	if reHeader.MatchString(line) {
		stars := reHeader.ExpandString([]byte{}, "$Level", line,
			reHeader.FindStringSubmatchIndex(line))
		curHeader := OrgHeader{
			Stars: string(stars),
			Level: len(stars),
			Text:  line[len(stars)+1:],
		}

		latest := fp.getHeader()
		for latest.Level >= curHeader.Level {
			fp.popHeader()
			latest = fp.getHeader()
		}

		fp.pushHeader(curHeader)

		return nil, nil
	}

	// second match link element
	if reLink.MatchString(line) {
		// TODO find all links in single line
		link := reLink.FindString(line)
		indexes := reLink.FindStringSubmatchIndex(line)
		typ := reLink.ExpandString([]byte{}, "$Type", line, indexes)
		path := string(reLink.ExpandString(
			[]byte{}, "$Path", line, indexes))
		path = filepath.Join(filepath.Dir(fp.file), path)

		return &OrgLink{
			File:    fp.FilePath(),
			Headers: fp.cloneHeader(),
			Link:    link,
			Type:    string(typ),
			Path:    path,
		}, nil
	}

	return nil, nil

}
