a = 2
if a < 2
    println(false)
elif a > 2
    println(false)
else
    println(true)
end

b = 2
if a == 2
    if a**b == 4
        println(true)
    elif a**b == 8
        println(false)
    end
elif a == 3
    if a ** b == 9
        println(false)
    elif a ** b == 27
        println(false)
    end
end