println((lambda x, y: (x**2+y**2)**(2))(1, 1) == 4)
a = (lambda x, y: (x**2+y**2)**(2))
println(a(1, 1) == 4)
c = [1, 2, 3, (lambda x, y: (x**2+y**2)**(2))]
println(c[3](1, 1) == 4)
a = (lambda x, y: (x**2+y**2)**(2))
a.b = 1
a.c = 1
println(a(a.b, a.c) == 4)
println((lambda: 1)() == 1)