# Basics

## Assignments

You can assign values to:

- Identifiers:

```ruby
my_variable = "Hello"
```

- Selectors:

```ruby
my_variable.my_property = "Hello"
```

- Indexes:

```ruby
my_array[2] = "Hello"

my_hash["Antonio"] = "Hello"
```

## Deleting symbols

Use the `delete` statement to delete symbols, selectors and indexes:

```ruby
delete my_variable              # Remove the symbol `my_variable`
delete my_variable.my_property  # Remove `my_property` from `my_variable` object
delete my_hash['Work']          # Remove the Key `Work` from `my_hash`
```

## Unary operators

Not: `not`, `!`
Negate bits: `~`

## Binary operators

- Boolean And: `and`, `&&`
- Boolean Or: `or`, `||`
- Boolean Xor: `xor`
- In: `in`
- Is: `is`
- Implements: `implements`
- Equals: `==`
- Not Equals: `!=`
- Greater than: `>`
- Greater or equal than: `>=`
- Less than: `<`
- Less or equal than: `<=`
- Bitwise or: `|`
- Bitwise and: `&`
- Bitwise xor: `^`
- Bitwise left: `<<`
- Bitwise right: `>>`
- Add: `+`
- Sub: `-`
- Div: `/`
- Floor division: `//`
- Mul: `*`
- Mod: `%`
- Pow: `**`

## Array expressions

Arrays can be defined using `[` and `]`, separating its internal elements with commas:

```ruby
my_array = [0, 1, 2, 3, "A string", [0, 1, 2]]
```

Arrays have specific behavior methods:

- `append(value)`: appends a value to the array
- `clear()`: clears the contents of the array
- `index(value)`: returns the index of a value in the array, `-1` if not found
- `pop()`: remove and returns the last element of the array
- `insert(index, value)`: inserts at index a new value
- `remove(index)`: remove element at index

## Tuple expressions

Tuples can be defined using `(` and `)`, separating its internal elements with commas:

```ruby
my_tuple = (0, 1, 2, 3, "A string", [0, 1, 2])
```

Notice that tuples are immutable, meaning you can not modify them but the elements inside them can.

## Hash expressions

Hash expressions, also known as map can be defined with `{` and `}`:

```ruby
my_hash = {
    "Antonio":  "Developer",
    "Victor":   "Administrator"
}
```

## String and bytes expressions

Strings can be defined of 3 ways:

```ruby
single_quote = 'My string'
double_quote = "My string"
back_quote   = `My string`
```

Bytes work the same as String but prepending a letter `b` before the first quote:

```ruby
single_quote = b'My bytes'
double_quote = b"My bytes"
back_quote   = b`My bytes`
```

Special methods for both:

- `join(tuple|array)`: returns a string with the content of the container but separating them with the contents of the
  original string
- `split(sep)`: returns a tuple with the string separated using the pattern of `sep`
- `upper()`: returns a new string but uppercase
- `lower()`: returns a new string but lowercase
- `count(pattern)`: counts how many times a pattern is inside the string
- `index(pattern)`: returns the index of the first pattern in the string, `-1` if not found

## Numbers

Plasma has integers and float that can be operated between both:

```ruby
my_int = 10
my_float = 2.0
my_result = my_int ** my_float
```

Special methods of both types are:

- `to_big()`: returns a bytes string with the 64 bit big endian contents of the number
- `from_big(bytes)`: reconstruct the number from the big endian bytes of the string
- `to_little()`: returns a bytes string with the 64 bit little endian contents of the number
- `from_little(bytes)`: reconstruct the number from the little endian bytes of the string

## Booleans

There are to booleans `true` and `false`

```ruby
# Alert of blocking code
while true
    pass
end

until false
    pass
end
```

## None

```ruby
a = none
```

## Functions calls

```ruby
my_func()
(lambda x, y: x + y)(1, 2)
```

## One line expressions

- Conditions `if` and `unless`:

```ruby
a = 2
b = 5 if a == 2 else 10
b = 5 unless a == 2 else 10
```

- Generators:

```ruby
for pow in (number ** 2 for number in range(1, 10))
    println(pow)
end
```