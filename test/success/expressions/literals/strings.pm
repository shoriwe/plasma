# Single Quote Strings
println('Hello World' == "Hello World")
println('Hello
World' == "Hello\nWorld")
# Double Quote Strings
println("Hello World" == 'Hello World')
println("Hello
World" == 'Hello\nWorld')
# Byte Strings
println(b"Hello world"[2] == 108)
# Escaped chars
println("Hello\x41World" == "HelloAWorld")
println("\u0041ntonio" == "Antonio")
println("Hello\\x41World" == "Hello" + "\\" + "x41World")
println("500\u20ac" == "500€")
println("500\\u20ac" == "500" + "\\" + "u20ac")