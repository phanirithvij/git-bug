package board

import (
	"fmt"

	"github.com/MichaelMure/git-bug/entity"
	"github.com/MichaelMure/git-bug/entity/dag"
	"github.com/MichaelMure/git-bug/identity"
	"github.com/MichaelMure/git-bug/util/text"
)

var DefaultColumns = []string{"To Do", "In Progress", "Done"}

var _ dag.Operation = &CreateOperation{}

type CreateOperation struct {
	dag.OpBase
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Columns     []string `json:"columns"`
}

func NewCreateOp(author identity.Interface, unixTime int64, title string, description string, columns []string) *CreateOperation {
	return &CreateOperation{
		OpBase:      dag.NewOpBase(CreateOp, author, unixTime),
		Title:       title,
		Description: description,
		Columns:     columns,
	}
}

func (op *CreateOperation) Id() entity.Id {
	return dag.IdOperation(op, &op.OpBase)
}

func (op *CreateOperation) Validate() error {
	if err := op.OpBase.Validate(op, CreateOp); err != nil {
		return err
	}

	if text.Empty(op.Title) {
		return fmt.Errorf("title is empty")
	}
	if !text.SafeOneLine(op.Title) {
		return fmt.Errorf("title has unsafe characters")
	}

	if !text.SafeOneLine(op.Description) {
		return fmt.Errorf("description has unsafe characters")
	}

	if len(op.Columns) <= 0 {
		return fmt.Errorf("no columns")
	}
	for _, column := range op.Columns {
		if !text.SafeOneLine(column) {
			return fmt.Errorf("a columns has unsafe characters")
		}
		if len(column) > 100 {
			return fmt.Errorf("a columns is too long")
		}
	}

	return nil
}

func (op *CreateOperation) Apply(snap *Snapshot) {
	// sanity check: will fail when adding a second Create
	if snap.id != "" && snap.id != entity.UnsetId && snap.id != op.Id() {
		return
	}

	snap.id = op.Id()

	snap.Title = op.Title
	snap.Description = op.Description

	for _, name := range op.Columns {
		snap.Columns = append(snap.Columns, Column{
			Name:  name,
			Cards: nil,
		})
	}

	snap.addActor(op.Author())
}

func CreateDefaultColumns(author identity.Interface, unixTime int64, title, description string) (*Board, *CreateOperation, error) {
	return Create(author, unixTime, title, description, DefaultColumns)
}

func Create(author identity.Interface, unixTime int64, title, description string, columns []string) (*Board, *CreateOperation, error) {
	b := NewBoard()
	op := NewCreateOp(author, unixTime, title, description, columns)
	if err := op.Validate(); err != nil {
		return nil, op, err
	}
	b.Append(op)
	return b, op, nil
}
