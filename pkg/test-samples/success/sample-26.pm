a = 0
map = {}
while a < 100
    map[a] = a
    map[a.__string__()] = a
    a += 1
end
println(map["11"])
