a = 0
while True
    a += 1
    break
end
println(a == 1)

a = 0
while a < 100
    if a == 10
        break
    elif a != 0 and a % 3 == 0
        a += 2
    else
        a += 1
        continue
    end
    a += 1
end
println(a == 102)