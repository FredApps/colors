package space

// IsValid reports whether s is a valid color space.
func (s Space) IsValid() bool {
	return s >= FirstSpace && s <= LastSpace
}

// ChannelCount returns the number of channels in the color space.
func (s Space) ChannelCount() int {
	if uint(s) >= uint(len(spaceChannelCounts)) {
		return 0
	}
	return int(spaceChannelCounts[s])
}

// HueIndex returns the index of the hue channel,
// or -1 if the color space has no hue channel.
func (s Space) HueIndex() int {
	if uint(s) >= uint(len(spaceHueIndexes)) {
		return -1
	}
	return int(spaceHueIndexes[s])
}

// CoordinateSystem returns the coordinate system of the color space.
func (s Space) CoordinateSystem() CoordinateSystem {
	if uint(s) >= uint(len(spaceCoordinateSystems)) {
		return 0
	}
	return spaceCoordinateSystems[s]
}

// Info returns information about the color space, or nil if s is invalid.
func (s Space) Info() SpaceInfo {
	if !s.IsValid() {
		return nil
	}
	return spaceInfos[s]
}
