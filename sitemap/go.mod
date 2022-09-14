module github.com/AksAman/gophercises/sitemap

replace github.com/AksAman/gophercises/linkparser => ../linkparser

go 1.19

require (
	github.com/AksAman/gophercises/linkparser v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.23.0
)

require (
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/net v0.0.0-20220909164309-bea034e7d591 // indirect
)
