package fftg_test

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/sp301415/fftg"
	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/dsp/fourier"
)

var (
	logN    = []int{8, 9, 10, 11, 12, 13, 14, 15, 16}
	v, vOut []complex128
)

func TestFFT(t *testing.T) {
	logN := 10
	eps := math.Exp2(float64(logN)) * 1e-10

	fft := fftg.NewFourierTransformer(1 << logN)
	fftTest := fourier.NewCmplxFFT(1 << logN)

	v := make([]complex128, 1<<logN)
	for i := 0; i < 1<<logN; i++ {
		v[i] = complex(float64(2*i), float64(2*i+1))
	}

	t.Run("Forward", func(t *testing.T) {
		vForward := fft.Forward(v)
		fftg.BitReverseInPlace(vForward)

		vForwardTest := fftTest.Coefficients(make([]complex128, 1<<logN), v)

		for i := 0; i < 1<<logN; i++ {
			assert.InDelta(t, real(vForward[i]), real(vForwardTest[i]), eps)
			assert.InDelta(t, imag(vForward[i]), imag(vForwardTest[i]), eps)
		}
	})

	t.Run("Inverse", func(t *testing.T) {
		fftg.BitReverseInPlace(v)
		vInverse := fft.Inverse(v)
		fftg.BitReverseInPlace(v)

		vInverseTest := fftTest.Sequence(make([]complex128, 1<<logN), v)

		for i := 0; i < 1<<logN; i++ {
			assert.InDelta(t, real(vInverse[i]), real(vInverseTest[i]), eps)
			assert.InDelta(t, imag(vInverse[i]), imag(vInverseTest[i]), eps)
		}
	})
}

func BenchmarkFFT(b *testing.B) {
	for _, n := range logN {
		fft := fftg.NewFourierTransformer(1 << n)

		v = make([]complex128, 1<<n)
		vOut = make([]complex128, 1<<n)
		for i := 0; i < 1<<n; i++ {
			v[i] = complex(rand.Float64(), rand.Float64())
			vOut[i] = complex(rand.Float64(), rand.Float64())
		}

		b.Run(fmt.Sprintf("Forward/LogN=%v", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				fft.ForwardAssign(v, vOut)
			}
		})

		b.Run(fmt.Sprintf("Inverse/LogN=%v", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				fft.InverseAssign(vOut, v)
			}
		})
	}
}
