a = 0

until a > 100
    a += 1
end

println(a)

a = 0
until a > 100
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

println(a)