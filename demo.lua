-- Why does this function exist?
local function add(a, b)
  return a + b
end

print(add(3, 4))

local items = {5, 10, 15}

for i = 1, #items do
  print(add(i, items[i]))
end
