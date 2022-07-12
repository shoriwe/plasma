interface Person
    def Initialize(name)
        self.name = name
    end

    def ToString()
        return self.name
    end
end

class Engineer(Person)
    University = "MIT"
    def Initialize(name)
        self.name = "Engineer " + name + " From: " + self.University
    end
end

found = false
antonio = Engineer("Antonio")

for subClass in antonio.SubClasses()
    if subClass == Person
        found = true
        break
    end
end

println(found)