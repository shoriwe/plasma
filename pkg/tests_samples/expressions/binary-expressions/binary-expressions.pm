println(((1 + 2 / 3) == (1 + 2 / 3)) and ((1 + 2 / 3).ToString() == "1.6666666666666666667"))
println(25**(1/2) == 5)
println("Hello " * 5 == "Hello Hello Hello Hello Hello ")
println((1, 2, 4 + 5 / 6 ** 2, 10, "hello * 5 " * 0) == (1, 2, 4 + 5 / 6 ** 2, 10, ""))
println(1 and (1, 2, 3, 4))
println(1 or (1, 2, 3, 4))
println(1 xor (1, 2, 3, 4) == False)
println(1 in (1, 2, 3, 4))
println(1 // 2 == 0)
println((1, 2, 3, "Hello") * 2 == (1, 2, 3, "Hello", 1, 2, 3, "Hello"))

class A
    def Equals(other)
        return self.Class() == other.Class()
    end
end

println((A(), A(), A(), A()) * 2 == (A(), A(), A(), A(), A(), A(), A(), A()))
println(1 / 2 == 0.5)
println(1 / 2 ** 2 + 5 * 1 - 3 == 2.25)
println(1 / 2 == 0.5)