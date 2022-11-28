## Classes and interfaces

Interfaces work the same way as classes, but they only permit function and generator definitions

ALWAYS create `__init__` method inside your classes and interfaces:

```ruby
interface Vehicle
    def __init__()
        pass
    end
    
    def drive()
       println("driving")
    end
end

class Car(Vehicle)
    def __init__(name)
        self.name name
    end
end
```

# Modules

Modules are a special way to organize symbols

```ruby
module MyModule
    def calc(a)
        println(a**2)
    end
end

MyModule.calc(10)
```