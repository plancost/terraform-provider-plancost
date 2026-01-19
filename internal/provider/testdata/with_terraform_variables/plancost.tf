# Auto-generated plancost estimates

resource "plancost_estimate" "this" {
  working_directory = abspath(path.module)
  var_file = abspath("${path.module}/terraform.tfvars")
}