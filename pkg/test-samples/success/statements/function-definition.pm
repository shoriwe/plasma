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

println(help() == none)

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
        b = 100
        while b < 110
            return b
        end
        return a
        a += 1
    end
    return a
end

def return_do_while()
    a = 0
    do
        b = 100
        do
            return b
            b += 1
        while a < 110
        return a
        a += 1
    while a < 10
    return a
end

def return_if()
    if true
        return true
    end
    return false
end

def return_unless()
    unless false
        return true
    end
    return false
end

def fib(n)
    if n == 0
        return 0
    end
    if n == 1
        return 1
    end
    return fib(n-1) + fib(n-2)
end

println(return_for() == 100)
println(return_while() == 100)
println(return_do_while() == 100)
println(return_if())
println(return_unless())
println(fib(10) == 55)