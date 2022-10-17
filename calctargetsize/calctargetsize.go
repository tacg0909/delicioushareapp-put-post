package calctargetsize

import "math/big"

func CalcTargetSize(width int, height int, maxLength int) (targetWidth int, targetHeight int) {
    gcd := int(new(big.Int).GCD(nil, nil, big.NewInt(int64(width)), big.NewInt(int64(height))).Int64())
    widthRate := width / gcd
    heightRate := height / gcd
    if widthRate > heightRate {
        rate := maxLength / widthRate
        targetWidth = widthRate * rate
        targetHeight = heightRate * rate
        return
    }
    rate := maxLength / heightRate
    targetWidth = widthRate * rate
    targetHeight = heightRate * rate
    return
}
