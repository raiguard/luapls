-- NUMBERS

local num_1 = 3
local num_2 = 3.0
local num_3 = 3.1416
local num_4 = 314.16e-2
local num_5 = 0.31416E1

local num_6 = 0xff
local num_7 = 0x0.1E
local num_8 = 0xA23p-4
local num_9 = 0X1.921FB54442D18P+1

local foo = 2 + -3
local bar = 'the quick fox'

local exp = -1.53e-5-4
local exp_2 = 0xA32e.CDp6432

print(foo, bar)

if foo < 5 then
  local baz = foo local foo = bar
  print("greater") print("lesser")
end

for i = 1.532, #bar do
  print(bar[i])
  goto continue
end

::continue::
