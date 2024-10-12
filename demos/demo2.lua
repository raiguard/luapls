local foo, bar

local i, j = 0;;

local l, m

local helpers = {}

--- Returns the sum of two numbers.
--- @param a number
--- @param b number
--- @return number
function helpers.add(a, b)
  return a + b;
end

--- This is a comment
  print("foo bar") -- Comment after

  do
for i = 1, 10 do
  if i > 5 then
    break
  end
end
  end

repeat
  print("hallo " .. i)
  i = i + 1
  if i > 20 then
    return
  end
until i == 10

i = 0

while i < 5 do
  i = i + 0.1
end
