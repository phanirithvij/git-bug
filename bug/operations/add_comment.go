package operations

import (
	"github.com/MichaelMure/git-bug/bug"
	"time"
)

var _ bug.Operation = AddCommentOperation{}

type AddCommentOperation struct {
	bug.OpBase
	Message string
	Author  bug.Person
	Time    int64
}

func NewAddCommentOp(author bug.Person, message string) AddCommentOperation {
	return AddCommentOperation{
		OpBase:  bug.OpBase{OperationType: bug.AddCommentOp},
		Message: message,
		Author:  author,
		Time:    time.Now().Unix(),
	}
}

func (op AddCommentOperation) Apply(snapshot bug.Snapshot) bug.Snapshot {
	comment := bug.Comment{
		Message: op.Message,
		Author:  op.Author,
		Time:    op.Time,
	}

	snapshot.Comments = append(snapshot.Comments, comment)

	return snapshot
}