package fail

import _ "embed"

var (
	//go:embed sample-1.pm
	sample1 string
	//go:embed sample-2.pm
	sample2 string
	//go:embed sample-3.pm
	sample3 string
)
var Samples = map[string]string{
	"sample-1.pm": sample1,
	"sample-2.pm": sample2,
	"sample-3.pm": sample3,
}
