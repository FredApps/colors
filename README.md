# colors

[![Go Reference](https://pkg.go.dev/badge/github.com/thmalt/colors.svg)](https://pkg.go.dev/github.com/thmalt/colors)
![Go Version](https://img.shields.io/github/go-mod/go-version/thmalt/colors?color=506070)
[![License: MIT](https://img.shields.io/badge/License-MIT-506070.svg)](LICENSE)

A Go library for color space conversions, interpolation, and gradients.

## Installation

```sh
go get github.com/thmalt/colors
```

## Usage

```go
import "github.com/thmalt/colors"
```

### Example

```go
package main

import (
	"fmt"

	"github.com/thmalt/colors"
	"github.com/thmalt/colors/space"
)

func main() {
	c := colors.Rgba(255, 0, 0, .5)

	fmt.Println(c.Hex())
	fmt.Println(c.To(space.Oklab))
	fmt.Println(c.Oklab())
}
```

## Color Conversion

The `github.com/thmalt/colors/convert` package provides direct conversion functions between color spaces.

Direct conversion functions are generated for the supported color-space pairs.
The `convert` package provides low-overhead conversion APIs, but it is
**not necessarily the fastest possible implementation**.

```go
l, a, b := convert.SrgbToOklab(r, g, b)
```

The root `github.com/thmalt/colors` package builds on `github.com/thmalt/colors/convert` and provides two conversion modes.

### `colors` build modes

#### Default build

The default build uses a hub-based conversion strategy to reduce binary size:

```bash
go build
```

Cross-family conversions may be routed through the `XyzD50` or `XyzD65` conversion hubs,
while conversions within the same color-space family remain direct.

#### Full conversion build

For applications where conversion performance is more important than binary size,
the full direct conversion graph can be enabled with the `colors_full` build tag:

```bash
go build -tags=colors_full
```

This generates direct conversion paths between color spaces instead of routing
cross-family conversions through the XYZ hubs.

| Package / Build tag            | Conversion strategy | Binary size | Performance |
| ------------------------------ | ------------------- | ----------- | ----------- |
| `colors/convert`               | Direct conversions  | On demand   | Higher      |
| `colors` (`-tags=colors_full`) | Direct paths        | Larger      | Medium      |
| `colors` (default)             | Hub-based paths     | Smaller     | Lower       |

### Binary size

Conversion code is included in the final binary based on how the conversion API is used.

Specific methods such as `Color.Srgb()`, `Color.Oklab()`, and `Color.XyzD50()`
only reference conversions required for their respective color spaces.
Unused conversions can be removed by the Go linker.

Generic methods such as `Color.To(dst space.Space)` and `Color.MustTo(dst space.Space)`
reference the complete generated conversion dispatch for the selected build mode.

For smaller production binaries, use `-s` and `-w` to strip symbol and DWARF information.
`-trimpath` can also be used to remove local filesystem paths:

```bash
go build -ldflags="-s -w" -trimpath
```

## Examples

For complete, runnable examples, see the [`examples`](./examples) directory.

### Color and Color Spaces

```go
package main

import (
	"fmt"

	"github.com/thmalt/colors"
	"github.com/thmalt/colors/space"
)

func main() {
	c := colors.Hex("#323c46fe")

	// Convert a color to different color spaces.
	fmt.Println("sRGB: ", c)
	fmt.Println("Oklab:", c.To(space.Oklab))
	fmt.Println("Oklch:", c.To(space.Oklch))
	fmt.Println("Lab:  ", c.To(space.Lab))
	fmt.Println("LCh:  ", c.To(space.Lch))

	fmt.Println()

	// Serialize the color as hexadecimal.
	fmt.Println("Hex:           ", c.WithAlpha(1).Hex())
	fmt.Println("Hex with alpha:", c.Hex())

	fmt.Println()

	// Inspect color space metadata.
	info := c.Space().Info()
	if info == nil {
		fmt.Println("invalid color space:", c.Space())
		return
	}

	fmt.Println("Color space:", info.DisplayName())

	whitePoint := info.WhitePoint()
	fmt.Println("White point:", whitePoint.Name)
	fmt.Println("  X:", whitePoint.X)
	fmt.Println("  Y:", whitePoint.Y)
	fmt.Println("  Z:", whitePoint.Z)

	fmt.Println("Channels:", info.ChannelCount())
	for i, ch := range info.Channels() {
		fmt.Printf(
			"  %d: %s [%g, %g], unrestricted=%t\n",
			i, ch.Name, ch.Min, ch.Max, ch.Unrestricted,
		)
	}
}
```

### Gradient Examples

```go
blue   = colors.Oklab(0.42, -0.02, -0.15)
purple = colors.Oklab(0.52, 0.08, -0.14)
orange = colors.Oklab(0.72, 0.12, 0.12)

// CSS:
// linear-gradient(
//   to bottom right in oklab,
//   #1b469e,
//   #504baa,
//   #784db6,
//   #b56c8a,
//   #f87c3e
// )
linearGradient = colors.NewGradientBuilder().
	AddStop(blue).
	AddStop(blue.Mix(purple, 0.5), 0.25).
	AddStop(purple, 0.5).
	AddStop(purple.Mix(orange, 0.5), 0.75).
	AddStop(orange, 1).
	Build()

linearSpec := gradient.NewLinearSpec()
linearSpec.SetSize(width, height)
linearSpec.SetDirection(gradient.ToBottomRight)

linear, err := linearSpec.Build()
```

<p align="center">
 <img src="./examples/gradient/output/linear.png" alt="Linear">
</p>

```go
// CSS:
// conic-gradient(
//   from 0.25turn at 50% 30% in oklab,
//   #f69d3c,
//   10deg,
//   #3f87a6,
//   350deg,
//   #ebf8e1
// )
conicGradient = colors.NewGradientBuilder().
	AddStop(colors.Hex("#f69d3c")).
	AddHint(gradient.DegToTurn(10)).
	AddStop(colors.Hex("#3f87a6")).
	AddHint(gradient.DegToTurn(350)).
	AddStop(colors.Hex("#ebf8e1")).
	Build()

conicSpec := gradient.NewConicSpec()
conicSpec.SetSize(width, height)
conicSpec.SetCenter(0.5, 0.3)
conicSpec.SetStartAngle(gradient.ToRight)

conic, err := conicSpec.Build()
```

<p align="center">
 <img src="./examples/gradient/output/conic.png" alt="Conic">
</p>

```go
// CSS:
// radial-gradient(
//   farthest-corner in oklab,
//   #3f5efb,
//   #00ff88,
//   #ce436d,
//   #0413bf 50%,
//   #fc466b
// )
radialGradient = colors.NewGradientBuilder().
	AddStop(colors.Hex("#3f5efb")).
	AddStop(colors.Hex("#00ff88")).
	AddStop(colors.Hex("#ce436d")).
	AddStop(colors.Hex("#0413bf"), 0.5).
	AddStop(colors.Hex("#fc466b")).
	Build()

radialSpec := gradient.NewRadialSpec()
radialSpec.SetSize(width, height)

radialSpec.SetShape(gradient.RadialCircle)
circle, err := radialSpec.Build()

radialSpec.SetShape(gradient.RadialEllipse)
ellipse, err := radialSpec.Build()
```

<p align="center">
 <img src="./examples/gradient/output/radial-circle.png" alt="Radial circle">
 <img src="./examples/gradient/output/radial-ellipse.png" alt="Radial ellipse">
</p>

More examples are available in the [`examples`](./examples) directory.

## Development

Run the package tests and vet checks with the default build and each LUT build tag:

```sh
go test ./...
go vet ./...
for tag in colors_full static st; do
    go test -tags="$tag" ./...
    go vet -tags="$tag" ./...
done
```

Regenerate the source with `go generate ./gen`. To also regenerate the lookup
tables, run `GEN_LUT=1 go generate ./gen` (or `sh gen.sh all`). The generator's
internal sRGB bootstrap tables use the same runtime-initialization template as
the public `convert` package, but always build without tag constraints so the
generator can compile with `colors_full`, `static`, or `st` enabled.

AI tools are used for research, problem solving, including advanced topics in color science and rendering.
