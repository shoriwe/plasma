//go:build ignore

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
)

var (
	number = regexp.MustCompile("\\d+")
	chars  = regexp.MustCompile("\\w+")
)

func removeCRLF(file string) {
	contents, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(file, bytes.ReplaceAll(contents, []byte("\r"), []byte{}), 0777)
	if err != nil {
		panic(err)
	}
}

func dir(name string) {
	file, err := os.Create(filepath.Join(name, "samples.go"))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	fmt.Fprintf(file, "package %s\n", name)
	fmt.Fprintln(file, `import _ "embed"`)
	fmt.Fprintln(file, "var(")
	files, err := os.ReadDir(name)
	if err != nil {
		panic(err)
	}
	for _, entry := range files {
		if entry.Name() == "samples.go" {
			continue
		}
		removeCRLF(filepath.Join(name, entry.Name()))
		fmt.Fprintf(file, "//go:embed %s\n", entry.Name())
		fmt.Fprintf(file, "sample%s string\n", number.FindString(entry.Name()))
	}
	fmt.Fprintln(file, ")")
	fmt.Fprintln(file, "var Samples = map[string]string{")
	for _, entry := range files {
		if entry.Name() == "samples.go" {
			continue
		}
		fmt.Fprintf(file, "\"%s\": %s,\n", entry.Name(), fmt.Sprintf("sample%s", number.FindString(entry.Name())))
	}
	fmt.Fprintln(file, "}")
}

func compare(name string) {
	file, err := os.Create(filepath.Join(name, "samples.go"))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	fmt.Fprintf(file, "package %s\n", name)
	fmt.Fprintln(file, `
import (
	_ "embed"
	test_samples "github.com/shoriwe/plasma/pkg/test-samples"
)
`)
	fmt.Fprintln(file, "var(")
	files, err := os.ReadDir(name)
	if err != nil {
		panic(err)
	}
	max := 0
	for _, entry := range files {
		if entry.Name() == "samples.go" {
			continue
		}
		removeCRLF(filepath.Join(name, entry.Name()))
		n := number.FindString(entry.Name())
		if nn, _ := strconv.Atoi(n); nn > max {
			max = nn
		}
		fmt.Fprintf(file, "//go:embed %s\n", entry.Name())
		fmt.Fprintf(file, "%s%s string\n", chars.FindString(entry.Name()), n)
	}
	fmt.Fprintln(file, ")")
	fmt.Fprintln(file, "var Samples = map[string]test_samples.Script{")
	for i := 0; i < max; i++ {
		fmt.Fprintf(file, `
"sample-%d.pm": {
	Code: sample%d,
	Result: result%d,
},
`, i+1, i+1, i+1)
	}
	fmt.Fprintln(file, "}")
}

func main() {
	dir("basic")
	dir("fail")
	compare("success")
	exec.Command("go", "fmt", "./...").Run()
}
