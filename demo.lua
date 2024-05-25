--- @diagnostic disable:unused-local,unused-function

local num = 3
local other_num = num
local other_other_num = other_num
local other_other_other_num = other_other_num

--- @type number
local uninitialized

local str, bool = "foo", true
local unknown = nil
local func = function(foo, bar) end

local type_error = num + str

--- Returns the sum of two numbers.
--- @param a number
--- @param b number
--- @return number
local function add(a, b)
  return a + b
end

local baz = add(true, 3)
baz = "lorem ipsum"
