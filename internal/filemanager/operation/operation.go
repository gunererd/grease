package operation

import "github.com/gunererd/grease/internal/filemanager/types"

type operation struct {
	opType types.OperationType
	source string
	target string
}

func New(opType types.OperationType, source string, target string) types.Operation {
	return &operation{
		opType: opType,
		source: source,
		target: target,
	}
}

func (o *operation) Type() types.OperationType {
	return o.opType
}

func (o *operation) Source() string {
	return o.source
}

func (o *operation) Target() string {
	return o.target
}
