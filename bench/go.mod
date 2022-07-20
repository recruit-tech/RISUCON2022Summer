module github.com/recruit-tech/RISUCON2022Summer/bench

go 1.17

replace github.com/recruit-tech/RISUCON2022Summer/snapshots/generator => ../snapshots/generator

require (
	github.com/google/uuid v1.3.0
	github.com/isucon/isucandar v0.0.0-20210921070917-929eaae2f9cb
	github.com/pkg/errors v0.9.1
	github.com/recruit-tech/RISUCON2022Summer/snapshots/generator v0.0.0-00010101000000-000000000000
	golang.org/x/exp v0.0.0-20220215214139-058d147d01d4
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
)

require (
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/pquerna/cachecontrol v0.0.0-20200819021114-67c6ae64274f // indirect
	golang.org/x/net v0.0.0-20211015210444-4f30a5c0130f // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
)
