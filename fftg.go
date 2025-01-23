// package fftg impelments the Fast Fourier Transform (FFT) in Go.
package fftg

import (
	"math"
)

const (
	// MinDegree is the minimum degree of polynomial that Transformers can handle.
	// Currently, this is set to 16 = 2^4, because AVX2 implementation of FFT and inverse FFT
	// handles first/last two loops separately.
	MinDegree = 1 << 4
)

// FourierTransformer computes FFT/IFFT over C[X]/(X^N-1).
// N should be a power of two.
//
// Note the following caveats:
//   - During forward transform, the output vector is in bit-reversed order.
//     Consequently, the input vector of inverse transform should be in bit-reversed order.
//     This is enough for computing convolutions.
//     If you need the result in natural order, use [BitReverseInPlace].
//   - The output vector of inverse transform is NOT normalized.
type FourierTransformer struct {
	// degree is the degree of the polynomial.
	degree int

	// tw is the twiddle factors for FFT.
	// Unlike other complex128 slices, tw is in natural representation.
	tw []complex128
	// twInv is the twiddle factors for IFFT.
	// Unlike other complex128 slices, twInv is in natural representation.
	twInv []complex128

	// buffFloat4 is a buffer for float4 representation of complex128 vector.
	buffFloat4 []float64
}

// NewFourierTransformer creates a new CyclicTransformer.
func NewFourierTransformer(N int) *FourierTransformer {
	switch {
	case !IsPowerOfTwo(N):
		panic("N should be a power of two")
	case N < MinDegree:
		panic("N should be at least MinDegree")
	}

	// Compute twiddle factors.
	twRef := make([]complex128, N/2)
	twInvRef := make([]complex128, N/2)
	for i := 0; i < N/2; i++ {
		e := -2 * math.Pi * float64(i) / float64(N)
		twRef[i] = complex(math.Cos(e), math.Sin(e))
		twInvRef[i] = complex(math.Cos(-e), math.Sin(-e))
	}
	BitReverseInPlace(twRef)
	BitReverseInPlace(twInvRef)

	tw := make([]complex128, 0, N)
	for m := 1; m < N; m <<= 1 {
		for i := 0; i < m; i++ {
			tw = append(tw, twRef[i])
		}
	}

	twInv := make([]complex128, 0, N)
	for m := N; m > 1; m >>= 1 {
		for i := 0; i < m>>1; i++ {
			twInv = append(twInv, twInvRef[i])
		}
	}

	return &FourierTransformer{
		degree: N,

		tw:    tw,
		twInv: twInv,

		buffFloat4: make([]float64, 2*N),
	}
}

// Degree returns the degree of this transformer.
func (fft *FourierTransformer) Degree() int {
	return fft.degree
}

// Forward computes the fourier transform of v and returns it.
// Input should be in natural order, and output will be in bit-reversed order.
func (fft *FourierTransformer) Forward(v []complex128) []complex128 {
	if len(v) != fft.degree {
		panic("invalid length of input vector")
	}

	vOut := make([]complex128, fft.degree)
	fft.ForwardAssign(v, vOut)
	return vOut
}

// ForwardAssign computes the fourier transform of v and writes it to vOut.
// Input should be in natural order, and output will be in bit-reversed order.
func (fft *FourierTransformer) ForwardAssign(v []complex128, vOut []complex128) {
	if len(v) != fft.degree && len(vOut) != fft.degree {
		panic("invalid length of input vector")
	}

	fftInPlace(v, fft.buffFloat4, fft.tw, vOut)
}

// Inverse computes the inverse fourier transform of v and returns it.
// Input should be in bit-reversed order, and output will be in natural order.
// The output vector is NOT normalized.
func (fft *FourierTransformer) Inverse(v []complex128) []complex128 {
	if len(v) != fft.degree {
		panic("invalid length of input vector")
	}

	vOut := make([]complex128, fft.degree)
	fft.InverseAssign(v, vOut)
	return vOut
}

// InverseAssign computes the inverse fourier transform of v and writes it to vOut.
// Input should be in bit-reversed order, and output will be in natural order.
// The output vector is NOT normalized.
func (fft *FourierTransformer) InverseAssign(v []complex128, vOut []complex128) {
	if len(v) != fft.degree && len(vOut) != fft.degree {
		panic("invalid length of input vector")
	}

	invfftInPlace(v, fft.buffFloat4, fft.twInv, vOut)
}
