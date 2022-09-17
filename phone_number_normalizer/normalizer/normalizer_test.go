package normalizer

import "testing"

type PhoneTestCase struct {
	caseName              string
	phoneNumber           string
	normalizedPhoneNumber string
}

func compare(t *testing.T, got, want any) {
	if got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func TestNormalizePhoneNumberSimple(t *testing.T) {
	testCase := "1234567890"
	want := "1234567890"
	got := NormalizePhoneNumber(testCase)
	compare(t, got, want)
}

func TestNormalizePhoneNumberAll(t *testing.T) {
	testcases := []PhoneTestCase{
		{caseName: "normal", phoneNumber: "1234567890", normalizedPhoneNumber: "1234567890"},
		{caseName: "spaces", phoneNumber: "123 456 7891", normalizedPhoneNumber: "1234567891"},
		{caseName: "dashes", phoneNumber: "123-456-7894", normalizedPhoneNumber: "1234567894"},
		{caseName: "parenthesis-spaces", phoneNumber: "(123) 456 7892", normalizedPhoneNumber: "1234567892"},
		{caseName: "parenthesis-dashes", phoneNumber: "(123)456-7892", normalizedPhoneNumber: "1234567892"},
		{caseName: "parenthesis-spaces-dashes", phoneNumber: "(123) 456-7893", normalizedPhoneNumber: "1234567893"},
	}

	for _, testcase := range testcases {
		t.Run(
			testcase.caseName,
			func(t *testing.T) {
				got := NormalizePhoneNumber(testcase.phoneNumber)
				want := testcase.normalizedPhoneNumber
				compare(t, got, want)
			},
		)
	}
}
