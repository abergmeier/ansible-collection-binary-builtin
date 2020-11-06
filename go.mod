module github.com/abergmeier/ansible-collection-binary-builtin

go 1.15

replace github.com/go-git/go-git/v5 v5.2.0 => github.com/abergmeier/go-git/v5 v5.0.0-20201106103301-6be754ba63e9

require (
	github.com/go-git/go-git/v5 v5.2.0
	golang.org/x/crypto v0.0.0-20200302210943-78000ba7a073
)
