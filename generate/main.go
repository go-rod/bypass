package main

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ysmood/kit"
)

var slash = filepath.FromSlash

func main() {

	build := kit.S(`// generated by running "go generate" on project root

package bypass

// JS for bypass
const JS = {{.js}}
`,
		"js", encode(fetchBypassJS()),
	)

	kit.E(kit.OutputFile(slash("assets.go"), build, nil))
}

func fetchBypassJS() string {
	kit.Exec("npx", "extract-stealth-evasions").MustDo()

	code, err := kit.ReadString("stealth.min.js")
	kit.E(err)

	// since the npx already mentioned extract-stealth-evasions, we don't have to do it again.
	code = regexp.MustCompile(`\A/\*\![\s\S]+?\*/`).ReplaceAllString(code, "")

	return fmt.Sprintf(";(() => {\n%s\n})();", code)
}

// not using encoding like base64 or gzip because of they will make git diff every large for small change
func encode(s string) string {
	return "`" + strings.ReplaceAll(s, "`", "` + \"`\" + `") + "`"
}
