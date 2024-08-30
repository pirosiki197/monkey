// Code generated by "stringer -type ObjectType -linecomment object.go"; DO NOT EDIT.

package object

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[INTEGER_OBJ-1]
	_ = x[BOOLEAN_OBJ-2]
	_ = x[NULL_OBJ-3]
	_ = x[RETURN_VALUE_OBJ-4]
}

const _ObjectType_name = "INTEGERBOOLEANNULLRETURN_VALUE"

var _ObjectType_index = [...]uint8{0, 7, 14, 18, 30}

func (i ObjectType) String() string {
	i -= 1
	if i < 0 || i >= ObjectType(len(_ObjectType_index)-1) {
		return "ObjectType(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _ObjectType_name[_ObjectType_index[i]:_ObjectType_index[i+1]]
}
