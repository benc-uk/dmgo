# DMGO - A Gameboy emulator written in Go

Another Gameboy emulator written in Go, using [Ebitengine](https://ebitengine.org/) for rendering and display

![screen](./etc/screens/tetris-0.png)

## Status

- Boots some ROMs, and runs the Gameboy boot ROM if present
- Tetris is playable!
- 100% of the CPU opcodes working and passing [Blargg's tests](https://github.com/retrio/gb-test-roms)
- PPU & LCD: Functional rendering but needs major work
- Nearly all interrupts
- Timing & HALT: Passes Blargg's interrupt test ROM
- No sound

## Todo Next

- Other interrupts: LCD STAT & serial
- Correct & update STAT register
- Render correctly per scanline

Longer term

- Support MBC1
- Remove dependency on ebiten from Gameboy package, in PPU

## Reference Collection

Docs, so many docs

- https://gbdev.io/
- https://gbdev.io/pandocs/
- https://gbdev.io/gb-opcodes/optables/dark
- https://github.com/Gekkio/gb-ctr
- https://gbdev.io/pandocs/CPU_Instruction_Set.html

Inspiration projects

- https://github.com/jacoblister/emuboy
- https://github.com/Humpheh/goboy
- https://github.com/ArcticXWolf/AXWGameboy

Cool videos

- https://www.youtube.com/watch?v=HyzD8pNlpwI
