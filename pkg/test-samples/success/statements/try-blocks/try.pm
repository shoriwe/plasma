try
    Fail("Function not defined")
except ObjectWithNameNotFoundError
    println(True)
end

try
    println("Invalid Number of arguments", 2)
except ObjectWithNameNotFoundError as Error
    println(False)
except InvalidNumberOfArgumentsError as error
    println(True)
end

a = 0
try
    println("Invalid Number of arguments", 2)
except ObjectWithNameNotFoundError as Error
    println(False)
except InvalidNumberOfArgumentsError as error
    a = 2
    println(True)
finally
    println(a == 2)
end

a = 0
try
    println("Invalid Number of arguments", 2)
except ObjectWithNameNotFoundError, ObjectConstructionError as error
    println(False)
else
    a = 2
    println(True)
finally
    println(a == 2)
end

a = 0
try
    println("Invalid Number of arguments", 2)
except ObjectWithNameNotFoundError, RuntimeError, InvalidNumberOfArgumentsError as error
    a += 1
    try
        error + 1
    except
        a += 1
        println(True)
    finally
        println(a == 2)
        a += 1
    end
else
    println(False)
finally
    println(a == 3)
end

try
    raise InvalidNumberOfArgumentsError(0, 10)
except InvalidNumberOfArgumentsError as error
    println(True)
end

try
    raise InvalidNumberOfArgumentsError(0, 10)
except InvalidNumberOfArgumentsError as
        error
    println(True)
end