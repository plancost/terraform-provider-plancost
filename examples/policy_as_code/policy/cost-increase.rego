package main
 
deny contains msg if {
    resource := input.resource_changes[_]
    resource.type == "plancost_estimate"
    resource.name == "this"
    
    after_cost := to_number(resource.change.after.monthly_cost)
    before_cost := get_before_cost(resource)
    diff := after_cost - before_cost
    
    diff >= 1
    msg := sprintf("Monthly cost increase is too high: $%v (limit: $1)", [diff])
}
 
get_before_cost(resource) := cost if {
    resource.change.before != null
    cost := to_number(resource.change.before.monthly_cost)
}
 
get_before_cost(resource) := 0 if {
    resource.change.before == null
}