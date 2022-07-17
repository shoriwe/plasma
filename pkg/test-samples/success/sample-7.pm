for number in (1, 2, 3)
    println(number)
end

for number in [1, 2, 3]
    println(number)
end

for char in "hello"
    println(char)
end

for a, b in ((1, 2), (2, 1))
    println(a, b)
end

n = none
for value in range(1, 100000, 1)
    n = value
end
println(n)

n = none
for value in range(1, 7.3, 1)
    n = value
end
println(n)