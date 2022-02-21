package flags

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Checkmarx/kics/internal/constants"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestFlags_GetAllFlags(t *testing.T) {
	boolFlag := true
	flagsBoolReferences["bool"] = &boolFlag
	intFlag := 5
	flagsIntReferences["int"] = &intFlag
	multiStrFlag := []string{"test"}
	flagsMultiStrReferences["multi"] = &multiStrFlag
	strFlag := "test"
	flagsStrReferences["str"] = &strFlag

	expectedFlags := map[string]interface{}{
		"bool":  &boolFlag,
		"int":   &intFlag,
		"multi": &multiStrFlag,
		"str":   &strFlag,
	}

	gotFlags := GetAllFlags()

	require.Equal(t, expectedFlags, gotFlags)
}

func TestFlags_InitJSONFlags(t *testing.T) {
	mockCmd := &cobra.Command{
		Use:   "mock",
		Short: "Mock cmd",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	tests := []struct {
		name                    string
		cmd                     *cobra.Command
		flagsListContent        string
		persintentFlag          bool
		supportedPlatforms      []string
		supportedCloudProviders []string
		wantErr                 bool
	}{
		{
			name: "should initialize flags without error",
			cmd:  mockCmd,
			flagsListContent: `{"log-level": {
				"flagType": "str",
				"shorthandFlag": "",
				"defaultValue": "INFO",
				"usage": "determines log level (${supportedLogLevels})",
				"validation": "validateStrEnum"
			},"preview-lines": {
				"flagType": "int",
				"shorthandFlag": "",
				"defaultValue": "3",
				"usage": "number of lines to be display in CLI results (min: 1, max: 30)"
			},"queries-path": {
				"flagType": "multiStr",
				"shorthandFlag": "q",
				"defaultValue": "./assets/queries",
				"usage": "paths to directory with queries"
			},"verbose": {
				"flagType": "bool",
				"shorthandFlag": "v",
				"defaultValue": "false",
				"usage": "write logs to stdout too (mutually exclusive with silent)"
			},"disable-cis-descriptions": {
				"flagType": "bool",
				"shorthandFlag": "",
				"defaultValue": "false",
				"usage": "disable request for CIS descriptions and use default vulnerability descriptions",
				"hidden": true,
				"deprecated": true,
				"deprecatedInfo": "use --disable-full-descriptions instead"
			}}`,
			persintentFlag:          true,
			supportedPlatforms:      []string{"terraform"},
			supportedCloudProviders: []string{"aws"},
			wantErr:                 false,
		},
		{
			name: "should throw error due to wrong json marshal on flagListContent",
			cmd:  mockCmd,
			flagsListContent: `"verbose": {
				"flagType": "bool",
				"shorthandFlag": "v",
				"defaultValue": "false",
				"usage": "write logs to stdout too (mutually exclusive with silent)"
			}`,
			persintentFlag:          true,
			supportedPlatforms:      []string{"terraform"},
			supportedCloudProviders: []string{"aws"},
			wantErr:                 true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := InitJSONFlags(test.cmd, test.flagsListContent, test.persintentFlag, test.supportedCloudProviders, test.supportedCloudProviders)
			if !test.wantErr {
				require.NoError(t, got)
			} else {
				require.Error(t, got)
			}
		})
	}
}

func TestFlags_GetStrFlag(t *testing.T) {
	tests := []struct {
		name     string
		flagName string
		expected string
	}{
		{
			name:     "should return value for valid flag",
			flagName: "test",
			expected: "exists",
		},
		{
			name:     "should not return value for invalid flag",
			flagName: "undefined",
			expected: "",
		},
	}
	existValue := "exists"
	flagsStrReferences["test"] = &existValue
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := GetStrFlag(test.flagName)
			require.Equal(t, test.expected, got)
		})
	}
}

