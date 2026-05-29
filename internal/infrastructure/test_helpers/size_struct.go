package test_helpers

import (
	"reflect"
	"sort"
	"testing"
	"unsafe"
)

func ActualAndGreedySizeOfStruct(t *testing.T, opts any) (actualSize uintptr, greedySize uintptr) {
	t.Helper()

	actualSize = unsafe.Sizeof(opts)

	typ := reflect.TypeOf(opts)

	var fields []reflect.StructField
	for i := 0; i < typ.NumField(); i++ {
		fields = append(fields, typ.Field(i))
	}

	sorted := make([]reflect.StructField, len(fields))
	copy(sorted, fields)

	sort.Slice(sorted, func(i, j int) bool {
		ti, tj := sorted[i].Type, sorted[j].Type
		ai, aj := ti.Align(), tj.Align()
		if ai != aj {
			return ai > aj
		}
		return ti.Size() > tj.Size()
	})

	greedySize = layoutSize(t, sorted)
	return actualSize, greedySize
}

func layoutSize(t *testing.T, fields []reflect.StructField) uintptr {
	t.Helper()

	var offset uintptr
	var maxAlign uintptr

	for _, f := range fields {
		size := f.Type.Size()
		align := uintptr(f.Type.Align())
		if align > maxAlign {
			maxAlign = align
		}
		if rem := offset % align; rem != 0 {
			offset += align - rem
		}
		offset += size
	}

	if maxAlign == 0 {
		return 0
	}

	if rem := offset % maxAlign; rem != 0 {
		offset += maxAlign - rem
	}

	return offset
}
