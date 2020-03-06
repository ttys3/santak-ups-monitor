module bitbucket.org/8ox86/santak-monitor

require (
	github.com/apex/log v1.1.0
	github.com/fatih/color v1.7.0 // indirect
	github.com/mattn/go-colorable v0.0.9 // indirect
	github.com/mattn/go-isatty v0.0.4 // indirect
	github.com/tarm/serial v0.0.0-20180830185346-98f6abe2eb07
	golang.org/x/sys v0.0.0-20181213200352-4d1cda033e06 // indirect
)

replace github.com/apex/log v1.1.0 => ../apex-log

go 1.13
