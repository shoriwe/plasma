println("" is String)
println("" implements String)

class Vehicle
    def __init__()
        self.model = "fastest"
    end
end

class Car(Vehicle)
    def __init__()
        self.type = "car"
    end
end


c = Car()
println(c is Car)
println(c implements Car)
println(c implements Vehicle)
