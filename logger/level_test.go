package logger

import "testing"

func TestContains(t *testing.T) {
	testCases := [...]struct {
		name   string
		levels Level
		level  Level
		want   bool
	}{
		{
			name:   "LTrace",
			levels: LstdLevel,
			level:  LTrace,
			want:   true,
		},
		{
			name:   "LDebug",
			levels: LstdLevel,
			level:  LDebug,
			want:   true,
		},
		{
			name:   "LInfo",
			levels: LstdLevel,
			level:  LInfo,
			want:   true,
		},
		{
			name:   "LWarn",
			levels: LstdLevel,
			level:  LWarn,
			want:   true,
		},
		{
			name:   "LPanic",
			levels: LstdLevel,
			level:  LPanic,
			want:   true,
		},
		{
			name:   "LFatal",
			levels: LstdLevel,
			level:  LFatal,
			want:   true,
		},
		{
			name:   "LHTTP",
			levels: LstdLevel,
			level:  LHTTP,
			want:   true,
		},
		{
			name:   "false",
			levels: LError | LPanic | LFatal,
			level:  LDebug,
			want:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.levels.contains(tc.level); got != tc.want {
				t.Errorf(`want:%v, got:%v`, tc.want, got)
			}
		})
	}
}
