package google

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

func resourceComputeHttpsHealthCheck() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeHttpsHealthCheckCreate,
		Read:   resourceComputeHttpsHealthCheckRead,
		Delete: resourceComputeHttpsHealthCheckDelete,
		Update: resourceComputeHttpsHealthCheckUpdate,

		Schema: map[string]*schema.Schema{
			"check_interval_sec": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"healthy_threshold": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2,
			},

			"host": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"port": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  80,
			},

			"request_path": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "/",
			},

			"self_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"timeout_sec": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},

			"unhealthy_threshold": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2,
			},
		},
	}
}

func resourceComputeHttpsHealthCheckCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	// Build the parameter
	hchk := &compute.HttpsHealthCheck{
		Name: d.Get("name").(string),
	}
	// Optional things
	if v, ok := d.GetOk("description"); ok {
		hchk.Description = v.(string)
	}
	if v, ok := d.GetOk("host"); ok {
		hchk.Host = v.(string)
	}
	if v, ok := d.GetOk("request_path"); ok {
		hchk.RequestPath = v.(string)
	}
	if v, ok := d.GetOk("check_interval_sec"); ok {
		hchk.CheckIntervalSec = int64(v.(int))
	}
	if v, ok := d.GetOk("healthy_threshold"); ok {
		hchk.HealthyThreshold = int64(v.(int))
	}
	if v, ok := d.GetOk("port"); ok {
		hchk.Port = int64(v.(int))
	}
	if v, ok := d.GetOk("timeout_sec"); ok {
		hchk.TimeoutSec = int64(v.(int))
	}
	if v, ok := d.GetOk("unhealthy_threshold"); ok {
		hchk.UnhealthyThreshold = int64(v.(int))
	}

	log.Printf("[DEBUG] HttpsHealthCheck insert request: %#v", hchk)
	op, err := config.clientCompute.HttpsHealthChecks.Insert(
		config.Project, hchk).Do()
	if err != nil {
		return fmt.Errorf("Error creating HttpsHealthCheck: %s", err)
	}

	// It probably maybe worked, so store the ID now
	d.SetId(hchk.Name)

	// Wait for the operation to complete
	w := &OperationWaiter{
		Service: config.clientCompute,
		Op:      op,
		Project: config.Project,
		Type:    OperationWaitGlobal,
	}
	state := w.Conf()
	state.Timeout = 2 * time.Minute
	state.MinTimeout = 1 * time.Second
	opRaw, err := state.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for HttpsHealthCheck to create: %s", err)
	}
	op = opRaw.(*compute.Operation)
	if op.Error != nil {
		// The resource didn't actually create
		d.SetId("")

		// Return the error
		return OperationError(*op.Error)
	}

	return resourceComputeHttpsHealthCheckRead(d, meta)
}

func resourceComputeHttpsHealthCheckUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	// Build the parameter
	hchk := &compute.HttpsHealthCheck{
		Name: d.Get("name").(string),
	}
	// Optional things
	if v, ok := d.GetOk("description"); ok {
		hchk.Description = v.(string)
	}
	if v, ok := d.GetOk("host"); ok {
		hchk.Host = v.(string)
	}
	if v, ok := d.GetOk("request_path"); ok {
		hchk.RequestPath = v.(string)
	}
	if v, ok := d.GetOk("check_interval_sec"); ok {
		hchk.CheckIntervalSec = int64(v.(int))
	}
	if v, ok := d.GetOk("healthy_threshold"); ok {
		hchk.HealthyThreshold = int64(v.(int))
	}
	if v, ok := d.GetOk("port"); ok {
		hchk.Port = int64(v.(int))
	}
	if v, ok := d.GetOk("timeout_sec"); ok {
		hchk.TimeoutSec = int64(v.(int))
	}
	if v, ok := d.GetOk("unhealthy_threshold"); ok {
		hchk.UnhealthyThreshold = int64(v.(int))
	}

	log.Printf("[DEBUG] HttpsHealthCheck patch request: %#v", hchk)
	op, err := config.clientCompute.HttpsHealthChecks.Patch(
		config.Project, hchk.Name, hchk).Do()
	if err != nil {
		return fmt.Errorf("Error patching HttpsHealthCheck: %s", err)
	}

	// It probably maybe worked, so store the ID now
	d.SetId(hchk.Name)

	// Wait for the operation to complete
	w := &OperationWaiter{
		Service: config.clientCompute,
		Op:      op,
		Project: config.Project,
		Type:    OperationWaitGlobal,
	}
	state := w.Conf()
	state.Timeout = 2 * time.Minute
	state.MinTimeout = 1 * time.Second
	opRaw, err := state.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for HttpsHealthCheck to patch: %s", err)
	}
	op = opRaw.(*compute.Operation)
	if op.Error != nil {
		// The resource didn't actually create
		d.SetId("")

		// Return the error
		return OperationError(*op.Error)
	}

	return resourceComputeHttpsHealthCheckRead(d, meta)
}

func resourceComputeHttpsHealthCheckRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	hchk, err := config.clientCompute.HttpsHealthChecks.Get(
		config.Project, d.Id()).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			// The resource doesn't exist anymore
			d.SetId("")

			return nil
		}

		return fmt.Errorf("Error reading HttpsHealthCheck: %s", err)
	}

	d.Set("host", hchk.Host)
	d.Set("request_path", hchk.RequestPath)
	d.Set("check_interval_sec", hchk.CheckIntervalSec)
	d.Set("health_threshold", hchk.HealthyThreshold)
	d.Set("port", hchk.Port)
	d.Set("timeout_sec", hchk.TimeoutSec)
	d.Set("unhealthy_threshold", hchk.UnhealthyThreshold)
	d.Set("self_link", hchk.SelfLink)

	return nil
}

func resourceComputeHttpsHealthCheckDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	// Delete the HttpsHealthCheck
	op, err := config.clientCompute.HttpsHealthChecks.Delete(
		config.Project, d.Id()).Do()
	if err != nil {
		return fmt.Errorf("Error deleting HttpsHealthCheck: %s", err)
	}

	// Wait for the operation to complete
	w := &OperationWaiter{
		Service: config.clientCompute,
		Op:      op,
		Project: config.Project,
		Type:    OperationWaitGlobal,
	}
	state := w.Conf()
	state.Timeout = 2 * time.Minute
	state.MinTimeout = 1 * time.Second
	opRaw, err := state.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for HttpsHealthCheck to delete: %s", err)
	}
	op = opRaw.(*compute.Operation)
	if op.Error != nil {
		// Return the error
		return OperationError(*op.Error)
	}

	d.SetId("")
	return nil
}
