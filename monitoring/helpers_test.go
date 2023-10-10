package libpack_monitoring

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func (suite *MonitoringTestSuite) Test_validate_metrics_name() {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test_validate_metrics_name - valid",
			args: args{
				name: "test_metrics_name",
			},
			wantErr: false,
		},
		{
			name: "Test_validate_metrics_name - invalid",
			args: args{
				name: "test metrics name",
			},
			wantErr: true,
		},
		{
			name: "Test_validate_metrics_name - invalid - special chars",
			args: args{
				name: "test_metrics_name!",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			err := validate_metrics_name(tt.args.name)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func (suite *MonitoringTestSuite) TestValidateMetricsName() {
	tests := []struct {
		name     string
		expected error
	}{
		{
			name:     "valid_name",
			expected: nil,
		},
		{
			name:     "name with spaces",
			expected: fmt.Errorf("Invalid metric name: %s, expected %s", "name with spaces", "name_with_spaces"),
		},
		{
			name:     "name with non-alphanumeric characters",
			expected: fmt.Errorf("Invalid metric name: %s, expected %s", "name with non-alphanumeric characters", "name_with_nonalphanumeric_characters"),
		},
		{
			name:     "name__with__consecutive__underscores",
			expected: fmt.Errorf("Invalid metric name: %s, expected %s", "name__with__consecutive__underscores", "name_with_consecutive_underscores"),
		},
		{
			name:     "_name_with_leading_or_trailing_underscores_",
			expected: fmt.Errorf("Invalid metric name: %s, expected %s", "_name_with_leading_or_trailing_underscores_", "name_with_leading_or_trailing_underscores"),
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			err := validate_metrics_name(tt.name)
			assert.Equal(t, tt.expected, err)
		})
	}
}
