--- @type number
local uninitialized

local num = 3
local str = "foo"

local should_error = num + str -- Cannot add a 'number' with a 'string'.
