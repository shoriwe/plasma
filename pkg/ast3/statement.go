package ast3

type (
	Statement interface {
		S3()
	}
	Assignment struct {
		Statement
		Left  Assignable
		Right Expression
	}
	Label struct {
		Statement
		Code int
	}
	Jump struct {
		Statement
		Target *Label
	}
	IfJump struct {
		Statement
		Condition Expression
		Target    *Label
	}
	Function struct {
		Statement
		Name      *Identifier
		Arguments []*Identifier
		Body      []Node
	}
	Return struct {
		Statement
		Result Expression
	}
)
