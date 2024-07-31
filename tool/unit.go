package tool

import "math/big"

var (
	weiUnit        = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil) // 10^0
	kweiUnit       = new(big.Int).Exp(big.NewInt(10), big.NewInt(15), nil) // 10^3
	mweiUnit       = new(big.Int).Exp(big.NewInt(10), big.NewInt(12), nil) // 10^6
	gweiUnit       = new(big.Int).Exp(big.NewInt(10), big.NewInt(9), nil)  // 10^9
	microetherUnit = new(big.Int).Exp(big.NewInt(10), big.NewInt(6), nil)  // 10^12
	millietherUnit = new(big.Int).Exp(big.NewInt(10), big.NewInt(3), nil)  // 10^15
	etherUnit      = new(big.Int).Exp(big.NewInt(10), big.NewInt(0), nil)  // 10^18
)

func ConvertWeiToUnit(wei *big.Int, _unit string) string {
	var unit *big.Int
	switch _unit {
	case "WEI":
		unit = weiUnit
	case "KWEI":
		unit = kweiUnit
	case "MWEI":
		unit = mweiUnit
	case "GWEI":
		unit = gweiUnit
	case "MICROETHER":
		unit = microetherUnit
	case "MILLIETHER":
		unit = millietherUnit
	case "ETHER":
		unit = etherUnit
	case "FRI":
		unit = weiUnit
	}
	weiFloat := new(big.Float).SetInt(wei)
	unitFloat := new(big.Float).SetInt(unit)
	return new(big.Float).Quo(weiFloat, unitFloat).Text('f', 18)
}
