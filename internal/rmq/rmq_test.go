package rmq

import (
	"testing"

	"github.com/popovous/rabbitmq-healthcheck/internal/clusterinfo"
)

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

func TestIsInCluster(t *testing.T) {
	var testData = []struct {
		testName          string
		clusterInfo       []clusterinfo.Members
		hostname          string
		expectedIsRunning bool
		expectedIsAlone   bool
	}{
		{
			testName:    "empty cluster info",
			clusterInfo: nil,
		},
		{
			testName: "all running nodes",
			clusterInfo: []clusterinfo.Members{
				{
					Name:    "rabbit@1",
					Running: true,
				},
				{
					Name:    "rabbit@2",
					Running: true,
				},
			},
			expectedIsRunning: true,
			hostname:          "2",
			expectedIsAlone:   false,
		},
		{
			testName: "only one running node",
			clusterInfo: []clusterinfo.Members{
				{
					Name:    "rabbit@1",
					Running: false,
				},
				{
					Name:    "rabbit@2",
					Running: false,
				},
				{
					Name:    "rabbit@3",
					Running: true,
				},
			},
			hostname:          "3",
			expectedIsRunning: true,
			expectedIsAlone:   true,
		},
		{
			testName: "one not-running node in cluster",
			clusterInfo: []clusterinfo.Members{
				{
					Name:    "rabbit@1",
					Running: true,
				},
				{
					Name:    "rabbit@2",
					Running: true,
				},
				{
					Name:    "rabbit@3",
					Running: false,
				},
			},
			hostname:          "2",
			expectedIsRunning: true,
			expectedIsAlone:   false,
		},
		{
			testName: "current not-running node",
			clusterInfo: []clusterinfo.Members{
				{
					Name:    "rabbit@1",
					Running: true,
				},
				{
					Name:    "rabbit@2",
					Running: false,
				},
				{
					Name:    "rabbit@3",
					Running: true,
				},
			},
			hostname:          "2",
			expectedIsRunning: false,
			expectedIsAlone:   false,
		},
	}

	for _, v := range testData {
		v := v
		t.Run(v.testName, func(t *testing.T) {
			isRunning, isAlone := isInCluster(v.clusterInfo, v.hostname)
			if isRunning != v.expectedIsRunning {
				t.Errorf("got isRunning = %v, expected %v", isRunning, v.expectedIsRunning)
			}
			if isAlone != v.expectedIsAlone {
				t.Errorf("got isAlone = %v, expected %v", isAlone, v.expectedIsAlone)
			}
		})
	}
}
