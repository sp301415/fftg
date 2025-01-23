//go:build amd64 && !purego

package fftg

func fftInPlaceAVX2(coeffs []float64, tw []complex128)

func invfftInPlaceAVX2(coeffs []float64, twInv []complex128)
