module Business
    class Person(Object)
        def Initialize(name)
            self.name = name
        end

        def ToString()
            return self.name
        end
    end
end

antonio = Business.Person("Antonio")
println(antonio)