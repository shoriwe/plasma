# Simple Array
array = [1, 2, 3, 4]
println(array.__string__())
println(array[0])
println(array[array.__len__()-1])
# Nested Arrays
println([1, [[[[1, 2, 3]]]]].__string__())

a = [
1, 2,
3,      4,
                   5
                   ,
                   6
]

println(a.__string__())