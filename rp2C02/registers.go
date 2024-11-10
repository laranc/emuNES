package rp2C02

type StatusRegister struct {
	Unused         uint8
	SpriteOverflow uint8
	SpriteZeroHit  uint8
	VerticalBlank  uint8
	Reg            uint8 // Read only
}

type MaskRegister struct {
	Grayscale            uint8
	RenderBackgroundLeft uint8
	RenderSpritesLeft    uint8
	RenderBackground     uint8
	RenderSprite         uint8
	EnchanceRed          uint8
	EnchanceGreen        uint8
	EnchanceBlue         uint8
	Reg                  uint8 // Read only
}

type ControlRegister struct {
	NameTableX        uint8
	NameTableY        uint8
	IncrementMode     uint8
	PatternSprite     uint8
	PatternBackground uint8
	SpriteSize        uint8
	SlaveMode         uint8
	EnableNMI         uint8
	Reg               uint8 // Read only
}

func MakeStatusRegister() StatusRegister {
	return StatusRegister{
		Unused:         5,
		SpriteOverflow: 1,
		SpriteZeroHit:  1,
		VerticalBlank:  1,
		Reg:            0b00000111,
	}
}

func (r *StatusRegister) Update() {
	r.Reg = (r.Unused << 4) | (r.SpriteOverflow << 5) | (r.SpriteZeroHit << 6) | (r.VerticalBlank << 7)
}

func MakeMaskRegister() MaskRegister {
	return MaskRegister{
		Grayscale:            1,
		RenderBackgroundLeft: 1,
		RenderSpritesLeft:    1,
		RenderBackground:     1,
		RenderSprite:         1,
		EnchanceRed:          1,
		EnchanceGreen:        1,
		EnchanceBlue:         1,
		Reg:                  0b11111111,
	}
}

func (r *MaskRegister) Update() {
	r.Reg = (r.Grayscale << 0) | (r.RenderBackgroundLeft << 1) | (r.RenderSpritesLeft << 2) | (r.RenderBackground << 3) | (r.RenderSprite << 4) | (r.EnchanceRed << 5) | (r.EnchanceGreen << 6) | (r.EnchanceBlue << 7)
}

func MakeControlRegister() ControlRegister {
	return ControlRegister{
		NameTableX:        1,
		NameTableY:        1,
		IncrementMode:     1,
		PatternSprite:     1,
		PatternBackground: 1,
		SpriteSize:        1,
		SlaveMode:         1,
		EnableNMI:         1,
		Reg:               0b11111111,
	}
}

func (r *ControlRegister) Update() {
	r.Reg = (r.NameTableX << 0) | (r.NameTableY << 1) | (r.IncrementMode << 2) | (r.PatternSprite << 3) | (r.PatternBackground << 4) | (r.SpriteSize << 5) | (r.SlaveMode << 6) | (r.EnableNMI << 7)
}
