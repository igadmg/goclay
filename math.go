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

type Color colorex.RGBA
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

func (c Color) IsZero() bool {
	return c.R == 0 &&
		c.G == 0 &&
		c.B == 0 &&
		c.A == 0
}
