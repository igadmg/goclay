module github.com/igadmg/goclay

go 1.24

replace github.com/igadmg/gamemath => ../gamemath

require (
	github.com/igadmg/gamemath v0.0.0-20250404225837-2cb2391130a0
	golang.org/x/exp v0.0.0-20250305212735-054e65f0b394
)

require github.com/chewxy/math32 v1.11.1 // indirect
