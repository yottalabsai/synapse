package utils

import (
	"errors"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/shopspring/decimal"
	"math/big"
	"strconv"
)

func GenerateAccuracyDivisor(n int) (int64, error) {
	if n <= 2 {
		return 0, errors.New("n must gt 2")
	}

	v := "1"
	for i := 0; i < n; i++ {
		v += "0"
	}

	return strconv.ParseInt(v, 10, 64)
}

// USDTConvertFunc 转换hex字符串到10进制，按accuracy保留位数/**
func USDTConvertFunc(amount string) (float64, error) {
	balance, _ := math.ParseBig256(amount)

	divisor, err := GenerateAccuracyDivisor(6)
	if err != nil {
		return 0, err
	}

	return float64(new(big.Int).Div(balance, big.NewInt(divisor)).Int64()), nil
}

func ETHConvertFunc(amount string) (float64, error) {
	v, _ := ToDecimal(amount, 18)
	return v, nil
}

func ToDecimal(iValue string, decimals int) (float64, bool) {
	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	num, _ := decimal.NewFromString(iValue)
	result := num.Div(mul)

	return result.Float64()
}

func ToWei(amount interface{}, decimals int) *big.Int {
	weiAmount := decimal.NewFromFloat(0)
	switch v := amount.(type) {
	case string:
		weiAmount, _ = decimal.NewFromString(v)
	case float64:
		weiAmount = decimal.NewFromFloat(v)
	case int64:
		weiAmount = decimal.NewFromFloat(float64(v))
	case decimal.Decimal:
		weiAmount = v
	case *decimal.Decimal:
		weiAmount = *v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	result := weiAmount.Mul(mul)

	wei := new(big.Int)
	wei.SetString(result.String(), 10)

	return wei
}
