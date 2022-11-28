# Control flow

- Nop expression: `pass` used to do nothing
- `if` and `unless`:

If and unless blocks are the basic condition control flow of the language.

Unless works the same way as `if` but it previously negates the expression.

```ruby
if my_condition
    pass
elif other_condition
    pass
else
    pass
end

unless my_condition
    pass
elif other_condition
    pass
else
    pass
end
```

- `switch` statements

```ruby
a = 1
switch a
case 1
    pass
case 2
    pass
case 3
    pass
default
    pass
end
```

- `BEGIN` and `END` blocks

This two blocks are executed before the string `BEGIN` and at the end of the script `END`:

```ruby
BEGIN
    println("first")
end

println("middle")

END
    println("end")
end
```