continue # 1

def my_function()
    continue # 2
end

gen my_generator()
    continue # 3
end

for a in range(100)
    continue
end

for value in range(2000)
    gen __anonymous()
        continue # 4
    end
end

if False
    continue # 5
end