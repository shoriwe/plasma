a = 0
println((1 unless a == 2 else 0 unless a == 0 else false) == 1)
println(((1 if a + 1 * 2 == 2) unless 0) == 1)