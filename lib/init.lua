-- golang-resource-kafka-library main module.
-- Renders a Kafka sync producer setup into the service's messaging package.
--
-- The calling archetype is responsible for adding the corresponding
-- Go module dependency:
--   github.com/IBM/sarama
--
-- API (called from a parent archetype):
--   local kafka = require("golang-resource-kafka")
--   kafka.render(context, { destination = context:get("project-name") })
--
-- Context contract (no required keys beyond what the calling archetype provides).

local M = {}

function M.render(context, opts)
    opts = opts or {}
    local d = opts.destination
    if d and d ~= "" then
        directory.render("contents", context, { destination = d })
    else
        directory.render("contents", context)
    end
    return context
end

return M
