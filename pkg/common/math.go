package common

import (
	"math/big"
)

type Float big.Float

var floatPrecision = uint(32)

func SetFloatPrecision(precision uint) {
	floatPrecision = precision
}

func float() *big.Float {
	return new(big.Float).SetPrec(floatPrecision)
}

func NewFloat(value float64) Float {
	return Float(*float().SetFloat64(value))
}

func AddFloat(a, b big.Float) Float {
	return Float(*float().Add(&a, &b))
}

func SubFloat(a, b big.Float) Float {
	return Float(*float().Sub(&a, &b))
}

func MulFloat(a, b big.Float) Float {
	return Float(*float().Mul(&a, &b))
}

func DivFloat(a, b big.Float) Float {
	return Float(*float().Quo(&a, &b))
}

func (f Float) Float64() float64 {
	ff := big.Float(f)
	res, _ := (&ff).Float64()
	return res
}
