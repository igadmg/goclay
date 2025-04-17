module github.com/igadmg/goclay

go 1.24

replace github.com/igadmg/gamemath => ../gamemath

require (
	github.com/igadmg/gamemath v0.0.0-20250410222204-28d83654fdf2
	golang.org/x/exp v0.0.0-20250408133849-7e4ce0ab07d0
)

require github.com/chewxy/math32 v1.11.1 // indirect
