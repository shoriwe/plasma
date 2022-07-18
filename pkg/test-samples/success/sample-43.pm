a = b", ".join((b"Antonio", b"Juan"))
println(a)

b = a.split(b", ")

println(b.__string__())

println(b"welcome".upper())
println(b"WELCOME".lower())

println(b"AAABBBCDDDDEEEEEEEEFF".count(b"C"))
println(b"AAABBBCDDDDEEEEEEEEFF".index(b"C"))
println(b"AAABBBCDDDDEEEEEEEEFF".index(b"Z"))