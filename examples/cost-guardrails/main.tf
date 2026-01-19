resource "azurerm_resource_group" "example" {
  name     = "exampleRG1"
  location = "eastus"
}

resource "azurerm_public_ip" "example" {
  name                = "example-public-ip"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "plancost_estimate" "this" {
  working_directory = abspath(path.module)

  guardrail {
    # Block when estimated cost exceeds $3.
    condition = "monthly_cost_budget"
    threshold = 3
    action    = "block"
  }

  guardrail {
    # Warn when estimated cost increases by more than $200.
    condition = "monthly_cost_increase_amount"
    threshold = 200
    action    = "warning"
  }

  guardrail {
    # Warn when estimated cost increases by more than 10%.
    condition = "monthly_cost_increase_percentage"
    threshold = 10
    action    = "warning"
  }
}
