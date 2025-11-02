package gosplit

import (
	"testing"
	"time"
)

func TestNormalizeTimeString(t *testing.T) {
	var tests = []struct {
		name  string
		input string
		want  string
	}{
		// the table itself
		{
			name:  "NormalizeTimeString",
			input: "00.000",
			want:  "0.000",
		},
		{
			name:  "NormalizeTimeString",
			input: "0",
			want:  "0.000",
		},
		{
			name:  "NormalizeTimeString",
			input: "0.0",
			want:  "0.000",
		},
		{
			name:  "NormalizeTimeString",
			input: "1.0",
			want:  "1.000",
		},
		{
			name:  "NormalizeTimeString",
			input: "1:0.0",
			want:  "1:00.000",
		},
		{
			name:  "NormalizeTimeString",
			input: "01:1.01",
			want:  "1:01.010",
		},
		{
			name:  "NormalizeTimeString",
			input: "04:3:2.10",
			want:  "4:03:02.100",
		},
		{
			name:  "NormalizeTimeString",
			input: "3:2:1.123",
			want:  "3:02:01.123",
		},
		{
			name:  "NormalizeTimeString",
			input: "30:20:10.123",
			want:  "30:20:10.123",
		},
		{
			name:  "NormalizeTimeString",
			input: "999:80:20.123123",
			want:  "999:80:20.123",
		},
	}
	for _, tt := range tests {
		Splits = []Segment{}
		t.Run(tt.name, func(t *testing.T) {
			ret := normalizeTimeString(tt.input)
			if ret != tt.want {
				t.Errorf("\nInput:\t%s\nExpected:\t%s\nGot:\t%s\n", tt.input, tt.want, ret)
			}
		})
	}
}

func TestTimeAsDuration(t *testing.T) {
	var tests = []struct {
		name  string
		input string
		want  time.Duration
	}{
		{
			name:  "TimeAsDuration",
			input: "0",
			want:  time.Duration(0),
		},
		{
			name:  "TimeAsDuration",
			input: "0.0",
			want:  time.Duration(0),
		},
		{
			name:  "TimeAsDuration",
			input: "1.0",
			want:  time.Duration(1) * time.Second,
		},
		{
			name:  "TimeAsDuration",
			input: "1:0.0",
			want:  time.Duration(1) * time.Minute,
		},
		{
			name:  "TimeAsDuration",
			input: "01:1.01",
			want: time.Duration(1)*time.Minute +
				time.Duration(1)*time.Second +
				time.Duration(10)*time.Millisecond,
		},
		{
			name:  "TimeAsDuration",
			input: "04:3:2.10",
			want: time.Duration(4)*time.Hour +
				time.Duration(3)*time.Minute +
				time.Duration(2)*time.Second +
				time.Duration(100)*time.Millisecond,
		},
		{
			name:  "TimeAsDuration",
			input: "3:2:1.123",
			want: time.Duration(3)*time.Hour +
				time.Duration(2)*time.Minute +
				time.Duration(1)*time.Second +
				time.Duration(123)*time.Millisecond,
		},
		{
			name:  "TimeAsDuration",
			input: "30:20:10.123",
			want: time.Duration(30)*time.Hour +
				time.Duration(20)*time.Minute +
				time.Duration(10)*time.Second +
				time.Duration(123)*time.Millisecond,
		},
		{
			name:  "TimeAsDuration",
			input: "999:80:20.123123",
			want: time.Duration(999)*time.Hour +
				time.Duration(80)*time.Minute +
				time.Duration(20)*time.Second +
				time.Duration(123)*time.Millisecond,
		},
	}
	for _, tt := range tests {
		Splits = []Segment{}
		t.Run(tt.name, func(t *testing.T) {
			ret := timeAsDuration(tt.input)
			if ret != tt.want {
				t.Errorf("\nInput:\t%s\nExpected:\t%s\nGot:\t%s\n", tt.input, tt.want, ret)
			}
		})
	}
}

func TestTimeAsString(t *testing.T) {
	var tests = []struct {
		name  string
		input time.Duration
		want  string
	}{
		{
			name:  "TimeAsString",
			input: time.Duration(0),
			want:  "0.000",
		},
		{
			name:  "TimeAsString",
			input: time.Duration(0),
			want:  "0.000",
		},
		{
			name:  "TimeAsString",
			input: time.Duration(1) * time.Second,
			want:  "1.000",
		},
		{
			name:  "TimeAsString",
			input: time.Duration(1) * time.Minute,
			want:  "1:00.000",
		},
		{
			name: "TimeAsString",
			input: time.Duration(1)*time.Minute +
				time.Duration(1)*time.Second +
				time.Duration(10)*time.Millisecond,
			want: "1:01.010",
		},
		{
			name: "TimeAsString",
			input: time.Duration(4)*time.Hour +
				time.Duration(3)*time.Minute +
				time.Duration(2)*time.Second +
				time.Duration(100)*time.Millisecond,
			want: "4:03:02.100",
		},
		{
			name: "TimeAsString",
			input: time.Duration(3)*time.Hour +
				time.Duration(2)*time.Minute +
				time.Duration(1)*time.Second +
				time.Duration(123)*time.Millisecond,
			want: "3:02:01.123",
		},
		{
			name: "TimeAsString",
			input: time.Duration(30)*time.Hour +
				time.Duration(20)*time.Minute +
				time.Duration(10)*time.Second +
				time.Duration(123)*time.Millisecond,
			want: "30:20:10.123",
		},
		{
			name: "TimeAsDuration",
			input: time.Duration(999)*time.Hour +
				time.Duration(80)*time.Minute +
				time.Duration(20)*time.Second +
				time.Duration(123)*time.Millisecond,
			want: "1000:20:20.123",
		},
	}
	for _, tt := range tests {
		Splits = []Segment{}
		t.Run(tt.name, func(t *testing.T) {
			ret := timeAsString(tt.input)
			if ret != tt.want {
				t.Errorf("\nInput:\t%s\nExpected:\t%s\nGot:\t%s\n", tt.input, tt.want, ret)
			}
		})
	}
}
