gen my_gen()
    for a in range(1, 20)
        yield a
    end
    if ref == 100
        return true
    end
    return false
end

ref = 100
for a in my_gen()
    println(a)
end
ref = 0
for a in my_gen()
    println(a)
end