package rmq

import "testing"

func TestParseHostname(t *testing.T) {
	var testData = []struct {
		testName string
		nodeName string
		expected string
	}{
		{
			testName: "empty string",
			nodeName: "",
			expected: "",
		},
		{
			testName: "empty host",
			nodeName: "rabbit@",
			expected: "",
		},
		{
			testName: "valid host",
			nodeName: "rabbit@some_name",
			expected: "some_name",
		},
	}

	for _, v := range testData {
		v := v
		t.Run(v.testName, func(t *testing.T) {
			got := parseHostname(v.nodeName)
			if got != v.expected {
				t.Errorf("got %q, expected %q", got, v.expected)
			}
		})
	}
}
