package gosplit

import (
	"fmt"
	"time"
)

type Segment struct {
	SegmentName string
	SplitTime string
	SegmentTime string
	BestSegment string
}


func setSplitTime(t string, i int) error {
	if i >= len(Splits) {
		return fmt.Errorf("Index out of bounds. Cannot edit split %d with current number of splits %d", i, len(Splits))
	}

	Splits[i].SplitTime = t

	for j, split := range Splits {
		if j != 0 {
			currSplitTime := timeAsDuration(Splits[j].SplitTime)
			prevSplitTime := timeAsDuration(Splits[j-1].SplitTime)

			if currSplitTime < prevSplitTime {
				currSplitTime = prevSplitTime
			}
			Splits[j].SplitTime = timeAsString(currSplitTime)

			currSegmentTime := currSplitTime - prevSplitTime
			Splits[j].SegmentTime = timeAsString(currSegmentTime)
		} else {
			Splits[j].SegmentTime = split.SplitTime
		}
	}

	updateBestSegment()

	return nil
}

func setSegmentTime(t string, i int) error {
	if i >= len(Splits) {
		return fmt.Errorf("Index out of bounds. Cannot edit split %d with current number of splits %d", i, len(Splits))
	}

	Splits[i].SegmentTime = t

	segmentTimeSum := time.Duration(0)
	for j, split := range Splits {
		segTime := timeAsDuration(split.SegmentTime)
		segmentTimeSum += segTime
		Splits[j].SplitTime = timeAsString(segmentTimeSum)
	}

	updateBestSegment()

	return nil
}

func setBestSegment(t string, i int) {
	Splits[i].BestSegment = t
	updateBestSegment()
}

func updateBestSegment() {
	for j, _ := range Splits {
		bestSeg := timeAsDuration(Splits[j].BestSegment)
		segTime := timeAsDuration(Splits[j].SegmentTime)
		if bestSeg > segTime {
			Splits[j].BestSegment = Splits[j].SegmentTime
		}
	}
}

