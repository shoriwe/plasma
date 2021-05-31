a = 2
if a < 2
    println("Less")
elif a > 2
    println("Greater")
else
    println("Equal")
end

b = 2
if a == 2
    if a**b == 4
        println((2, 2))
    elif a**b == 8
        println((2, 3))
    end
elif a == 3
    if a ** b == 9
        println((3, 2))
    elif a ** b == 27
        println((3, 3))
    end
end