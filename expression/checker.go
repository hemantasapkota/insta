package expression

// Expression Checker ...
type ExpChecker struct {
	Node *Node
}

//IsLoop ...
func (e *ExpChecker) IsLoop() bool {
	if e.Node == nil {
		return false
	}
	src := e.Node.Src
	length := len(src)
	return length > 0 && string(src[:4]) == "loop"
}

//IsPool ...
func (e *ExpChecker) IsPool() bool {
	if e.Node.Src == nil {
		return false
	}
	length := len(e.Node.Src)
	return length > 0 && string(e.Node.Src) == "pool"
}
