package flag_parse

import (
	"strconv"
	"testing"

	_ "github.com/stretchr/testify"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/infrastructure/test_helpers"
)

func TestMyInt16_Set(t *testing.T) {
	t.Parallel()
	type TestCase struct {
		name string
		arg  string
		want MyInt16
		err  error
	}
	tableCases := []TestCase{
		{
			name: "simple test",
			arg:  "100",
			want: MyInt16(100),
			err:  nil,
		},
		{
			name: "simple test",
			arg:  "-100",
			want: MyInt16(-100),
			err:  nil,
		},
		{
			name: "error test",
			arg:  "fffff",
			want: 0,
			err:  strconv.ErrSyntax,
		},
	}
	for _, tc := range tableCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var a MyInt16
			err := a.Set(tc.arg)

			if tc.err != nil {
				assert.ErrorIs(t, err, tc.err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, a)
		})
	}
}

func TestMyInt16_String(t *testing.T) {
	t.Parallel()
	type TestCase struct {
		name string
		arg  MyInt16
		want string
		err  error
	}
	tableCases := []TestCase{
		{
			name: "simple test",
			want: "100",
			arg:  MyInt16(100),
			err:  nil,
		},
		{
			name: "simple test",
			want: "-100",
			arg:  MyInt16(-100),
			err:  nil,
		},
		{
			name: "error test",
			want: "0",
			arg:  0,
			err:  strconv.ErrSyntax,
		},
	}
	for _, tc := range tableCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var a = tc.arg
			str := a.String()

			assert.Equal(t, tc.want, str)

		})
	}
}

func TestAffineParams_Set(t *testing.T) {
	t.Parallel()
	type TestCase struct {
		name string
		arg  string
		want AffineParams
		err  string
	}
	tableCases := []TestCase{
		{
			name: "simple test",
			arg:  "1,2,3,4,5,6",
			want: AffineParams{{
				A: 1,
				B: 2,
				C: 3,
				D: 4,
				E: 5,
				F: 6,
			}},
			err: "",
		},
		{
			name: "simple test",
			arg:  "1.0,   2.5,   3.3, 4.4, 5.111, 6",
			want: AffineParams{{
				A: 1.0,
				B: 2.5,
				C: 3.3,
				D: 4.4,
				E: 5.111,
				F: 6,
			}},
			err: "",
		},
		{
			name: "error test",
			arg:  "fffff",
			want: AffineParams{},
			err:  "invalid affine block",
		},
		{
			name: "error test",
			arg:  "f, f, f, f, f, f",
			want: AffineParams{},
			err:  "invalid float",
		},
		{
			name: "error test",
			arg:  "f, f, f, f, f, f, f",
			want: AffineParams{},
			err:  "invalid affine block",
		},
	}
	for _, tc := range tableCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var a AffineParams
			err := a.Set(tc.arg)

			if tc.err != "" {
				assert.ErrorContains(t, err, tc.err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, a)

		})
	}
}

func TestAffineParams_String(t *testing.T) {
	t.Parallel()
	type TestCase struct {
		name string
		arg  AffineParams
		want string
		err  string
	}
	tableCases := []TestCase{
		{
			name: "simple test",
			want: "1,2,3,4,5,6",
			arg: AffineParams{{
				A: 1,
				B: 2,
				C: 3,
				D: 4,
				E: 5,
				F: 6,
			}},
			err: "",
		},
		{
			name: "simple test",
			want: "1,2.5,3.3,4.4,5.111,6",
			arg: AffineParams{{
				A: 1.0,
				B: 2.5,
				C: 3.3,
				D: 4.4,
				E: 5.111,
				F: 6,
			}},
			err: "",
		},
	}
	for _, tc := range tableCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var a = tc.arg
			str := a.String()

			assert.Equal(t, tc.want, str)
		})
	}
}

func TestFuncsParams_Set(t *testing.T) {
	t.Parallel()
	type TestCase struct {
		name string
		arg  string
		want FuncsParams
		err  string
	}
	tableCases := []TestCase{
		{
			name: "simple test",
			arg:  "swirl:1.5,horseshoe:1",
			want: FuncsParams{{
				Name:   "swirl",
				Weight: 1.5,
			}, {
				Name:   "horseshoe",
				Weight: 1,
			}},
			err: "",
		},
		{
			name: "simple test",
			arg:  "swirl:131",
			want: FuncsParams{{
				Name:   "swirl",
				Weight: 131,
			}},
			err: "",
		},
		{
			name: "error test",
			arg:  "",
			want: FuncsParams{},
			err:  "functions string is empty",
		},
		{
			name: "error test",
			arg:  "f, f, f, f, f, f",
			want: FuncsParams{},
			err:  "invalid function spec",
		},
		{
			name: "error test",
			arg:  "swirl 1.swirl 2",
			want: FuncsParams{},
			err:  "invalid function spec",
		},
	}
	for _, tc := range tableCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var a FuncsParams
			err := a.Set(tc.arg)

			if tc.err != "" {
				assert.ErrorContains(t, err, tc.err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, a)

		})
	}
}

func TestFuncsParams_String(t *testing.T) {
	t.Parallel()
	type TestCase struct {
		name string
		want string
		arg  FuncsParams
		err  string
	}
	tableCases := []TestCase{
		{
			name: "simple test",
			want: "swirl:1.5,horseshoe:1",
			arg: FuncsParams{{
				Name:   "swirl",
				Weight: 1.5,
			}, {
				Name:   "horseshoe",
				Weight: 1,
			}},
			err: "",
		},
		{
			name: "simple test",
			want: "swirl:131",
			arg: FuncsParams{{
				Name:   "swirl",
				Weight: 131,
			}},
			err: "",
		},
	}
	for _, tc := range tableCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var a = tc.arg
			str := a.String()

			assert.Equal(t, tc.want, str)
		})
	}
}
func TestOpts_LayoutIsReasonable(t *testing.T) {
	t.Helper()
	t.Parallel()
	var opts Opts
	actualSize, greedySize := test_helpers.ActualAndGreedySizeOfStruct(t, opts)

	t.Logf("actual size: %d, greedy size: %d", actualSize, greedySize)

	if greedySize < actualSize {
		t.Fatalf("suboptimal layout for Opts: actual=%d, greedy=%d",
			actualSize, greedySize)
	}
}
