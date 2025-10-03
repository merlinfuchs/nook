package thing

type Operation string

var (
	OperationOverwrite Operation = "overwrite"
	OperationAppend    Operation = "append"
	OperationPrepend   Operation = "prepend"
	OperationIncrement Operation = "increment"
	OperationDecrement Operation = "decrement"
)
