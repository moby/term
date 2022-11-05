module github.com/moby/term/test

go 1.18

require (
	github.com/creack/pty v1.1.18
	github.com/moby/term v0.0.0-00010101000000-000000000000 // replaced
)

require (
	github.com/Azure/go-ansiterm v0.0.0-20250102033503-faa5f7b0171c // indirect
	golang.org/x/sys v0.1.0 // indirect
)

replace github.com/moby/term => ../
