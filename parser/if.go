package parser

type IfBranch struct {
	condition Expression
	body      []Statement
}

func (b IfBranch) Condition() Expression {
	return b.condition
}

func (b IfBranch) Body() []Statement {
	return b.body
}

type Else struct {
	body []Statement
}

func (e Else) Body() []Statement {
	return e.body
}

type If struct {
	ifBranch     IfBranch
	elifBranches []IfBranch
	elseBranch   Else
}

func (i If) StatementType() StatementType {
	return STATEMENT_TYPE_IF
}

func (i If) IfBranch() IfBranch {
	return i.ifBranch
}

func (i If) ElseIfBranches() []IfBranch {
	return i.elifBranches
}

func (i If) Else() Else {
	return i.elseBranch
}

func (i If) HasElse() bool {
	return len(i.Else().Body()) > 0
}
