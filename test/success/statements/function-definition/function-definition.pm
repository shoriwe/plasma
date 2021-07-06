def pow(a, b)
    return a**b
end

println(pow(1, 2) == 1)

println(pow(2, 1) == 2)

println(pow(2, 2) == 4)

def special(a, b, c)
    return a ** 2 ** b - c
end

println(special(1, 2, 3) == -2)

println(special(2, 1, 3) == 1)

def help()
end

println(help() == None)

def return_for()
    for a in range(0, 10, 1)
        for b in range(100, 110, 1)
            return b
        end
        return a
    end
    return a
end

def return_while()
    a = 0
    while a < 10
        return a
        a += 1
    end
    return a
end

def return_do_while()
    a = 0
    do
        return a
        a += 1
    while a < 10
    return a
end

println(return_for() == 0)
println(return_while() == 0)
println(return_do_while() == 0)
