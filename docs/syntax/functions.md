# Anonymous functions

Define anonymous functions, meaning you don't need to assign them to a symbol:

```ruby
result = (lambda x: x ** 2)(5)
println(result)
```

Expected output:

```
25
```

# Functions

To define functions you will make use of the keyword `def`:

```ruby
def my_func(a, b)
    return a ** b
end

println(my_func(5, 2))
```

Expected output:

```
25
```

# Generators

Generators are special functions that can be used to simplify the implementation iterator objects:

To define generators you will make use of the keyword `genb`, When `yield` is used the function `__next__()` function of
the object returns the iterated value, on `return` the iterator ends:

```ruby
gen my_gen(pow)
    for number in range(1, 5)
        yield number ** pow
    end
    return pow ** (1/2)
end

for result in my_gen(25)
    println(result)
end
```

The expected output will be:

```
1
33554432
847288609443
1125899906842624
5.000000
```

# Defer statement

The `defer` statement is used to call a function just after the `return` statement is evaluated.

```ruby
def my_func()
    defer println("third")      # This will print third
    defer println("last")       # This will print last
    println("first")            # This will print first
    return println("second")    # This till print second
end

my_func()
```

Expected output:

```
first
second
third
last
```
