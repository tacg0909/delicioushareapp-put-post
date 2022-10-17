package calctargetsize

import "testing"

func TestCalcTargetSizeWidthLongerThanHeight(t *testing.T) {
    w, h := CalcTargetSize(3840, 2160, 1000)
    if w / 16 != h / 9 {
        t.Fatalf(`CalcTargetSize(3840, 2160, 1000) = %d, %d, want match 16:9`, w, h)
    }
}

func TestCalcTargetSizeHeightLongerThanWidth(t *testing.T) {
    w, h := CalcTargetSize(2160, 3840, 1000)
    if w / 9 != h / 16 {
        t.Fatalf(`CalcTargetSize(2160, 3840, 1000) = %d, %d, want match 9:16`, w, h)
    }
}
