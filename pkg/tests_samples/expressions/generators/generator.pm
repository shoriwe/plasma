for a in (char for char in ((char2 * 10) for char2 in "Hello"))
    println(a)
end

println("Then...")

for a in ((y, x) for x, y in  [(1, 2), (3, 4)])
    println(a)
end