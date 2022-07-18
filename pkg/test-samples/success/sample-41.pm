def my_func()
    defer println("last")
    println("first")
    return 1
end

gen my_gen()
    defer println("once")
    for i in range(1, 3)
        yield i
    end
    return 3
end

println(my_func())

for number in my_gen()
    println(number)
end
