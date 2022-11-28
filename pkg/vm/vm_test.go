package vm

import (
	"bytes"
	"fmt"
	"github.com/shoriwe/plasma/pkg/test-samples/fail"
	"github.com/shoriwe/plasma/pkg/test-samples/success"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSuccessSampleScripts(t *testing.T) {
	for index := 1; index <= len(success.Samples); index++ {
		func(i int) {
			sampleScript := fmt.Sprintf("sample-%d.pm", i)
			script := success.Samples[sampleScript]
			out := &bytes.Buffer{}
			v := NewVM(nil, out, nil)
			rCh, errCh, _ := v.ExecuteString(script.Code)
			defer close(errCh)
			defer close(rCh)
			assert.Nil(t, <-errCh)
			<-rCh
			s := out.String()
			assert.Equal(t, script.Result, out.String(), fmt.Sprintf("Expecting:\n%s\nBut received:\n%s", script.Result, s))
		}(index)
	}
}

func TestFailSampleScripts(t *testing.T) {
	for index := 1; index <= len(fail.Samples); index++ {
		func(i int) {
			sampleScript := fmt.Sprintf("sample-%d.pm", i)
			script := fail.Samples[sampleScript]
			v := NewVM(nil, nil, nil)
			rCh, errCh, _ := v.ExecuteString(script)
			defer close(errCh)
			defer close(rCh)
			assert.NotNil(t, <-errCh)
		}(index)
	}
}
