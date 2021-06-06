println("Starting!")
do
    println("Printed at least once")
while False
a = 0
do
    b = 0
    do
        if (a + b) % 2 == 0
            println(a.ToString() + " + " + b.ToString() + ' = ' + (a+b).ToString())
            break
        end
        b += 1
    while b < 10
    a += 1
while a < 10
println("Done")