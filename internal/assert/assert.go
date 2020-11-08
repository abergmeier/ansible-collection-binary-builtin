package assert

import (
	"fmt"
	"strings"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"go.starlark.net/syntax"

	"github.com/abergmeier/ansible-collection-binary-builtin/internal/ansible"
	starlarkrw "github.com/abergmeier/ansible-collection-binary-builtin/internal/starlark"
)

// Parameters to feed to this module
type Parameters struct {
	FailMsg    string
	SuccessMsg string
	Quiet      bool
	That       []string
}

// Run assert validation
func Run(p *Parameters) (bool, error) {

	thread := &starlark.Thread{
		Name: "Assert",
	}

	for i, t := range p.That {
		truth, err := evalExprTruth(thread, i, t)
		if err != nil {
			return false, err
		}
		if !truth {
			return false, nil
		}
	}

	return true, nil
}

func evalExprTruth(t *starlark.Thread, i int, in string) (bool, error) {

	clean := strings.ReplaceAll(in, " is ", " == ")

	expr, err := syntax.ParseExpr(fmt.Sprintf("assert.%d.star", i), clean, syntax.RetainComments)
	if err != nil {
		return false, err
	}

	expr = starlarkrw.Rewrite(expr, ansiblePythonToStarklark).(syntax.Expr)

	env := starlark.StringDict{
		"ansible_version": ansibleVersion(),
		"directory":       starlark.NewBuiltin("directory", ansible.Directory),
		"exists":          starlark.NewBuiltin("exists", ansible.Exists),
		"file":            starlark.NewBuiltin("file", ansible.File),
		"version":         starlark.NewBuiltin("version", ansible.Version),
		"version_compare": starlark.NewBuiltin("version_compare", ansible.Version),
	}
	val, err := starlark.EvalExpr(t, expr, env)
	if err != nil {
		return false, err
	}

	if !val.Truth() {
		return false, nil
	}

	return true, nil
}

func ansibleVersion() starlark.Value {

	d := starlark.StringDict{
		"full": starlark.String("2.12.1"),
	}
	return starlarkstruct.FromStringDict(starlarkstruct.Default, d)
}

func ansiblePythonToStarklark(n syntax.Node) starlarkrw.NodeAction {
	switch et := n.(type) {
	case *syntax.BinaryExpr:

		callExpr := getPythonIsWithCall(et)
		if callExpr == nil {
			break
		}

		callExpr.Args = append([]syntax.Expr{et.X}, callExpr.Args...)

		return starlarkrw.NodeAction{
			Replace:   true,
			ReplaceBy: callExpr,
		}
	}

	return starlarkrw.NodeAction{
		Recurse: true,
	}
}

func getPythonIsWithCall(b *syntax.BinaryExpr) *syntax.CallExpr {
	// We implicitly know that is has been replaced with == before
	if b.Op != syntax.EQL {
		return nil
	}

	call, ok := b.Y.(*syntax.CallExpr)
	if ok {
		return pythonCall(call)
	}
	ident, ok := b.Y.(*syntax.Ident)
	if ok {
		return pythonIdent(ident)
	}

	return nil
}

func pythonIdent(ident *syntax.Ident) *syntax.CallExpr {
	if ident.Name != "directory" && ident.Name != "exists" && ident.Name != "file" && ident.Name != "link" {
		return nil
	}

	return &syntax.CallExpr{
		Fn: &syntax.Ident{
			Name: ident.Name,
		},
	}
}

func pythonCall(call *syntax.CallExpr) *syntax.CallExpr {

	ident, ok := call.Fn.(*syntax.Ident)
	if !ok {
		return nil
	}

	if ident.Name != "version" && ident.Name != "version_compare" {
		return nil
	}

	return call
}
