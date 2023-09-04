::beginning::

local i = 0

while i < 100 do
  i = i + 1
  if i % 2 == 0 then
    i = i + 1
  end
end

repeat
  i = i - 2
until i == 50

if i == 50 then
  goto beginning
end

-- Not valid Lua. Oh well!
for i = 1, 100, -1 do
  i = (i - 1 / 4) ^ 16 % 3
end

for key, value in table do
  return false
end
