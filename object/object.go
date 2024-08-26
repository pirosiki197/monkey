package object

import "fmt"

type ObjectType int

type Object interface {
	Type() ObjectType
	Inspect() string
}

//go:generate stringer -type ObjectType -linecomment object.go
const (
	_           ObjectType = iota
	INTEGER_OBJ            // INTEGER
	BOOLEAN_OBJ            // BOOLEAN
	NULL_OBJ               // NULL
)

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }
