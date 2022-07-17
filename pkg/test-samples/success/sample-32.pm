interface Person
    def __init__(name)
        self.name = name
    end

    def __string__()
        return self.name
    end
end

class Engineer(Person)
    University = "MIT"
    def __init__(name)
        self.name = "Engineer " + name + " From: " + self.University
    end
end

found = false
antonio = Engineer("Antonio")

for subClass in antonio.__sub_classes__()
    if subClass == Person
        found = true
        break
    end
end

println(found)