package pgtype

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"strconv"

	"github.com/jackc/pgx/pgio"
)

type Float4 struct {
	Float  float32
	Status Status
}

func (dst *Float4) Set(src interface{}) error {
	switch value := src.(type) {
	case float32:
		*dst = Float4{Float: value, Status: Present}
	case float64:
		*dst = Float4{Float: float32(value), Status: Present}
	case int8:
		*dst = Float4{Float: float32(value), Status: Present}
	case uint8:
		*dst = Float4{Float: float32(value), Status: Present}
	case int16:
		*dst = Float4{Float: float32(value), Status: Present}
	case uint16:
		*dst = Float4{Float: float32(value), Status: Present}
	case int32:
		f32 := float32(value)
		if int32(f32) == value {
			*dst = Float4{Float: f32, Status: Present}
		} else {
			return fmt.Errorf("%v cannot be exactly represented as float32", value)
		}
	case uint32:
		f32 := float32(value)
		if uint32(f32) == value {
			*dst = Float4{Float: f32, Status: Present}
		} else {
			return fmt.Errorf("%v cannot be exactly represented as float32", value)
		}
	case int64:
		f32 := float32(value)
		if int64(f32) == value {
			*dst = Float4{Float: f32, Status: Present}
		} else {
			return fmt.Errorf("%v cannot be exactly represented as float32", value)
		}
	case uint64:
		f32 := float32(value)
		if uint64(f32) == value {
			*dst = Float4{Float: f32, Status: Present}
		} else {
			return fmt.Errorf("%v cannot be exactly represented as float32", value)
		}
	case int:
		f32 := float32(value)
		if int(f32) == value {
			*dst = Float4{Float: f32, Status: Present}
		} else {
			return fmt.Errorf("%v cannot be exactly represented as float32", value)
		}
	case uint:
		f32 := float32(value)
		if uint(f32) == value {
			*dst = Float4{Float: f32, Status: Present}
		} else {
			return fmt.Errorf("%v cannot be exactly represented as float32", value)
		}
	case string:
		num, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return err
		}
		*dst = Float4{Float: float32(num), Status: Present}
	default:
		if originalSrc, ok := underlyingNumberType(src); ok {
			return dst.Set(originalSrc)
		}
		return fmt.Errorf("cannot convert %v to Float8", value)
	}

	return nil
}

func (dst *Float4) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst.Float
	case Null:
		return nil
	default:
		return dst.Status
	}
}

func (src *Float4) AssignTo(dst interface{}) error {
	return float64AssignTo(float64(src.Float), src.Status, dst)
}

func (dst *Float4) DecodeText(src []byte) error {
	if src == nil {
		*dst = Float4{Status: Null}
		return nil
	}

	n, err := strconv.ParseFloat(string(src), 32)
	if err != nil {
		return err
	}

	*dst = Float4{Float: float32(n), Status: Present}
	return nil
}

func (dst *Float4) DecodeBinary(src []byte) error {
	if src == nil {
		*dst = Float4{Status: Null}
		return nil
	}

	if len(src) != 4 {
		return fmt.Errorf("invalid length for float4: %v", len(src))
	}

	n := int32(binary.BigEndian.Uint32(src))

	*dst = Float4{Float: math.Float32frombits(uint32(n)), Status: Present}
	return nil
}

func (src Float4) EncodeText(w io.Writer) (bool, error) {
	switch src.Status {
	case Null:
		return true, nil
	case Undefined:
		return false, errUndefined
	}

	_, err := io.WriteString(w, strconv.FormatFloat(float64(src.Float), 'f', -1, 32))
	return false, err
}

func (src Float4) EncodeBinary(w io.Writer) (bool, error) {
	switch src.Status {
	case Null:
		return true, nil
	case Undefined:
		return false, errUndefined
	}

	_, err := pgio.WriteInt32(w, int32(math.Float32bits(src.Float)))
	return false, err
}
