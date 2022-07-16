module Business
    class Person( Value)
        def Initialize(name)
            self.name = name
        end

        def ToString()
            return self.name
        end
    end
end

class Person
end

antonio = Business.Person("Antonio")
println(antonio.ToString() == "Antonio")
