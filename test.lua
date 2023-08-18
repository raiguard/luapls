-- Numbers

local num_1 = 3
local num_2 = 3.0
local num_3 = 3.1416
local num_4 = 314.16e-2
local num_5 = 0.31416E1

local num_6 = 0xff
local num_7 = 0x0.1E
local num_8 = 0xA23p-4
local num_9 = 0X1.921FB54442D18P+1

-- Strings

local single_string = 'the quick fox'
local double_string = "the quick fox"

local bool_1 = false
local bool_2 = true
local nillit = nil

local exp = -1.53e-5-4
local exp_2 = 0xA32e.CDp-6432-5

-- Multiple assignment

local first, second = "second", "first"

-- Other

print(num_1, num_2)

if num_2 < 5 then
  print("not greater") print("lesser")
end

for i = 1, #double_string do
  print(double_string[i])
  goto continue
end

function foo(param)
  print(param, num_1)
end

local raw_string = [===[[[==[=]]==]====]
	[=[]=]===]

local other_raw = [[




food
]]

--[=[
--
--
-- foo bar [[ hello world ]]
]=

::continue::
