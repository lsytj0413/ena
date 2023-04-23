package ena

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestPointerTo(t *testing.T) {
	g := NewWithT(t)

	v1 := "123"
	g.Expect(PointerTo("123")).To(Equal(&v1))

	v2 := []int{1, 2}
	g.Expect(PointerTo(v2)).To(Equal(&v2))
}

func TestValueTo(t *testing.T) {
	g := NewWithT(t)

	v1 := "123"
	g.Expect(ValueTo(&v1)).To(Equal(v1))

	v2 := []int{1, 2}
	g.Expect(ValueTo(&v2)).To(Equal(v2))

	var v3 *string
	g.Expect(ValueTo(v3)).To(Equal(""))
}

func TestValueToOrDefault(t *testing.T) {
	g := NewWithT(t)

	v1 := "123"
	g.Expect(ValueToOrDefault(&v1, "456")).To(Equal(v1))

	var v3 *string
	g.Expect(ValueToOrDefault(v3, "456")).To(Equal("456"))
}

func TestEnforcePtr(t *testing.T) {
	type testCase struct {
		desp string
		obj  interface{}

		err string
	}
	testCases := []testCase{
		{
			desp: "normal test",
			obj:  &testCase{},
			err:  ``,
		},
		{
			desp: "not pointer",
			obj:  testCase{},
			err:  `expected pointer, but got ena.testCase type: UnexpectPointer`,
		},
		{
			desp: "invalid kind",
			obj:  nil,
			err:  `expected pointer, but got invalid kind: UnexpectPointer`,
		},
		{
			desp: "nil value",
			obj:  (*testCase)(nil),
			err:  `expected pointer, but got nil: UnexpectPointer`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			_, err := EnforcePtr(tc.obj)
			if tc.err != `` {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
		})
	}
}
