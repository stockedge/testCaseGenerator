package tcg

import (
	parsec "github.com/prataprc/goparsec"
	"testing"
)

func TestConstraints(t *testing.T) {
	node, _ := constraints(parsec.NewScanner([]byte(`[ foo ] = 123;`)))
	if node == nil {
		t.Fail()
	}
}

func TestConstraint(t *testing.T) {
	node, _ := constraint(parsec.NewScanner([]byte(`[ foo ] = 123;`)))
	if node == nil {
		t.Fail()
	}
}

func TestPredicate(t *testing.T) {
	node, _ := predicate(parsec.NewScanner([]byte(`[ foo ] = 123`)))
	if node == nil {
		t.Fail()
	}
}

func TestClause(t *testing.T) {
	node, _ := clause(parsec.NewScanner([]byte(`[ foo ] = 123`)))
	if node == nil {
		t.Fail()
	}
}

func TestLogicalOperator(t *testing.T) {
	node, _ := logicalOperator(parsec.NewScanner([]byte(`and`)))
	if node == nil {
		t.Fail()
	}
}

func TestTerm(t *testing.T) {
	node, _ := term(parsec.NewScanner([]byte(`[ foo ] = 123`)))
	if node == nil {
		t.Fail()
	}
}

func TestParameterName(t *testing.T) {
	node, _ := parameterName(parsec.NewScanner([]byte(`[ foo ]`)))
	if node == nil {
		t.Fail()
	}
}

func TestRelation(t *testing.T) {
	node, _ := relation(parsec.NewScanner([]byte(`>=`)))
	if node == nil {
		t.Fail()
	}
}

func TestValue(t *testing.T) {
	node, _ := value(parsec.NewScanner([]byte(`123`)))
	if node == nil {
		t.Fail()
	}
	node, _ = value(parsec.NewScanner([]byte(`"foo"`)))
	if node == nil {
		t.Fail()
	}
}

func TestPatternString(t *testing.T) {
	//TODO
}

func TestValueSet(t *testing.T) {
	node, _ := valueSet(parsec.NewScanner([]byte(`"foo", 123`)))
	if node == nil {
		t.Fail()
	}
}
