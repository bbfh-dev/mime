# 📀 Vintage

Minecraft data-driven vanilla data & resource pack development kit with minimum boilerplate and abstraction.

> [!IMPORTANT]
> Check out the [Wiki](https://codeberg.org/bbfh/vintage/wiki) for information.

## Rationality

There are many [other pre-processors](https://gist.github.com/Ellivers/db296c438f9f87bbf9c79d24f940fe03), but all of them are fundamentally different from what Vintage is trying to accomplish:

1. Minimum code abstraction.[⁽¹⁾](https://codeberg.org/bbfh/vintage/wiki/Philosophy#1-minimum-code-abstraction)
1. No [vendor lock-in](https://en.wikipedia.org/wiki/Vendor_lock-in) to specific language.[⁽²⁾](https://codeberg.org/bbfh/vintage/wiki/Philosophy#2-no-vendor-lock-in-to-specific-language)
1. Power thanks to simplicity not over-engineering.[⁽³⁾](https://codeberg.org/bbfh/vintage/wiki/Philosophy#3-power-thanks-to-simplicity-not-over-engineering)
1. User-control and predictability.[⁽⁴⁾](https://codeberg.org/bbfh/vintage/wiki/Philosophy#4-user-control-and-predictability)

## Roadmap before v1

- [x] Inline templates
- [x] Generator templates
- [ ] Collector templates
- [ ] Better examples
- [ ] Patch should apply on the source code, so that function parsing and templates run on it.
- [ ] Overlays should be processed by the main pipeline, so that mcfunction and templates compile.
