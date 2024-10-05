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

--- @param input string
local function print(input) end

local tbl = {
  first = 1,
  second = "two",
  third = true,
}

tbl.fourth = function() end

print(tbl.fourth)

--- @param first string
function tbl.fifth(first) end

tbl.sixth = function(first, ...)
  print"bar"
end

script.on_event(defines.events.on_player_created, function(e)
  local player = game.get_player(e.player_index)
  if not player then
    return
  end
  player.print("Hello, world!")
end)

for i = 1, 10 do
  print(i)
end
