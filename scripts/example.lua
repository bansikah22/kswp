-- This is an example script for kswp.
-- It iterates over the list of unused resources and prints their names.

for i, resource in ipairs(resources) do
  print(resource.name)
end
