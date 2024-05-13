# DMGO - A Gameboy emulator written in Go

Another Gameboy emulator written in Go, using [Ebitengine](https://ebitengine.org/) for rendering and display

![screen](./etc/screens/tetris-0.png)

## Status

- Boots some ROMs, and runs the Gameboy boot ROM if present
- About 60% of the CPU opcodes
- 100% of the CB prefix opcodes
- PPU: Functional rendering but incomplete (e.g. scanlines)
- No interrupts other than vblank
- Timing: CPU and PPU mostly? clock correct
- No input
- No sound

## Todo

- Finish instructions
- Other interrupts, timer etc
- Input
- Render correctly
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
