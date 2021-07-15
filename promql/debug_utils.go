package promql

import (
	"bufio"
	"bytes"
	"regexp"
	"runtime/debug"
	"strings"
)

var patterns = []struct {
	re      *regexp.Regexp
	extract string
}{
	{regexp.MustCompile(`^\s+\/home.+eus\/prometheus\/(.*)`), "\t> prom:  ${1}"},
	{regexp.MustCompile(`^\s+\/home.+\/mod\/.+\/prometheus/(.*)`), "\t> mod :  ${1}"},
	{regexp.MustCompile(`^\s+\/home.+go\/src\/(.*)`), "\t> go  : ${1}"},
	{regexp.MustCompile(`^\s+\/home.+\/mod\/.+\/(.+)@v.{5,8}/(.*)`), "\t> mod :  ${1}/${2}"},
	{regexp.MustCompile(`\(0x.*\)`), ""},
	{regexp.MustCompile(`\s\+0x.{2,8}$`), ""},
	{regexp.MustCompile(`github.com\/.*\/prometheus\/(.*)`), "pkg: $1"},
}

func stacktrace(depths ...int) string {

	depth := 40
	if len(depths) > 0 {
		depth = depths[0]
	}

	br := bytes.NewReader(debug.Stack())
	buf := bufio.NewReader(br)

	for skip := 0; skip <= 4; skip++ {
		buf.ReadLine()
	}

	var stack strings.Builder
	newLine := true
	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			break
		}

		str := string(line)
		for _, p := range patterns {
			str = p.re.ReplaceAllString(str, p.extract)
		}
		stack.WriteString(str)
		newLine = !newLine
		if !newLine {
			continue
		}
		stack.WriteString("\n")
		depth--
		if depth == 0 {
			break
		}
	}

	return stack.String()
}
