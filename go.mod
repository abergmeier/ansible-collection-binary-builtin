module github.com/abergmeier/ansible-collection-binary-builtin

go 1.15

replace github.com/go-git/go-git/v5 v5.2.0 => github.com/abergmeier/go-git/v5 v5.0.0-20201106103301-6be754ba63e9

require (
	github.com/go-git/go-git/v5 v5.2.0
	github.com/hashicorp/go-version v1.2.1
	go.starlark.net v0.0.0-20201014215153-dff0ae5b4820
	golang.org/x/crypto v0.1.0
)
