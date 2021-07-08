unless False
    println(True)
else
    println(False)
end

a = 2
unless a < 2
    println(True)
elif a > 2
    println(False)
else
    println(False)
end

b = 2
unless a == 2
    unless a ** b == 4
        println(False)
    elif a ** b == 8
        println(False)
    end
elif a == 3
    unless a ** b == 9
        println(True)
    elif a ** b == 27
        println(False)
    end
end
