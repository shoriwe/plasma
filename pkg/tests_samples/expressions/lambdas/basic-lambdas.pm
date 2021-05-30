a = (lambda x, y: (x**2+y**2)**(1/2))
a(2, 1)
c = [1, 2, lambda x, y: x**y]
c[2](2, 2).Value = 2