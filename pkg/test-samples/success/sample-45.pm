a = 1
println(a.big_endian().__array__().__string__())
b = a.big_endian()
c = 0.from_big(b)
println(c)
a = 1
println(a.little_endian().__array__().__string__())
b = a.little_endian()
c = 0.from_little(b)
println(c)