package gameboy

type Buttons struct {
	butA  bool
	butB  bool
	sel   bool
	start bool
	right bool
	left  bool
	up    bool
	down  bool

	changed bool
}

func (b *Buttons) Set(name string, pressed bool) {
	switch name {
	case "A":
		b.butA = pressed
	case "B":
		b.butB = pressed
	case "Select":
		b.sel = pressed
	case "Start":
		b.start = pressed
	case "Right":
		b.right = pressed
	case "Left":
		b.left = pressed
	case "Up":
		b.up = pressed
	case "Down":
		b.down = pressed
	}

	b.changed = true
}

func (b *Buttons) Changed() bool {
	return b.changed
}

func (b *Buttons) ClearChanged() {
	b.changed = false
}

// NOTE! The buttons are active low!
// Note: 0 = pressed, 1 = not pressed 🙃
// https://gbdev.io/pandocs/Joypad_Input.html

func (b *Buttons) getPadState() byte {
	state := byte(0x0F)

	if b.down {
		state = bitReset(state, 3)
	}
	if b.up {
		state = bitReset(state, 2)
	}
	if b.left {
		state = bitReset(state, 1)
	}
	if b.right {
		state = bitReset(state, 0)
	}

	return state
}

func (b *Buttons) getButtonState() byte {
	state := byte(0x0F)

	if b.start {
		state = bitReset(state, 3)
	}
	if b.sel {
		state = bitReset(state, 2)
	}
	if b.butB {
		state = bitReset(state, 1)
	}
	if b.butA {
		state = bitReset(state, 0)
	}

	return state
}
