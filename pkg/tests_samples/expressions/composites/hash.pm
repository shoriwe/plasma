# Simple Hash
println({"Hello": "World", 1: 2}["Hello"] == "World")

# Tuple hashing
my_hash = {(1, 2): 1, (2, 1): 100, ("Hello", "World"): 0}
println(my_hash == {(1, 2): 1, (2, 1): 100, ("Hello", "World"): 0})

println({(1, 2, 3): 1, (3, 2, 34, 1): 100, ("Hello", "a", 0): 0}[(3, 2, 34, 1)] == 100)

# Nested Hash
my_nested_hash = {0: 0, 1: [1, 2, 3, {0: 1232456, "a word": 2, 3: [3, 2]}]}

println(my_nested_hash[1][3][3][0] == 3)

a = {
    1: 2,
    2: 3
}

println(a[1] == 2 and a[2] == 3)