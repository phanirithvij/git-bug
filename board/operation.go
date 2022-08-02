package board

import (
	"encoding/json"
	"fmt"

	"github.com/MichaelMure/git-bug/entity/dag"
	"github.com/MichaelMure/git-bug/identity"
)

// OperationType is an operation type identifier
type OperationType dag.OperationType

const (
	_ dag.OperationType = iota
	CreateOp
	SetTitleOp
	SetDescriptionOp
	AddItemOp
	MoveItemOp
	RemoveItemOp

	// TODO: change columns?
)

type Operation interface {
	dag.Operation
	// Apply the operation to a Snapshot to create the final state
	Apply(snapshot *Snapshot)
}

func operationUnmarshaller(raw json.RawMessage, resolver identity.Resolver) (dag.Operation, error) {
	var t struct {
		OperationType dag.OperationType `json:"type"`
	}

	if err := json.Unmarshal(raw, &t); err != nil {
		return nil, err
	}

	var op dag.Operation

	switch t.OperationType {
	case CreateOp:
		op = &CreateOperation{}
	case SetTitleOp:
		op = &SetTitleOperation{}
	case SetDescriptionOp:
		op = &SetDescriptionOperation{}
	default:
		panic(fmt.Sprintf("unknown operation type %v", t.OperationType))
	}

	err := json.Unmarshal(raw, &op)
	if err != nil {
		return nil, err
	}

	return op, nil
}
