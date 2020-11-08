package starlark

import "go.starlark.net/syntax"

type NodeAction struct {
	Recurse   bool
	Replace   bool
	ReplaceBy syntax.Node
}

// Rewrite traverses a syntax tree in depth-first order.
// It starts by calling f(n); n must not be nil.
// f then returns whether to Rewrite itself
// recursively or to replace n.
// Rewrite then calls f(nil).
func Rewrite(n syntax.Node, f func(syntax.Node) NodeAction) syntax.Node {
	if n == nil {
		panic("nil")
	}

	nr := f(n)
	if nr.Replace {
		return nr.ReplaceBy
	}
	if !nr.Recurse {
		return n
	}

	// TODO(adonovan): opt: order cases using profile data.
	switch n := n.(type) {
	case *syntax.File:
		replaceStmts(n.Stmts, f)

	case *syntax.ExprStmt:
		replace := rewriteExpr(n.X, f)
		if replace != nil {
			n.X = replace
		}

	case *syntax.BranchStmt:
		// no-op

	case *syntax.IfStmt:
		replace := rewriteExpr(n.Cond, f)
		if replace != nil {
			n.Cond = replace
		}
		replaceStmts(n.True, f)
		replaceStmts(n.False, f)

	case *syntax.AssignStmt:
		n.LHS = rewriteExpr(n.LHS, f)
		n.RHS = rewriteExpr(n.RHS, f)

	case *syntax.DefStmt:
		n.Name = rewriteIdent(n.Name, f)
		for i, param := range n.Params {
			n.Params[i] = rewriteExpr(param, f)

		}
		replaceStmts(n.Body, f)

	case *syntax.ForStmt:
		n.Vars = rewriteExpr(n.Vars, f)
		n.X = rewriteExpr(n.X, f)
		replaceStmts(n.Body, f)

	case *syntax.ReturnStmt:
		if n.Result != nil {
			n.Result = rewriteExpr(n.Result, f)
		}

	case *syntax.LoadStmt:
		n.Module = rewriteLiteral(n.Module, f)
		for i, from := range n.From {
			n.From[i] = rewriteIdent(from, f)
		}
		for i, to := range n.To {
			n.To[i] = rewriteIdent(to, f)
		}

	case *syntax.Ident, *syntax.Literal:
		// no-op

	case *syntax.ListExpr:
		for i, x := range n.List {
			n.List[i] = rewriteExpr(x, f)
		}

	case *syntax.ParenExpr:
		n.X = rewriteExpr(n.X, f)

	case *syntax.CondExpr:
		n.Cond = rewriteExpr(n.Cond, f)
		n.True = rewriteExpr(n.True, f)
		n.False = rewriteExpr(n.False, f)

	case *syntax.IndexExpr:
		n.X = rewriteExpr(n.X, f)
		n.Y = rewriteExpr(n.Y, f)

	case *syntax.DictEntry:
		n.Key = rewriteExpr(n.Key, f)
		n.Value = rewriteExpr(n.Value, f)

	case *syntax.SliceExpr:
		n.X = rewriteExpr(n.X, f)
		if n.Lo != nil {
			n.Lo = rewriteExpr(n.Lo, f)
		}
		if n.Hi != nil {
			n.Hi = rewriteExpr(n.Hi, f)
		}
		if n.Step != nil {
			n.Step = rewriteExpr(n.Step, f)
		}

	case *syntax.Comprehension:
		n.Body = rewriteExpr(n.Body, f)
		for i, clause := range n.Clauses {
			n.Clauses[i] = Rewrite(clause, f)
		}

	case *syntax.IfClause:
		n.Cond = rewriteExpr(n.Cond, f)

	case *syntax.ForClause:
		n.Vars = rewriteExpr(n.Vars, f)
		n.X = rewriteExpr(n.X, f)

	case *syntax.TupleExpr:
		for i, x := range n.List {
			n.List[i] = rewriteExpr(x, f)
		}

	case *syntax.DictExpr:
		for _, entry := range n.List {
			entry := entry.(*syntax.DictEntry)
			entry.Key = rewriteExpr(entry.Key, f)
			entry.Value = rewriteExpr(entry.Value, f)
		}

	case *syntax.UnaryExpr:
		if n.X != nil {
			n.X = rewriteExpr(n.X, f)
		}

	case *syntax.BinaryExpr:
		n.X = rewriteExpr(n.X, f)
		n.Y = rewriteExpr(n.Y, f)

	case *syntax.DotExpr:
		n.X = rewriteExpr(n.X, f)
		n.Name = rewriteIdent(n.Name, f)

	case *syntax.CallExpr:
		n.Fn = rewriteExpr(n.Fn, f)
		for i, arg := range n.Args {
			n.Args[i] = rewriteExpr(arg, f)
		}

	case *syntax.LambdaExpr:
		for i, param := range n.Params {
			n.Params[i] = rewriteExpr(param, f)
		}
		n.Body = rewriteExpr(n.Body, f)

	default:
		panic(n)
	}

	f(nil)
	return n
}

func rewriteExpr(n syntax.Expr, f func(syntax.Node) NodeAction) syntax.Expr {
	return Rewrite(n, f).(syntax.Expr)
}

func rewriteIdent(n *syntax.Ident, f func(syntax.Node) NodeAction) *syntax.Ident {
	return Rewrite(n, f).(*syntax.Ident)
}

func rewriteLiteral(n *syntax.Literal, f func(syntax.Node) NodeAction) *syntax.Literal {
	return Rewrite(n, f).(*syntax.Literal)
}

func rewriteStmt(n syntax.Stmt, f func(syntax.Node) NodeAction) syntax.Stmt {
	return Rewrite(n, f).(syntax.Stmt)
}

func replaceStmts(stmts []syntax.Stmt, f func(syntax.Node) NodeAction) {
	for i, stmt := range stmts {
		stmts[i] = rewriteStmt(stmt, f)
	}
}
