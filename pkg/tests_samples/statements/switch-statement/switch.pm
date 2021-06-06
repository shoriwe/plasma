number = 2

switch number
case 1
    println("Number 1")
case 2
    println("Number 2")
default
    println("Invalid number")
end

switch number
case 1
    println("Number 1")
case 2
    switch number + 1
    case 3
        println("Number 2 + 1 = 3")
    case 4
        println("Number 2 + 1 = 4")
    end
default
    println("Invalid number --")
end