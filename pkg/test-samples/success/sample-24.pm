module Business
    class Person
        def __init__(name)
            self.name = name
        end

        def __string__()
            return self.name
        end
    end
end

class Citizen(Business.Person)
    id = 0
    country = 'US'
    def __string__()
        return self.name + " - FROM: " + self.country
    end
end

antonio = Citizen("Antonio")
println( antonio.__string__())
