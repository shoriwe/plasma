a = 1.0
println(a.big_endian().__array__().__string__())
b = a.big_endian()
c = 0.0.from_big(b)
println(c)
a = 1.0
println(a.little_endian().__array__().__string__())
b = a.little_endian()
c = 0.0.from_little(b)
println(c)

m = [0, 0, 0, 0, 0, 0, 0, 0]
m[0] = 63
m[1] = 240

b = 0.0.from_big(m.__bytes__())
println(b)