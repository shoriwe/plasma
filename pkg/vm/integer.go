package vm

/*
IntegerValue
	Class: Integer TODO
	Methods:
	- Positive TODO
	- Negative TODO
	- NegateBits TODO
	- Equals: Integer == Any TODO
	- NotEqual: Integer != Any TODO
	- GreaterThan: Integer > Integer, Integer > Float TODO
	- GreaterOrEqualThan: Integer >= Integer, Integer >= Float TODO
	- LessThan: Integer < Integer, Integer < Float TODO
	- LessOrEqualThan: Integer <= Integer, Integer <= Float TODO
	- BitwiseOr: Integer | Integer TODO
	- BitwiseXor: Integer ^ Integer TODO
	- BitwiseAnd: Integer & Integer TODO
	- BitwiseLeft: Integer << Integer TODO
	- BitwiseRight: Integer >> Integer TODO
	- Add: Integer + Integer, Integer + Float TODO
	- Sub: Integer - Integer, Integer - Float TODO
	- Mul: Integer * Integer, Integer * Float, Integer * String, Integer * Bytes, Integer * Array TODO
	- Div: Integer / Integer, Integer / Float TODO
	- FloorDiv: Integer // Integer, Integer // Float TODO
	- Modulus: Integer % Integer TODO
	- PowerOf: Integer ** Integer, Integer ** Float TODO
	- Bool TODO
	- Class TODO
	- Copy TODO
	- String TODO
*/
func (ctx *Context) IntegerValue(i int64) *Value {
	value := ctx.NewValue()
	value.Int = i
	// TODO: init symbols
	panic("implement me!")
	return value
}
