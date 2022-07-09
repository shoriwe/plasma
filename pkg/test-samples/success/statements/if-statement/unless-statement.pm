unless false
    println(true)
else
    println(false)
end

a = 2
unless a < 2
    println(true)
elif a > 2
    println(false)
else
    println(false)
end

b = 2
unless a == 2
    unless a ** b == 4
        println(false)
    elif a ** b == 8
        println(false)
    end
elif a == 3
    unless a ** b == 9
        println(true)
    elif a ** b == 27
        println(false)
    end
end
