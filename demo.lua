-- --- @diagnostic disable:unused-local,unused-function

-- local num = 3
-- local other_num = num
-- local other_other_num = other_num
-- local other_other_other_num = other_other_num

-- --- @type number
-- local uninitialized

-- local str, bool = "foo", true
-- local unknown = nil
-- local func = function(foo, bar) end

-- local type_error = num + str

-- --- @class MyTable
-- --- @field foo string
-- --- @field bar string

local foo = "foo"
local bar = 123
foo = bar

do
  local hidden = true
end

hidden = false

local bar = "baz"

for i = 1, 10 do
  print(i)
end
