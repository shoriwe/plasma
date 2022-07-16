module Caller
    def hello(message)
        return message
    end
end

println(Caller.hello("Hello John") == "Hello John")