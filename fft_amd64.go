//go:build amd64 && !purego

package fftg

import "golang.org/x/sys/cpu"

// fftInPlace is a top-level function for FFT.
func fftInPlace(coeffs []complex128, buff []float64, tw []complex128, coeffsOut []complex128) {
	if cpu.X86.HasAVX2 && cpu.X86.HasFMA && len(coeffs) >= MinDegree {
		cmplxToFloat4AssignAVX2(coeffs, buff)
		fftInPlaceAVX2(buff, tw)
		float4ToCmplxAssignAVX2(buff, coeffsOut)
		return
	}

	copy(coeffsOut, coeffs)
	fftInPlaceGeneric(coeffsOut, tw)
}

// invfftInPlace is a top-level function for inverse FFT.
func invfftInPlace(coeffs []complex128, buff []float64, twInv []complex128, coeffsOut []complex128) {
	if cpu.X86.HasAVX2 && cpu.X86.HasFMA && len(coeffs) >= MinDegree {
		cmplxToFloat4AssignAVX2(coeffs, buff)
		invfftInPlaceAVX2(buff, twInv)
		float4ToCmplxAssignAVX2(buff, coeffsOut)
		return
	}

	copy(coeffsOut, coeffs)
	invfftInPlaceGeneric(coeffsOut, twInv)
}
