package ansible

import (
	"fmt"
	"strings"

	version "github.com/hashicorp/go-version"
	"go.starlark.net/starlark"
)

const (
	eq = iota
	lt
	le
	gt
	ge
	ne
)

var foo = `
def version_compare(value, version, operator='eq', strict=None, version_type=None):

    type_map = {
        'loose': LooseVersion,
        'strict': StrictVersion,
        'semver': SemanticVersion,
        'semantic': SemanticVersion,
    }

    if strict != None and version_type != None:
        fail("Cannot specify both 'strict' and 'version_type'")

    Version = LooseVersion
    if strict:
        Version = StrictVersion
    elif version_type:
        Version = type_map.get(version_type, None)
        if Version == None:
            fail("Invalid version type (%s). Must be one of %s" % (version_type, ', '.join(map(repr, type_map))))

    if operator in op_map:
        operator = op_map[operator]
    else:
        fail('Invalid operator type (%s). Must be one of %s' % (operator, ', '.join(map(repr, op_map))))

    method = getattr(py_operator, operator)
    return method(Version(to_text(value)), Version(to_text(version)))
    #raise errors.AnsibleFilterError('Version comparison failed: %s' % to_native(e))
`

// Version performs a version comparison on a value
func Version(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {

	opMap := map[string]int{
		"==": eq, "=": eq, "eq": eq,
		"<": lt, "lt": lt,
		"<=": le, "le": le,
		">": gt, "gt": gt,
		">=": ge, "ge": ge,
		"!=": ne, "<>": ne, "ne": ne,
	}
	value, ok := starlark.AsString(args[0])
	if !ok {
		panic(ok)
	}
	ver, ok := starlark.AsString(args[1])
	if !ok {
		panic(ok)
	}

	op := eq
	if len(args) > 2 {
		operator := args[2]
		var ok bool
		os, ok := starlark.AsString(operator)
		if !ok {
			panic(ok)
		}
		op, ok = opMap[os]
		if !ok {
			opMapKeys := make([]string, 0, len(opMap))
			for k := range opMap {
				opMapKeys = append(opMapKeys, k)
			}
			return nil, fmt.Errorf("Invalid operator type (%s). Must be one of %s", operator, strings.Join(opMapKeys, ", "))
		}
	}

	v1, err := version.NewVersion(value)
	if err != nil {
		return nil, err
	}
	v2, err := version.NewVersion(ver)
	if err != nil {
		return nil, err
	}

	var truth bool
	switch op {
	case eq:
		truth = v1.Equal(v2)
	case lt:
		truth = v1.LessThan(v2)
	case le:
		truth = v1.LessThanOrEqual(v2)
	case gt:
		truth = v1.GreaterThan(v2)
	case ge:
		truth = v1.GreaterThanOrEqual(v2)
	case ne:
		truth = !v1.Equal(v2)
	default:
		panic(op)
	}

	return starlark.Bool(truth), nil
}
