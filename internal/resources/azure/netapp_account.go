package azure

import (
	"github.com/plancost/terraform-provider-plancost/internal/schema"
)

type NetAppAccount struct {
	Address string
	Region  string
}

func (r *NetAppAccount) CoreType() string {
	return "NetAppAccount"
}

func (r *NetAppAccount) UsageSchema() []*schema.UsageItem {
	return []*schema.UsageItem{}
}

func (r *NetAppAccount) PopulateUsage(u *schema.UsageData) {
}

func (r *NetAppAccount) BuildResource() *schema.Resource {
	return &schema.Resource{
		Name:      r.Address,
		IsSkipped: true,
		NoPrice:   true,
	}
}
