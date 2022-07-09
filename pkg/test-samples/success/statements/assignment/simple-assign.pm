a = "Hello "
a += "Rosmary" # a = a + "Rosmary"
println(a == "Hello Rosmary")
a.Age = 48
a.Age.BornDate = 1972

println(a.Age == 48 and a.Age.BornDate == 1972)
# Overwrite the variable
a = 1
a *= 10 / 83475987
println(a == 1.1979492976824580702e-07)

array = [1, 2, 3, 4]
array[2] = 10000
println(array == [1, 2, 10000, 4])