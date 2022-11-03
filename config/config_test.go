package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFromEnv(t *testing.T) {
	type args struct {
		setEnvFunc func()
	}
	cases := []struct {
		name     string
		args     args
		expected Config
	}{
		{
			name: "env is filled",
			args: args{
				setEnvFunc: func() {
					err := os.Setenv("BINN_SEND_INTERVAL_SEC", "30")
					require.NoError(t, err)
				},
			},
			expected: Config{
				SendInterval: 30 * time.Second,
			},
		},
		{
			name: "env is not filled",
			args: args{
				setEnvFunc: func() {
					err := os.Setenv("BINN_SEND_INTERVAL_SEC", "")
					require.NoError(t, err)
				},
			},
			expected: Config{
				SendInterval: 10 * time.Second,
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.setEnvFunc()
			c := NewFromEnv()
			assert.Equal(t, tt.expected, c)
		})
	}
}
