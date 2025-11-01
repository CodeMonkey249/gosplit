package gosplit

import (
	"fmt"
	"testing"
)

func TestAddSplits(t *testing.T) {
	var tests = []struct {
		name  string
		input []string
		want  []Segment
	}{
		// the table itself
		{
			"SplitsWithAt",
			[]string{
				"sp add Segment1 SplitTime=2:00:3",
				"sp add Segment4 SplitTime=3:00:4.001",
				"sp add Segment5 SplitTime=4:00:00",
				"sp add Segment3 at=2 SplitTime=2:5:0",
				"sp add Segment2 at=2 SegmentTime=2:00.000",
			},
			[]Segment{
				{
					SegmentName: "Segment1",
					SplitTime:   "2:00:03.000",
					SegmentTime: "2:00:03.000",
					BestSegment: "2:00:03.000",
				},
				{
					SegmentName: "Segment2",
					SplitTime:   "2:02:03.000",
					SegmentTime: "2:00.000",
					BestSegment: "2:00.000",
				},
				{
					SegmentName: "Segment3",
					SplitTime:   "2:07:00.000",
					SegmentTime: "4:57.000",
					BestSegment: "4:57.000",
				},
				{
					SegmentName: "Segment4",
					SplitTime:   "3:02:04.001",
					SegmentTime: "55:04.001",
					BestSegment: "55:04.001",
				},
				{
					SegmentName: "Segment5",
					SplitTime:   "4:02:00.000",
					SegmentTime: "59:55.999",
					BestSegment: "59:55.999",
				},
			},
		},
		{
			"FullyKittedOutSplits",
			[]string{
				"sp add Segment1 SplitTime=2:00:3 SegmentTime=2:00:03.000 BestSegment=2:00.030",
				"sp add Segment2 SplitTime=2:02:03.000 SegmentTime=2:00.000 BestSegment=1:33.000",
				"sp add Segment3 SegmentTime=4:57.000 BestSegment=4:57.000 SplitTime=2:07:00.000",
				"sp add Segment4 SegmentTime=55:04.001 SplitTime=3:02:04.001",
				"sp add Segment5 SplitTime=4:02:00.000",
			},
			[]Segment{
				{
					SegmentName: "Segment1",
					SplitTime:   "2:00:03.000",
					SegmentTime: "2:00:03.000",
					BestSegment: "2:00.030",
				},
				{
					SegmentName: "Segment2",
					SplitTime:   "2:02:03.000",
					SegmentTime: "2:00.000",
					BestSegment: "1:33.000",
				},
				{
					SegmentName: "Segment3",
					SplitTime:   "2:07:00.000",
					SegmentTime: "4:57.000",
					BestSegment: "4:57.000",
				},
				{
					SegmentName: "Segment4",
					SplitTime:   "3:02:04.001",
					SegmentTime: "55:04.001",
					BestSegment: "55:04.001",
				},
				{
					SegmentName: "Segment5",
					SplitTime:   "4:02:00.000",
					SegmentTime: "59:55.999",
					BestSegment: "59:55.999",
				},
			},
		},
	}
	for _, tt := range tests {
		Splits = []Segment{}
		t.Run(tt.name, func(t *testing.T) {
			for _, cmd := range tt.input {
				_, _ = ParseCommands(cmd)
			}
			for i, seg := range tt.want {
				if seg.SegmentName != Splits[i].SegmentName ||
					seg.SplitTime != Splits[i].SplitTime ||
					seg.SegmentTime != Splits[i].SegmentTime ||
					seg.BestSegment != Splits[i].BestSegment {
					fmt.Println("Segments do not match expected:")
					fmt.Println(seg.SegmentName)
					fmt.Println("\nGot:")
					printSplits()
					fmt.Println("\nExpected:")
					Splits = tt.want
					printSplits()
					t.Errorf("Test Case Failed\n")
				}
			}
		})
	}
}
