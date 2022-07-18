number = 2

switch number
case 1
    println(false)
case 2
    println(true)
default
    println(false)
end

switch number
case 1
    println(false)
case 2
    switch number + 1
    case 3
        println(true)
    case 4
        println(false)
    end
default
    println(false)
end