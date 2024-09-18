package evaluator

import (
	"fmt"
	"unique"

	"github.com/pirosiki197/monkey/ast"
	"github.com/pirosiki197/monkey/object"
)

type Evaluator struct {
	env *object.Environment
}

func New() *Evaluator {
	return &Evaluator{
		env: object.NewEnvironment(),
	}
}

func NewWithEnv(env *object.Environment) *Evaluator {
	return &Evaluator{
		env: env,
	}
}

func (e *Evaluator) Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return e.evalProgram(node)
	case *ast.BlockStatement:
		return e.evalBlockStatements(node)
	case *ast.LetStatement:
		val := e.Eval(node.Value)
		if isError(val) {
			return val
		}
		e.env.Set(node.Name.Value, val)
		return nil
	case *ast.AssignStatement:
		val := e.Eval(node.Value)
		if isError(val) {
			return val
		}
		_, ok := e.env.Update(node.Name.Value, val)
		if !ok {
			return newError("identifier not found: %s", node.Name.Value)
		}
		return nil
	case *ast.ReturnStatement:
		val := e.Eval(node.ReturnValue)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.IfExpression:
		return e.evalIfExpression(node)
	case *ast.ExpressionStatement:
		return e.Eval(node.Expression)
	case *ast.PrefixExpression:
		right := e.Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := e.Eval(node.Left)
		if isError(left) {
			return left
		}
		right := e.Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.CallExpression:
		function := e.Eval(node.Function)
		if isError(function) {
			return function
		}
		args := e.evalExpressions(node.Arguments)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		evaluated := applyFunction(function, args)
		return unwrapReturnValue(evaluated)
	case *ast.Identifier:
		return e.evalIdentifier(node)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: unique.Make(node.Value)}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: e.env, Body: body}
	default:
		return nil
	}
}

func (e *Evaluator) evalProgram(program *ast.Program) object.Object {
	stmts := program.Statements
	var result object.Object
	for _, stmt := range stmts {
		result = e.Eval(stmt)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

func (e *Evaluator) evalBlockStatements(block *ast.BlockStatement) object.Object {
	stmts := block.Statements
	var result object.Object
	enclosedEvaluator := &Evaluator{env: object.NewEnclosedEnvironment(e.env)}
	for _, stmt := range stmts {
		result = enclosedEvaluator.Eval(stmt)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}
	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		if len(fn.Parameters) != len(args) {
			return newError("wrong length of arguments: %d parameters but called with %d arguments",
				len(fn.Parameters),
				len(args))
		}
		extendedEnv := extendFunctionEnv(fn, args)
		e := NewWithEnv(extendedEnv)
		evaluated := e.Eval(fn.Body)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

func (e *Evaluator) evalExpressions(exps []ast.Expression) []object.Object {
	result := make([]object.Object, 0, len(exps))

	for _, exp := range exps {
		evaluated := e.Eval(exp)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func (e *Evaluator) evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := e.Eval(ie.Condition)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return e.Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return e.Eval(ie.Alternative)
	} else {
		return NULL
	}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case both(left, right, object.INTEGER_OBJ):
		return evalIntegerInfixExpression(operator, left, right)
	case both(left, right, object.STRING_OBJ):
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: unique.Make(leftVal.Value() + rightVal.Value())}
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
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

func evalMinusOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func (e *Evaluator) evalIdentifier(node *ast.Identifier) object.Object {
	if val, ok := e.env.Get(node.Value); ok {
		return val
	}
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	return newError("identifier not found: %s", node.Value)
}

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func nativeBoolToBooleanObject(b bool) *object.Boolean {
	if b {
		return TRUE
	} else {
		return FALSE
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

func isError(obj object.Object) bool {
	return obj.Type() == object.ERROR_OBJ
}

func newError(format string, a ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func both(left, right object.Object, objType object.ObjectType) bool {
	return left.Type() == objType && right.Type() == objType
}
