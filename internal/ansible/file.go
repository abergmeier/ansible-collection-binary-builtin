package ansible

import (
	"os"

	"go.starlark.net/starlark"
)

func Directory(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	path, ok := starlark.AsString(args[0])
	if !ok {
		panic(ok)
	}
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return starlark.False, nil
		}
		return nil, err
	}
	return starlark.Bool(fi.IsDir()), nil
}

// Exists checks file
func Exists(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	path, ok := starlark.AsString(args[0])
	if !ok {
		panic(ok)
	}
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return starlark.False, nil
		}
		return nil, err
	}
	return starlark.True, nil
}

func File(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	path, ok := starlark.AsString(args[0])
	if !ok {
		panic(ok)
	}
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return starlark.False, nil
		}
		return nil, err
	}
	return starlark.Bool(!fi.IsDir()), nil
}
