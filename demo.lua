--- @type number
local uninitialized

local num, bool = 3, true
local str = "foo"
local unknown = nil
local func = function() end

local should_error = num + str -- Cannot add a 'number' with a 'string'.
