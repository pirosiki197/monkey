// Code generated by "stringer -type ObjectType -linecomment object.go"; DO NOT EDIT.

package object

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[INTEGER_OBJ-1]
	_ = x[STRING_OBJ-2]
	_ = x[BOOLEAN_OBJ-3]
	_ = x[NULL_OBJ-4]
	_ = x[RETURN_VALUE_OBJ-5]
	_ = x[FUNCTION_OBJ-6]
	_ = x[ERROR_OBJ-7]
	_ = x[BUILTIN_OBJ-8]
}

const _ObjectType_name = "INTEGERSTRINGBOOLEANNULLRETURN_VALUEFUNCTIONERRORBUILTIN"

var _ObjectType_index = [...]uint8{0, 7, 13, 20, 24, 36, 44, 49, 56}

func (i ObjectType) String() string {
	i -= 1
	if i < 0 || i >= ObjectType(len(_ObjectType_index)-1) {
		return "ObjectType(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _ObjectType_name[_ObjectType_index[i]:_ObjectType_index[i+1]]
}
