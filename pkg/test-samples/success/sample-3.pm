println((1 + 2 / 3))
println(25**(1/2))
println("Hello " * 5)
println((1, 2, 4 + 5 / 6 ** 2, 10, "hello * 5 " * 0).__string__() == "(1, 2, 4.138889, 10, )")
println(1 and (1, 2, 3, 4))
println(1 or (1, 2, 3, 4))
println(1 xor (1, 2, 3, 4))
println(1 in (1, 2, 3, 4))
println(1 // 2)
println((1, 2, 3, "Hello"))
println(1 in [1, 2, 3, 4, 5])
println(1 in (1, 2, 3, 4, 5))
println(1 in {1: 2, 2: 3, 3: 4, 4: 5, 5: 6})
println(((1 + 2) / 3))
println(1 / 2)
println(1 / 2 ** 2 + 5 * 1 - 3)
println(1 / 2)

class A
    def __init__()
        pass
    end
    def __equals__(other)
        return self.__class__() == other.__class__()
    end
end

println(A() == A())