func TestFlags_GetMultiStrFlag(t *testing.T) {
	tests := []struct {
		name     string
		flagName string
		expected []string
	}{
		{
			name:     "should return value for valid flag",
			flagName: "test",
			expected: []string{"exists"},
		},
		{
			name:     "should not return value for invalid flag",
			flagName: "undefined",
			expected: []string{},
		},
	}
	existValue := []string{"exists"}
	flagsMultiStrReferences["test"] = &existValue
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := GetMultiStrFlag(test.flagName)
			require.Len(t, test.expected, len(got))
			require.Equal(t, test.expected, got)
		})
	}
}

func TestFlags_GetBoolFlag(t *testing.T) {
	tests := []struct {
		name     string
		flagName string
		expected bool
	}{
		{
			name:     "should return value for valid flag",
			flagName: "test",
			expected: true,
		},
		{
			name:     "should not return value for invalid flag",
			flagName: "undefined",
			expected: false,
		},
	}
	existValue := true
	flagsBoolReferences["test"] = &existValue
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := GetBoolFlag(test.flagName)
			require.Equal(t, test.expected, got)
		})
	}
}

func TestFlags_GetIntFlag(t *testing.T) {
	tests := []struct {
		name     string
		flagName string
		expected int
	}{
		{
			name:     "should return value for valid flag",
			flagName: "test",
			expected: 1,
		},
		{
			name:     "should not return value for invalid flag",
			flagName: "undefined",
			expected: -1,
		},
	}
	existValue := 1
	flagsIntReferences["test"] = &existValue
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := GetIntFlag(test.flagName)
			require.Equal(t, test.expected, got)
		})
	}
}

func TestFlags_SetStrFlag(t *testing.T) {
	tests := []struct {
		name     string
		flagName string
		expected string
	}{
		{
			name:     "should return value for valid flag",
			flagName: "test",
			expected: "exists",
		},
		{
			name:     "should not return value for invalid flag",
			flagName: "undefined",
			expected: "",
		},
	}
	existValue := "test"
	flagsStrReferences["test"] = &existValue
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			SetStrFlag(test.flagName, "exists")
			got := GetStrFlag(test.flagName)
			require.Equal(t, test.expected, got)
		})
	}
}

func TestFlags_SetMultiStrFlag(t *testing.T) {
	tests := []struct {
		name     string
		flagName string
		expected []string
	}{
		{
			name:     "should return value for valid flag",
			flagName: "test",
			expected: []string{"exists"},
		},
		{
			name:     "should not return value for invalid flag",
			flagName: "undefined",
			expected: []string{},
		},
	}
	existValue := []string{"test"}
	flagsMultiStrReferences["test"] = &existValue
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			SetMultiStrFlag(test.flagName, []string{"exists"})
			got := GetMultiStrFlag(test.flagName)
			require.Len(t, test.expected, len(got))
			require.Equal(t, test.expected, got)
		})
	}
}

func TestFlags_evalUsage(t *testing.T) {
	tests := []struct {
		name               string
		supportedPlatforms []string
		supportedProviders []string
		usage              string
		expected           string
	}{
		{
			name:               "should return same text",
			usage:              "test",
			expected:           "test",
			supportedPlatforms: []string{""},
			supportedProviders: []string{""},
		},
		{
			name:               "should return message translated",
			usage:              "test ${supportedPlatforms}",
			supportedPlatforms: []string{"terraform", "dockerfile"},
			supportedProviders: []string{"aws", "azure"},
			expected:           fmt.Sprintf("test %s", strings.Join([]string{"terraform", "dockerfile"}, ", ")),
		},
		{
			name:               "should return message translated for multiple variables",
			usage:              "test ${supportedPlatforms} ${defaultLogFile}",
			supportedPlatforms: []string{"terraform", "dockerfile"},
			supportedProviders: []string{"aws", "azure"},
			expected:           fmt.Sprintf("test %s %s", strings.Join([]string{"terraform", "dockerfile"}, ", "), constants.DefaultLogFile),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := evalUsage(
				test.usage,
				test.supportedPlatforms,
				test.supportedProviders)
			require.Equal(t, test.expected, got)
		})
	}
}
