a = 0
map = {}
while a < 100
    map[a] = a
    map[a.ToString()] = a
    a += 1
end
println(map.ToTuple().GetLength())
