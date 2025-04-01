module github.com/igadmg/goclay

go 1.24

replace github.com/igadmg/gamemath => ../gamemath

require (
	github.com/igadmg/gamemath v0.0.0-00010101000000-000000000000
	golang.org/x/exp v0.0.0-20250305212735-054e65f0b394
)

require github.com/chewxy/math32 v1.11.1 // indirect
