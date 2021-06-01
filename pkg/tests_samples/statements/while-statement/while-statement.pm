while True
    println("1 Time")
    break
end
println("Done - 1")

a = 0
while a < 100
    if a == 10
        break
    elif a != 0 and a % 3 == 0
        println(a)
    else
        println("a not yet")
        a += 1
        continue
    end
    a += 1
end
println("Done - 2")