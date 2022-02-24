package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/prashantv/tf-test/internal/filestore"
)

func serviceResource() *schema.Resource {
	return &schema.Resource{
		Description: "Service resource",

		CreateContext: serviceCreate,
		ReadContext:   serviceRead,
		UpdateContext: serviceUpdate,
		DeleteContext: serviceDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "service name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"job": {
				Type:     schema.TypeSet,
				Elem:     newJobSchema(),
				MaxItems: 1,
				Required: true,
			},
		},
	}
}

func newJobSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"search_tags": {
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		},
	}
}

type service struct {
	Name string
	Job  job
}

type job struct {
	Namespace  string
	SearchTags map[string]string
}

func buildService(d *schema.ResourceData) *service {
	svc := &service{
		Name: d.Get("name").(string),
	}

	jobSet := d.Get("job").(*schema.Set)
	if jobSet.Len() == 0 {
		return svc
	}
	if jobSet.Len() > 1 {
		panic("jobSet has more elements than expected: " + fmt.Sprintf("%+v", jobSet.List()))
	}

	rawJob := jobSet.List()[0].(map[string]interface{})
	svc.Job.Namespace = rawJob["namespace"].(string)
	svc.Job.SearchTags = mapStrStr(rawJob["search_tags"].(map[string]interface{}))
	return svc
}

func serviceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id := fmt.Sprint(time.Now().UnixNano())
	svc := buildService(d)

	if err := filestore.Write(id, svc); err != nil {
		return diag.Errorf("failed to write: %v", err)
	}

	d.SetId(id)
	return nil
}

func serviceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var svc service
	err := filestore.Read(d.Id(), &svc)
	if err == filestore.ErrNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("read failed: %v", err)
	}

	if err := d.Set("name", svc.Name); err != nil {
		return diag.Errorf("failed to set name: %v", err)
	}

	if err := d.Set("job", []interface{}{map[string]interface{}{
		"namespace":   svc.Job.Namespace,
		"search_tags": svc.Job.SearchTags,
	}}); err != nil {
		return diag.Errorf("failed to set job: %v", err)
	}

	return nil
}

func serviceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	svc := buildService(d)

	if err := filestore.Write(d.Id(), svc); err != nil {
		return diag.Errorf("failed to write: %v", err)
	}

	return nil
}

func serviceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if err := filestore.Delete(d.Id()); err != nil {
		return diag.Errorf("failed to delete: %v", err)
	}

	return nil
}

func mapStrStr(m map[string]interface{}) map[string]string {
	out := make(map[string]string)
	for k, v := range m {
		out[k] = fmt.Sprint(v)
	}
	return out
}
