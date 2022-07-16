reference = ("Antonio", "Juan")
names = {"Antonio": "Developer", "Juan": "Analyst"}
index = 0
for name in names
    println(reference[index] == name)
    index += 1
end

reference = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9]
for number in range(0, 10, 1)
    println(reference[number] == number)
end