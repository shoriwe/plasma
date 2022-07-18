index = 0
reference = "Hello"
for a in (char for char in ((char2 * 10) for char2 in "Hello"))
    println(a)
    index += 1
end

println(index == 5)

reference = [(2, 1), (4, 3)]
index = 0
for a in ((y, x) for x, y in  [(1, 2), (3, 4)])
    println(a.__string__())
    index += 1
end

println(index == 2)