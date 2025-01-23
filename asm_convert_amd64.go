//go:build amd64 && !purego

package fftg

func cmplxToFloat4AssignAVX2(v []complex128, vOut []float64)

func float4ToCmplxAssignAVX2(v []float64, vOut []complex128)
