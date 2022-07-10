package transformations_1

import "github.com/shoriwe/gplasma/pkg/ast3"

func (transform *transformPass) nextLabel() *ast3.Label {
	label := transform.currentLabel
	transform.currentLabel++
	return &ast3.Label{
		Code:      label,
	}
}
