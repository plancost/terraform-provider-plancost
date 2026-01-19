package main
 
deny contains msg if {
	resource := input.resource_changes[_]
	resource.type == "plancost_estimate"
	resource.name == "this"
	cost := resource.change.after.monthly_cost
	to_number(cost) >= 2
	msg := sprintf("Monthly cost estimate is too high: $%v (limit: $2)", [cost])
}