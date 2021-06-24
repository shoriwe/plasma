number = 2

switch number
case 1
    println(False)
case 2
    println(True)
default
    println(False)
end

switch number
case 1
    println(False)
case 2
    switch number + 1
    case 3
        println(True)
    case 4
        println(False)
    end
default
    println(False)
end