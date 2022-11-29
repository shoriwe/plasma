# Index ranges

Plasma supports a python like indexing for ranges by using **:** inside your index expressions

```ruby
a = [ 1, 2, 3, 4, 5 ]
println(a[3:]) # [ 4, 5 ]
println(a[:3]) # [ 1, 2, 3 ]
println(a[1:4]) # [ 2, 3, 4]
```

