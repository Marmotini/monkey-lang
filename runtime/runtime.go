package runtime

import (
	"fmt"

	"github.com/nerdysquirrel/monkey-lang/ast"
	"github.com/nerdysquirrel/monkey-lang/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(
	node ast.Node,
	env *object.Environment) object.Object {

	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolean(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)
	case *ast.BlockStatement:
		return evalBlockStatements(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.ReturnStatement:
		value := Eval(node.ReturnValue, env)
		if isError(value) {
			return value
		}
		return &object.ReturnValue{Value: value}
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}

		env.Set(node.Name.Value, val)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	}

	return nil
}

func evalProgram(
	prog *ast.Program,
	env *object.Environment) object.Object {

	var results object.Object

	for _, statement := range prog.Statements {
		results = Eval(statement, env)

		switch r := results.(type) {
		case *object.ReturnValue:
			return r.Value
		case *object.Error:
			return r
		}
	}

	return results
}

func evalBlockStatements(
	block *ast.BlockStatement,
	env *object.Environment) object.Object {

	var results object.Object

	for _, statement := range block.Statements {
		results = Eval(statement, env)

		if results != nil {
			r := results.Type()

			if r == object.RETURN_VALUE_OBJ || r == object.ERROR_OBJ {
				return results
			}
		}
	}

	return results
}

func nativeBoolean(input bool) object.Object {
	if input {
		return TRUE
	}

	return FALSE
}

func evalPrefixExpression(
	operator string,
	right object.Object) object.Object {

	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	}

	// TODO
	return newError("unknown operator: %s%s", operator, right.Type())
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(
	operator string,
	left, right object.Object) object.Object {

	if left.Type() != right.Type() {
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	}

	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case  operator == "==":
		return nativeBoolean(left == right)
	case operator == "!=" :
		return nativeBoolean(left != right)
	}

	return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())

}

func evalIntegerInfixExpression(
	operator string,
	left, right object.Object) object.Object {

	l := left.(*object.Integer).Value
	r := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: l + r}
	case "-":
		return &object.Integer{Value: l - r}
	case "/":
		return &object.Integer{Value: l / r}
	case "*":
		return &object.Integer{Value: l * r}

	case "<":
		return nativeBoolean(l < r)
	case ">":
		return nativeBoolean(l > r)
	case "==":
		return nativeBoolean(l == r)
	case "!=":
		return nativeBoolean(l != r)
	}

	return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
}

func evalIfExpression(
	node *ast.IfExpression,
	env *object.Environment) object.Object {

	condition := Eval(node.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(node.Consequence, env)
	} else if node.Alternative != nil {
		return Eval(node.Alternative, env)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func newError(format string, s ... interface{}) object.Object {
	return &object.Error{Message: fmt.Sprintf(format, s)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}

	return false
}

func evalIdentifier(
	node *ast.Identifier,
	env *object.Environment) object.Object {

	val, ok := env.Get(node.Value)
	if !ok {
		return newError("identifier not found: %s", node.Value)
	}

	return val

}