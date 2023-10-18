local function distance_to(pos)
  return math.sqrt(pos.x^2 + pos.y^2)
end

print(distance_to({x = 3, y = 4}))
