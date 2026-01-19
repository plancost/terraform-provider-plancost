# Tagging Policy Example

This example demonstrates how to enforce tagging standards using the `tagging_policy` block in `plancost_estimate`.

It covers:
- Warning if a required tag (e.g., `Owner`) is missing.
- Blocking if a tag value is not in the allowed list (e.g., `Environment`).
- Blocking if a tag value does not match a regex pattern (e.g., email format for `Owner`).
