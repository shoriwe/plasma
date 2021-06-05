try
    Fail("Function not defined")
except ObjectWithNameNotFoundError
    println("Could Not find \"Fail\"")
end

try
    println("Function not defined", 2)
except ObjectWithNameNotFoundError as Error
    println("Not here")
except InvalidNumberOfArgumentsError as error
    println(error)
end

try
    println("Function not defined", 2)
except ObjectWithNameNotFoundError as Error
    println("Not here")
except InvalidNumberOfArgumentsError as error
    println(error)
finally
    println("Always the end!!!")
end

try
    println("Function not defined", 2)
except ObjectWithNameNotFoundError, RuntimeError, ObjectConstructionError as error
    println(error)
else
    print("No error matches the one received")
finally
    println("Always the end 2!!!")
end

try
    println("Function not defined", 2)
except ObjectWithNameNotFoundError, RuntimeError, InvalidNumberOfArgumentsError as error
    try
        error + 1
    except
        println("Nested try-except")
    finally
        println("Always the end 3!!!")
    end
else
    print("No error matches the one received")
finally
    println("Always the end 4!!!")
end
