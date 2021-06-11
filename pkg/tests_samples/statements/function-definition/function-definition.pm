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
