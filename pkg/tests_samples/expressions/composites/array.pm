# Simple Array
array = [1, 2, 3, 4]
println(array.ToString() == "[1, 2, 3, 4]")
println(array[0] == 1)
println(array[-1] == 4)
# Nested Arrays
println([1, [[[[1, 2, 3]]]]].ToString() == "[1, [[[[1, 2, 3]]]]]")

a = [
1, 2,
3,      4,
                   5
                   ,
                   6
]

println(a.ToString() == "[1, 2, 3, 4, 5, 6]")