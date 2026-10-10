// Copyright (C) 2008 Apple Inc.  All rights reserved.
// Copyright (C) 2015 Google Inc. All rights reserved.
//
// Portions of this file are based on Blink's CSSGradientValue
// implementation.
//
// Original source:
// third_party/blink/renderer/core/css/css_gradient_value.cc
//
// See third_party/LICENSES, section "Blink - CSSGradientValue".

package colors

import (
	"math"

	"github.com/thmalt/colors/space"
)

func resolveStops(stops []GradientStop, mixer Mixer) []GradientStop {
	// Stops may come from another gradient. Recompute all private metadata
	// from their public colors and offsets with this gradient's mixer.
	for i := range stops {
		stops[i] = GradientStop{Color: stops[i].Color, Offset: stops[i].Offset}
	}
	resolveStopOffsets(stops)
	convertStopColors(stops, mixer.Space())
	resolved := resolveHints(stops, mixer)
	precomputeInvRange(resolved)
	return resolved
}

func convertStopColors(stops []GradientStop, dst space.Space) {
	for i := range stops {
		if stops[i].IsHint() {
			continue
		}

		stops[i].Color.MutTo(dst)
	}
}

func precomputeInvRange(stops []GradientStop) {
	for i := 1; i < len(stops); i++ {
		d := stops[i].Offset - stops[i-1].Offset
		if d == 0 {
			stops[i-1].invRange = 0
		} else {
			stops[i-1].invRange = 1 / d
		}
	}
}

func resolveStopOffsets(stops []GradientStop) {
	count := len(stops)
	if count == 0 {
		return
	}

	// Resolve the first and last offsets.
	if !stops[0].HasOffset() {
		stops[0].Offset = 0
	}
	if count > 1 && !stops[count-1].HasOffset() {
		stops[count-1].Offset = 1
	}

	// Ensure offsets are non-decreasing.
	prev := 0
	for i := range stops {
		if !stops[i].HasOffset() {
			continue
		}

		if stops[i].Offset < stops[prev].Offset {
			stops[i].Offset = stops[prev].Offset
		}

		prev = i
	}

	// Evenly distribute runs of unspecified offsets.
	if count > 2 {
		start := -1
		for i := range count {
			if !stops[i].HasOffset() {
				if start < 0 {
					start = i
				}
				continue
			}

			if start < 0 {
				continue
			}

			beg := stops[start-1].Offset
			end := stops[i].Offset
			step := (end - beg) / float64(i-start+1)
			for j := start; j < i; j++ {
				stops[j].Offset = beg + float64(j-start+1)*step
			}

			start = -1
		}
	}
}

func resolveHints(stops []GradientStop, mixer Mixer) []GradientStop {
	if len(stops) < 2 {
		return stops
	}

	resolved := make([]GradientStop, 0, len(stops))
	resolved = append(resolved, stops[0])

	last := len(stops) - 1
	for i := 1; i < last; i++ {
		if !stops[i].IsHint() {
			resolved = append(resolved, stops[i])
			continue
		}

		if stops[i-1].IsHint() || stops[i+1].IsHint() {
			continue
		}

		resolved = resolveHint(resolved, stops[i-1], stops[i], stops[i+1], mixer)
	}

	resolved = append(resolved, stops[last])
	return resolved
}

// resolveHints resolves color hints according to the CSS gradient
// color-hint resolution algorithm.
//
// The algorithm is based on WebKit/Blink's CSSGradientValue implementation.
// See third_party/LICENSES, section "WebKit/Blink - CSSGradientValue".
func resolveHint(dst []GradientStop, left, hint, right GradientStop, mixer Mixer) []GradientStop {
	offsetLeft := left.Offset
	offset := hint.Offset
	offsetRight := right.Offset

	leftDist := offset - offsetLeft
	rightDist := offsetRight - offset

	if nearlyEqual(leftDist, rightDist) {
		return dst
	}

	segment := &gradientHint{left: left.Color, right: right.Color}
	hint.hint = segment

	if nearlyEqual(leftDist, 0.0) {
		hint.Color = right.Color
		hint.weight = 1
		return append(dst, hint)
	}

	dst[len(dst)-1].hint = segment
	dst[len(dst)-1].weight = 0

	if nearlyEqual(rightDist, 0.0) {
		hint.Color = left.Color
		return append(dst, hint)
	}

	var stops [9]GradientStop

	if leftDist > rightDist {
		for i := range 7 {
			stops[i].Offset = offsetLeft + leftDist*(float64(7+i)/13.0)
		}

		stops[7].Offset = offset + rightDist*(1.0/3.0)
		stops[8].Offset = offset + rightDist*(2.0/3.0)
	} else {
		stops[0].Offset = offsetLeft + leftDist*(1.0/3.0)
		stops[1].Offset = offsetLeft + leftDist*(2.0/3.0)

		for i := range 7 {
			stops[i+2].Offset = offset + rightDist*(float64(i)/13.0)
		}

	}

	totalDist := offsetRight - offsetLeft
	hintRelativeOffset := leftDist / totalDist
	exponent := math.Log(0.5) / math.Log(hintRelativeOffset)

	for i := range stops {
		t := (stops[i].Offset - offsetLeft) / totalDist
		weighting := math.Pow(t, exponent)

		if math.IsNaN(weighting) || math.IsInf(weighting, 0) {
			continue
		}

		stops[i].Color = mixer.Mix(left.Color, right.Color, weighting)
		stops[i].hint = segment
		stops[i].weight = weighting

		dst = append(dst, stops[i])
	}

	return dst
}

const nearlyEqualTolerance = 1.0 / 4096

func nearlyEqual(a, b float64) bool {
	return math.Abs(a-b) <= nearlyEqualTolerance
}
