//go:build !(amd64 && !purego)

package fftg

// fftInPlace is a top-level function for FFT.
func fftInPlace(coeffs []complex128, _ []float64, tw []complex128, coeffsOut []complex128) {
	copy(coeffsOut, coeffs)
	fftInPlaceGeneric(coeffsOut, tw)
}

// invfftInPlace is a top-level function for inverse FFT.
func invfftInPlace(coeffs []complex128, _ []float64, twInv []complex128, coeffsOut []complex128) {
	copy(coeffsOut, coeffs)
	invfftInPlaceGeneric(coeffsOut, twInv)
}
