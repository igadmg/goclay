module github.com/igadmg/goclay

go 1.24

replace github.com/igadmg/goex => ../goex

replace github.com/igadmg/raylib-go/raymath => ../raylib-go/raymath

require (
	github.com/igadmg/goex v0.0.0-20250321131421-ccb743b21181
	github.com/igadmg/raylib-go/raymath v0.0.0-00010101000000-000000000000
	golang.org/x/exp v0.0.0-20250305212735-054e65f0b394
)

require (
	deedles.dev/xiter v0.2.1 // indirect
	github.com/chewxy/math32 v1.11.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
