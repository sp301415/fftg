package fftg

// fftInPlaceGeneric is a top-level function for FFT.
func fftInPlaceGeneric(coeffs []complex128, tw []complex128) {
	N := len(coeffs)

	t := N
	for m := 1; m < N; m <<= 1 {
		t >>= 1
		for i := 0; i < m; i++ {
			j1 := i * t << 1
			j2 := j1 + t
			for j := j1; j < j2; j++ {
				U, V := coeffs[j], coeffs[j+t]*tw[i]
				coeffs[j], coeffs[j+t] = U+V, U-V
			}
		}
	}
}

// invfftInPlaceGeneric is a top-level function for inverse FFT.
func invfftInPlaceGeneric(coeffs []complex128, twInv []complex128) {
	N := len(coeffs)

	t := 1
	for m := N; m > 1; m >>= 1 {
		j1 := 0
		h := m >> 1
		for i := 0; i < h; i++ {
			j2 := j1 + t
			for j := j1; j < j2; j++ {
				U, V := coeffs[j], coeffs[j+t]
				coeffs[j], coeffs[j+t] = U+V, (U-V)*twInv[i]
			}
			j1 += t << 1
		}
		t <<= 1
	}
}
