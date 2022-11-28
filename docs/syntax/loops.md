# Loops

- `do-while`: execute the code before evaluating the exit condition:

```ruby
do
    println("at least once")
while my_condition()
```

- `while` and `until`:

```ruby
while my_condition()
    # do something
    pass
end
```

Until works similar to `while`, the key difference is that it negates the condition prior evaluation:

```ruby
until false
    # do something
    pass
end
```

- `for`: iterates using `__next__` and `__has_next__` results:

```ruby
for number in [1, 2, 3]
    println(a)
end
```

## Loop controls

- `continue`
- `break`
