redo # 1

def my_function()
    redo # 2
end

gen my_generator()
    redo # 3
end

for a in range(100)
    redo
end

for value in range(2000)
    gen __anonymous()
        redo # 4
    end
end

if False
    redo # 5
end