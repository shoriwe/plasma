def fib(n)
    if n <= 1
        return 1
    end
    return (fib(n - 1) + fib(n - 2))
end
a = 20
println(fib(a))