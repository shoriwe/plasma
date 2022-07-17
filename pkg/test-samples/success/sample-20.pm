# Simple Tuple
array = (1, 2, 3, 4)
println(array.__string__() == "(1, 2, 3, 4)")
println(array[0] == 1)
println(array[array.__len__()-1] == 4)
# Nested Tuples
println((1, ((((1, 2, 3),),),)).__string__())

a = (
1, 2,
3,      4,
                   5
                   ,
                   6
)

println(a.__string__() == "(1, 2, 3, 4, 5, 6)")