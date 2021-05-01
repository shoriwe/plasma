package vm


const  (
	AddOP uint16 = iota
	SubOP
	DivOP
	MulOP
	PowOP
	ModOP
	NegateBitsOP
	BitAndOP
	BitOrOP
	BitXorOP
	BitLeftOP
	BitRightOP
)