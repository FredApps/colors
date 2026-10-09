package codegen

import (
	"math"

	"github.com/thmalt/colors/gen/codegen/data"
	"github.com/thmalt/colors/gen/codegen/model"
)

var (
	HueInterpolation = [...]string{
		"HueShorter",
		"HueLonger",
		"HueIncreasing",
		"HueDecreasing",
	}

	xChannelName = channelName("X", "x", "X", "X")
	yChannelName = channelName("Y", "y", "Y", "Y")
	zChannelName = channelName("Z", "z", "Z", "Z")

	lightnessChannelName = channelName("Lightness", "l", "L", "Lightness")
	chromaChannelName    = channelName("Chroma", "c", "C", "Chroma")
	hueChannelName       = channelName("Hue", "h", "h", "Hue")

	uChannelName = channelName("U", "u", "u", "Green-Red Opponent")
	vChannelName = channelName("V", "v", "v", "Blue-Yellow Opponent")

	aChannelName = channelName("A", "a", "a", "Green-Red Opponent")
	bChannelName = channelName("B", "b", "b", "Blue-Yellow Opponent")

	rgbChannels = []model.Channel{
		extendChannel(channelName("Red", "r", "R", "Red"), numberChannel(0, 1, 6)),
		extendChannel(channelName("Green", "g", "G", "Green"), numberChannel(0, 1, 6)),
		extendChannel(channelName("Blue", "b", "B", "Blue"), numberChannel(0, 1, 6)),
	}

	xyzD50Channels = []model.Channel{
		extendChannel(xChannelName, numberChannel(0, data.D50Xyz[0], 8), unrestrictedChannel),
		extendChannel(yChannelName, numberChannel(0, data.D50Xyz[1], 8), unrestrictedChannel),
		extendChannel(zChannelName, numberChannel(0, data.D50Xyz[2], 8), unrestrictedChannel),
	}

	xyzD65Channels = []model.Channel{
		extendChannel(xChannelName, numberChannel(0, data.D65Xyz[0], 8), unrestrictedChannel),
		extendChannel(yChannelName, numberChannel(0, data.D65Xyz[1], 8), unrestrictedChannel),
		extendChannel(zChannelName, numberChannel(0, data.D65Xyz[2], 8), unrestrictedChannel),
	}

	xyzAbsD65Channels = []model.Channel{
		extendChannel(xChannelName, numberChannel(0, math.Round(data.D65Xyz[0]*data.XyzAbsD65Scale), 6), unrestrictedChannel),
		extendChannel(yChannelName, numberChannel(0, math.Round(data.D65Xyz[1]*data.XyzAbsD65Scale), 6), unrestrictedChannel),
		extendChannel(zChannelName, numberChannel(0, math.Round(data.D65Xyz[2]*data.XyzAbsD65Scale), 6), unrestrictedChannel),
	}

	xyYChannels = []model.Channel{
		extendChannel(channelName("Chromaticity x", "x", "x", "x"), numberChannel(0, 1, 8)),
		extendChannel(channelName("Chromaticity y", "y", "y", "y"), numberChannel(0, 1, 8)),
		extendChannel(channelName("Luminance", "luminance", "Y", "Y"), numberChannel(0, 1, 8), unrestrictedChannel),
	}

	labChannels = []model.Channel{
		extendChannel(lightnessChannelName, numberChannel(0, 100, 4)),
		extendChannel(aChannelName, numberChannel(-125, 125, 4), unrestrictedChannel),
		extendChannel(bChannelName, numberChannel(-125, 125, 4), unrestrictedChannel),
	}

	lchChannels = []model.Channel{
		extendChannel(lightnessChannelName, numberChannel(0, 100, 4)),
		extendChannel(chromaChannelName, numberChannel(0, 150, 4), unrestrictedChannel),
		extendChannel(hueChannelName, degreeChannel),
	}

	luvChannels = []model.Channel{
		extendChannel(lightnessChannelName, numberChannel(0, 100, 4)),
		extendChannel(uChannelName, numberChannel(-134, 220, 4), percentScaleChannel(215), unrestrictedChannel),
		extendChannel(vChannelName, numberChannel(-140, 122, 4), percentScaleChannel(215), unrestrictedChannel),
	}

	lchuvChannels = []model.Channel{
		extendChannel(lightnessChannelName, numberChannel(0, 100, 4)),
		extendChannel(chromaChannelName, numberChannel(0, 180, 4), percentScaleChannel(220), unrestrictedChannel),
		extendChannel(hueChannelName, degreeChannel),
	}

	oklabChannels = []model.Channel{
		extendChannel(lightnessChannelName, numberChannel(0, 1, 6)),
		extendChannel(aChannelName, numberChannel(-0.4, 0.4, 6), unrestrictedChannel),
		extendChannel(bChannelName, numberChannel(-0.4, 0.4, 6), unrestrictedChannel),
	}

	oklchChannels = []model.Channel{
		extendChannel(lightnessChannelName, numberChannel(0, 1, 6)),
		extendChannel(chromaChannelName, numberChannel(0, 0.4, 6), unrestrictedChannel),
		extendChannel(hueChannelName, degreeChannel),
	}

	hueChannel        = extendChannel(channelName("Hue", "h", "H", "Hue"), degreeChannel, precisionChannel(2))
	saturationChannel = extendChannel(channelName("Saturation", "s", "S", "Saturation"), percentChannel(0, 1, 4))
	lightnessChannel  = extendChannel(channelName("Lightness", "l", "L", "Lightness"), percentChannel(0, 1, 4))
	valueChannel      = extendChannel(channelName("Value", "v", "V", "Value"), percentChannel(0, 1, 4))
	whitenessChannel  = extendChannel(channelName("Whiteness", "w", "W", "Whiteness"), percentChannel(0, 1, 4))
	blacknessChannel  = extendChannel(channelName("Blackness", "b", "B", "Blackness"), percentChannel(0, 1, 4))

	hslChannels = []model.Channel{hueChannel, saturationChannel, lightnessChannel}
	hsvChannels = []model.Channel{hueChannel, saturationChannel, valueChannel}
	hwbChannels = []model.Channel{hueChannel, whitenessChannel, blacknessChannel}

	Spaces = [...]model.Space{
		{
			Name:        "Srgb",
			Family:      "RGB",
			Base:        "LinearSrgb",
			DisplayName: "sRGB",
			CssName:     "srgb",
			WhitePoint:  "D65",
			Channels:    rgbChannels,

			UseGenericColorFunction: true,
		},
		{
			Name:        "LinearSrgb",
			Family:      "RGB",
			Base:        "XyzD65",
			DisplayName: "Linear sRGB",
			CssName:     "srgb-linear",
			WhitePoint:  "D65",
			Channels:    rgbChannels,
			SnakeName:   "srgb_linear",

			UseGenericColorFunction: true,
		},
		{
			Name:        "DisplayP3",
			Family:      "RGB",
			Base:        "LinearDisplayP3",
			DisplayName: "Display P3",
			CssName:     "display-p3",
			WhitePoint:  "D65",
			Channels:    rgbChannels,

			UseGenericColorFunction: true,
		},
		{
			Name:        "LinearDisplayP3",
			Family:      "RGB",
			Base:        "XyzD65",
			DisplayName: "Linear Display P3",
			CssName:     "display-p3-linear",
			WhitePoint:  "D65",
			Channels:    rgbChannels,
			SnakeName:   "display_p3_linear",

			UseGenericColorFunction: true,
		},
		{
			Name:        "A98",
			Family:      "RGB",
			Base:        "LinearA98",
			DisplayName: "Adobe RGB (1998)",
			CssName:     "a98-rgb",
			WhitePoint:  "D65",
			Channels:    rgbChannels,

			UseGenericColorFunction: true,
		},
		{
			Name:        "LinearA98",
			Family:      "RGB",
			Base:        "XyzD65",
			DisplayName: "Linear Adobe RGB (1998)",
			CssName:     "a98-rgb-linear",
			WhitePoint:  "D65",
			Channels:    rgbChannels,
			SnakeName:   "a98_linear",

			UseGenericColorFunction: true,
		},
		{
			Name:        "ProPhoto",
			Family:      "RGB",
			Base:        "LinearProPhoto",
			DisplayName: "ProPhoto",
			CssName:     "prophoto-rgb",
			WhitePoint:  "D50",
			Channels:    rgbChannels,
			SnakeName:   "prophoto",

			UseGenericColorFunction: true,
		},
		{
			Name:        "LinearProPhoto",
			Family:      "RGB",
			Base:        "XyzD50",
			DisplayName: "Linear ProPhoto",
			CssName:     "prophoto-rgb-linear",
			WhitePoint:  "D50",
			Channels:    rgbChannels,
			SnakeName:   "prophoto_linear",

			UseGenericColorFunction: true,
		},
		{
			Name:        "Rec2020",
			Family:      "RGB",
			Base:        "LinearRec2020",
			DisplayName: "Rec. 2020",
			CssName:     "rec2020",
			WhitePoint:  "D65",
			Channels:    rgbChannels,

			UseGenericColorFunction: true,
		},
		{
			Name:        "Rec2020OETF",
			Family:      "RGB",
			Base:        "LinearRec2020",
			DisplayName: "Rec. 2020 Scene Referred",
			CssName:     "rec2020-oetf",
			WhitePoint:  "D65",
			Channels:    rgbChannels,
			SnakeName:   "rec2020_oetf",

			UseGenericColorFunction: true,
		},
		{
			Name:        "LinearRec2020",
			Family:      "RGB",
			Base:        "XyzD65",
			DisplayName: "Linear Rec. 2020",
			CssName:     "rec2020-linear",
			WhitePoint:  "D65",
			Channels:    rgbChannels,
			SnakeName:   "rec2020_linear",

			UseGenericColorFunction: true,
		},
		{
			Name:        "LinearRec2100",
			Equivalent:  "LinearRec2020",
			Family:      "RGB",
			Base:        "XyzD65",
			DisplayName: "Linear Rec. 2100",
			CssName:     "rec2100-linear",
			WhitePoint:  "D65",
			Channels:    rgbChannels,
			SnakeName:   "rec2100_linear",

			UseGenericColorFunction: true,
		},
		{
			Name:        "Rec2100PQ",
			Family:      "RGB",
			Base:        "XyzD65",
			DisplayName: "Rec. 2100 PQ",
			CssName:     "rec2100-pq",
			WhitePoint:  "D65",
			Channels:    rgbChannels,
			SnakeName:   "rec2100_pq",

			UseGenericColorFunction: true,
		},
		{
			Name:        "Rec2100HLG",
			Family:      "RGB",
			Base:        "XyzD65",
			DisplayName: "Rec. 2100 HLG",
			CssName:     "rec2100-hlg",
			WhitePoint:  "D65",
			Channels:    rgbChannels,
			SnakeName:   "rec2100_hlg",

			UseGenericColorFunction: true,
		},
		{
			Name:        "XyzD50",
			Family:      "XYZ",
			Base:        "XyzD65",
			DisplayName: "CIE XYZ D50",
			CssName:     "xyz-d50",
			WhitePoint:  "D50",
			Channels:    xyzD50Channels,

			UseGenericColorFunction: true,
		},
		{
			Name:        "XyzD65",
			Aliases:     []string{"Xyz"},
			Family:      "XYZ",
			DisplayName: "CIE XYZ D65",
			CssName:     "xyz-d65",
			CssAliases:  []string{"xyz"},
			WhitePoint:  "D65",
			Channels:    xyzD65Channels,

			UseGenericColorFunction: true,
		},
		{
			Name:        "XyzAbsD65",
			Family:      "XYZ",
			DisplayName: "Absolute XYZ D65",
			CssName:     "xyz-abs-d65",
			WhitePoint:  "D65",
			Channels:    xyzAbsD65Channels,
			SnakeName:   "xyz_abs_d65",

			UseGenericColorFunction: true,
		},
		{
			Name:        "XyYD50",
			Family:      "XYZ",
			Base:        "XyzD50",
			DisplayName: "CIE xyY",
			CssName:     "xyy-d50",
			WhitePoint:  "D50",
			Channels:    xyYChannels,
			SnakeName:   "xyy_d50",

			UseGenericColorFunction: true,
		},
		{
			Name:        "XyYD65",
			Aliases:     []string{"XyY"},
			Family:      "XYZ",
			Base:        "XyzD65",
			DisplayName: "CIE xyY",
			CssName:     "xyy-d65",
			CssAliases:  []string{"xyy"},
			WhitePoint:  "D65",
			Channels:    xyYChannels,
			SnakeName:   "xyy_d65",

			UseGenericColorFunction: true,
		},
		{
			Name:        "LabD50",
			Aliases:     []string{"Lab"},
			Family:      "Lab",
			Base:        "XyzD50",
			DisplayName: "CIE Lab D50",
			CssName:     "lab",
			CssAliases:  []string{"lab-d50"},
			WhitePoint:  "D50",
			Channels:    labChannels,
			ParseKind:   model.ParseLabLike,
		},
		{
			Name:        "LchD50",
			Aliases:     []string{"Lch"},
			Family:      "Lab",
			Base:        "LabD50",
			DisplayName: "CIE LCh D50",
			CssName:     "lch",
			CssAliases:  []string{"lch-d50"},
			WhitePoint:  "D50",
			Coordinate:  model.Polar,
			Channels:    lchChannels,
			ParseKind:   model.ParseLchLike,
		},
		{
			Name:        "LabD65",
			Family:      "Lab",
			Base:        "XyzD65",
			DisplayName: "CIE Lab D65",
			CssName:     "lab-d65",
			WhitePoint:  "D65",
			Channels:    labChannels,
			ParseKind:   model.ParseLabLike,

			UseGenericColorFunction: true,
		},
		{
			Name:        "LchD65",
			Family:      "Lab",
			Base:        "LabD65",
			DisplayName: "CIE LCh D65",
			CssName:     "lch-d65",
			WhitePoint:  "D65",
			Coordinate:  model.Polar,
			Channels:    lchChannels,
			ParseKind:   model.ParseLchLike,

			UseGenericColorFunction: true,
		},
		{
			Name:        "LuvD50",
			Aliases:     []string{"Luv"},
			Family:      "Luv",
			Base:        "XyzD50",
			DisplayName: "CIE Luv D50",
			CssName:     "luv",
			CssAliases:  []string{"luv-d50"},
			WhitePoint:  "D50",
			Coordinate:  model.Polar,
			Channels:    luvChannels,
			ParseKind:   model.ParseLabLike,
		},
		{
			Name:        "LchuvD50",
			Aliases:     []string{"Lchuv"},
			Family:      "Luv",
			Base:        "LuvD50",
			DisplayName: "CIE LChuv D50",
			CssName:     "lchuv",
			CssAliases:  []string{"lchuv-d50"},
			WhitePoint:  "D50",
			Coordinate:  model.Polar,
			Channels:    lchuvChannels,
			ParseKind:   model.ParseLchLike,
		},
		{
			Name:        "LuvD65",
			Family:      "Luv",
			Base:        "XyzD65",
			DisplayName: "CIE Luv D65",
			CssName:     "luv-d65",
			WhitePoint:  "D65",
			Coordinate:  model.Polar,
			Channels:    luvChannels,
			ParseKind:   model.ParseLabLike,

			UseGenericColorFunction: true,
		},
		{
			Name:        "LchuvD65",
			Family:      "Luv",
			Base:        "LuvD65",
			DisplayName: "CIE LChuv D65",
			CssName:     "lchuv-d65",
			WhitePoint:  "D65",
			Coordinate:  model.Polar,
			Channels:    lchuvChannels,
			ParseKind:   model.ParseLchLike,

			UseGenericColorFunction: true,
		},
		{
			Name:        "Oklab",
			Family:      "Oklab",
			Base:        "XyzD65",
			DisplayName: "Oklab",
			CssName:     "oklab",
			WhitePoint:  "D65",
			Channels:    oklabChannels,
			ParseKind:   model.ParseLabLike,
		},
		{
			Name:        "Oklch",
			Family:      "Oklab",
			Base:        "Oklab",
			DisplayName: "Oklch",
			CssName:     "oklch",
			WhitePoint:  "D65",
			Coordinate:  model.Polar,
			Channels:    oklchChannels,
			ParseKind:   model.ParseLchLike,
		},
		{
			Name:        "Hsl",
			Family:      "RGB",
			Base:        "Srgb",
			DisplayName: "HSL",
			CssName:     "hsl",
			WhitePoint:  "D65",
			Coordinate:  model.Polar,
			Channels:    hslChannels,
		},
		{
			Name:        "Hsv",
			Family:      "RGB",
			Base:        "Srgb",
			DisplayName: "HSV",
			CssName:     "hsv",
			WhitePoint:  "D65",
			Coordinate:  model.Polar,
			Channels:    hsvChannels,
		},
		{
			Name:        "Hwb",
			Family:      "RGB",
			Base:        "Srgb",
			DisplayName: "HWB",
			CssName:     "hwb",
			WhitePoint:  "D65",
			Coordinate:  model.Polar,
			Channels:    hwbChannels,
		},
	}

	WhitePoints = [...]model.WhitePoint{
		{
			Name: "D50",
			X:    data.D50.X,
			Y:    data.D50.Y,
		},
		{
			Name: "D65",
			X:    data.D65.X,
			Y:    data.D65.Y,
		},
	}

	ConvertFuncs = [...]ConvertFunc{
		// Standard transfer conversion
		transferFunc("Srgb", "LinearSrgb", "SrgbDecode"),
		transferFunc("LinearSrgb", "Srgb", "SrgbEncode"),

		transferFunc("A98", "LinearA98", "A98Decode"),
		transferFunc("LinearA98", "A98", "A98Encode"),

		transferFunc("DisplayP3", "LinearDisplayP3", "SrgbDecode"),
		transferFunc("LinearDisplayP3", "DisplayP3", "SrgbEncode"),

		transferFunc("ProPhoto", "LinearProPhoto", "ProPhotoDecode"),
		transferFunc("LinearProPhoto", "ProPhoto", "ProPhotoEncode"),

		transferFunc("Rec2020", "LinearRec2020", "Rec2020Decode"),
		transferFunc("LinearRec2020", "Rec2020", "Rec2020Encode"),

		transferFunc("Rec2020OETF", "LinearRec2020", "Rec2020OETFDecode"),
		transferFunc("LinearRec2020", "Rec2020OETF", "Rec2020OETFEncode"),

		transferFunc("Rec2100HLG", "LinearRec2100", "Rec2100HLGDecode"),
		transferFunc("LinearRec2100", "Rec2100HLG", "Rec2100HLGEncode"),

		transferFunc("Rec2100PQ", "LinearRec2100", "Rec2100PQDecode"),
		transferFunc("LinearRec2100", "Rec2100PQ", "Rec2100PQEncode"),

		// Equivalent
		convertFunc("LinearRec2100", "LinearRec2020"),
		convertFunc("LinearRec2020", "LinearRec2100"),

		// Standard conversion
		implementedFunc("Srgb", "Hsl"),
		implementedFunc("Hsl", "Srgb"),

		implementedFunc("Srgb", "Hsv"),
		implementedFunc("Hsv", "Srgb"),

		implementedFunc("Srgb", "Hwb"),
		implementedFunc("Hwb", "Srgb"),

		implementedFunc("Hsl", "Hsv"),
		implementedFunc("Hsv", "Hsl"),

		implementedFunc("Hsl", "Hwb"),
		implementedFunc("Hwb", "Hsl"),

		implementedFunc("Hsv", "Hwb"),
		implementedFunc("Hwb", "Hsv"),

		// Standard conversion
		implementedFunc("LabD50", "XyzD50"),
		implementedFunc("XyzD50", "LabD50"),

		implementedFunc("LabD65", "XyzD65"),
		implementedFunc("XyzD65", "LabD65"),

		implementedFunc("LuvD50", "XyzD50"),
		implementedFunc("XyzD50", "LuvD50"),

		implementedFunc("LuvD65", "XyzD65"),
		implementedFunc("XyzD65", "LuvD65"),

		implementedFunc("XyzAbsD65", "XyzD65"),
		implementedFunc("XyzD65", "XyzAbsD65"),

		implementedFunc("XyY", "Xyz"),
		implementedFunc("Xyz", "XyY"),

		// Generate with Call Ops
		convertFunc("XyYD50", "XyzD50", opCall("XyY", "Xyz")),
		convertFunc("XyzD50", "XyYD50", opCall("Xyz", "XyY")),

		convertFunc("XyYD65", "XyzD65", opCall("XyY", "Xyz")),
		convertFunc("XyzD65", "XyYD65", opCall("Xyz", "XyY")),

		convertFunc("LabD50", "LchD50", opCall("lab", "Lch")),
		convertFunc("LchD50", "LabD50", opCall("Lch", "Lxy")),

		convertFunc("LabD65", "LchD65", opCall("lab", "Lch")),
		convertFunc("LchD65", "LabD65", opCall("Lch", "Lxy")),

		convertFunc("LuvD50", "LchuvD50", opCall("luv", "Lchuv")),
		convertFunc("LchuvD50", "LuvD50", opCall("Lch", "Lxy")),

		convertFunc("LuvD65", "LchuvD65", opCall("luv", "Lchuv")),
		convertFunc("LchuvD65", "LuvD65", opCall("Lch", "Lxy")),

		convertFunc("Oklab", "Oklch", opCall("oklab", "Oklch")),
		convertFunc("Oklch", "Oklab", opCall("Lch", "Lxy")),

		// Generate with Matrix Ops
		// Oklab
		convertFunc("Oklab", "XyzD65", opMatrix(data.OklabLabToLms), opCube(), opMatrix(data.OklabLmsToXyzD65)),
		convertFunc("XyzD65", "Oklab", opMatrix(data.OklabXyzD65ToLms), opCbrt(), opMatrix(data.OklabLmsToLab)),

		// Xyz*
		convertFunc("XyzD65", "XyzD50", opMatrix(data.XyzD65ToXyzD50)),
		convertFunc("XyzD50", "XyzD65", opMatrix(data.XyzD50ToXyzD65)),

		// Xyz* -> Linear*
		convertFunc("XyzD65", "LinearSrgb", opMatrix(data.XyzD65ToLinearSrgb)),
		convertFunc("LinearSrgb", "XyzD65", opMatrix(data.LinearSrgbToXyzD65)),

		convertFunc("XyzD65", "LinearDisplayP3", opMatrix(data.XyzD65ToLinearDisplayP3)),
		convertFunc("LinearDisplayP3", "XyzD65", opMatrix(data.LinearDisplayP3ToXyzD65)),

		convertFunc("XyzD65", "LinearA98", opMatrix(data.XyzD65ToLinearA98)),
		convertFunc("LinearA98", "XyzD65", opMatrix(data.LinearA98ToXyzD65)),

		convertFunc("XyzD50", "LinearProPhoto", opMatrix(data.XyzD50ToLinearProPhoto)),
		convertFunc("LinearProPhoto", "XyzD50", opMatrix(data.LinearProPhotoToXyzD50)),

		convertFunc("XyzD65", "LinearRec2020", opMatrix(data.XyzD65ToLinearRec2020)),
		convertFunc("LinearRec2020", "XyzD65", opMatrix(data.LinearRec2020ToXyzD65)),
	}
)
