package clay

import (
	gm "github.com/igadmg/gamemath"
	"github.com/igadmg/gamemath/rect2"
	"github.com/igadmg/gamemath/vector2"
	"github.com/igadmg/goex/image/colorex"
	"golang.org/x/exp/constraints"
)

type Coordinate interface {
	constraints.Integer | constraints.Float
}

type Axis = gm.Axis

const (
	AxisX = gm.AxisX
	AxisY = gm.AxisY
)

type Color = colorex.RGBA
type Vector2 = vector2.Float32
type Dimensions = vector2.Float32
type BoundingBox = rect2.Float32

func MakeDimensions[XT, YT Coordinate](x XT, y YT) Dimensions {
	return vector2.NewFloat32(x, y)
}

func MakeVector2[XT, YT Coordinate](x XT, y YT) Vector2 {
	return vector2.NewFloat32(x, y)
}

func MakeBoundingBox(position Vector2, size Dimensions) BoundingBox {
	return rect2.NewFloat32(position, size)
}

var CLAY__EPSILON float32 = 0.01

func floatEqual(left float32, right float32) bool {
	subtracted := left - right
	return subtracted < CLAY__EPSILON && subtracted > -CLAY__EPSILON
}
