package tcg

import (
	parsec "github.com/prataprc/goparsec"
)

type ConstraintExpr interface{}

type Constraints []Constraint

type Constraint parsec.Parser

type Predicate parsec.Parser

type Clause parsec.Parser

type LogicalOperator parsec.Parser

type Term parsec.Parser

type ParameterName parsec.Parser

type Relation parsec.Parser

type Value parsec.Parser

type PatternString parsec.Parser

type ValueSet parsec.Parser

type Env = interface{}

func eval(env Env) bool {
	return false
}

var constraints parsec.Parser

var constraint parsec.Parser

var predicate parsec.Parser

var clause parsec.Parser

var logicalOperator parsec.Parser

var term parsec.Parser

var parameterName parsec.Parser

var relation parsec.Parser

var value parsec.Parser

var patternString parsec.Parser

var valueSet parsec.Parser

func init() {
	//Constraints :=
	//Constraint
	//| Constraint Constraints
	constraints = parsec.Kleene(nil, &constraint)

	//Constraint    :: =
	//IF Predicate THEN Predicate ELSE Predicate;
	//| Predicate;
	constraint = parsec.OrdChoice(nil,
		parsec.And(nil, parsec.Atom("if", "IF"), &predicate, parsec.Atom("then", "THEN"),
			&predicate, parsec.Atom("else", "ELSE"), &predicate, parsec.Atom(";", ";")),
		parsec.And(nil, &predicate, parsec.Atom(";", ";")))

	//Predicate     :: =
	//Clause
	//| Clause LogicalOperator Predicate
	predicate = parsec.OrdChoice(nil, &clause, parsec.And(nil, &logicalOperator, &predicate))

	//Clause        :: =
	//Term
	//| ( Predicate )
	//| NOT Predicate
	clause = parsec.OrdChoice(nil,
		&term,
		parsec.And(nil, parsec.Atom("(", "("), &predicate, parsec.Atom(")", ")")),
		parsec.And(nil, parsec.Atom("not", "NOT"), &predicate))

	//Term          :: =
	//ParameterName Relation Value
	//| ParameterName LIKE PatternString
	//| ParameterName IN { ValueSet }
	//| ParameterName Relation ParameterName
	term = parsec.OrdChoice(nil,
		parsec.And(nil, &parameterName, &relation, &value),
		parsec.And(nil, &parameterName, parsec.Atom("like", "LIKE"), &patternString),
		parsec.And(nil, &parameterName, parsec.Atom("in", "IN"), parsec.Atom("{", "{"), &valueSet, parsec.Atom("}", "}")),
		parsec.And(nil, &parameterName, &relation, &parameterName))

	//ValueSet       :: =
	//Value
	//| Value, ValueSet
	valueSet = parsec.Kleene(nil, &value, parsec.Atom(",", ","))

	//LogicalOperator ::=
	//AND
	//| OR
	logicalOperator = parsec.OrdChoice(nil, parsec.Atom("and", "AND"), parsec.Atom("OR", "OR"))

	//Relation      :: =
	//=
	//| <>
	//| >
	//| >=
	//| <
	//| <=
	relation = parsec.OrdChoice(nil,
		parsec.Atom("=", "="),
		parsec.Atom("<>", "<>"),
		parsec.Atom(">", ">"),
		parsec.Atom(">=", ">="),
		parsec.Atom("<", "<"),
		parsec.Atom("<=", "<="),
	)

	//ParameterName ::= [String]
	parameterName = parsec.And(nil,
		parsec.Atom("[", "["),
		parsec.Token("[^\\]]+", "STRING"),
		parsec.Atom("]", "]"),
	)

	//Value         :: =
	//"String"
	//| Number
	value = parsec.OrdChoice(nil, parsec.String(), parsec.Int())

}
