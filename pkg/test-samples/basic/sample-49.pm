switch Token.Kind
case Numeric, CommandOutput
	# break
case String
	print("I am a String")
default
	print("errors")
end