package datashare

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datashare/mgmt/2019-11-01/datashare"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datashare/helper"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datashare/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datashare/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func dataSourceDataShareDatasetBlobStorage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmDataShareDatasetBlobStorageRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.DatashareDataSetName(),
			},

			"share_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.DataShareID,
			},

			"container_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"storage_account_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"storage_account_resource_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"storage_account_subscription_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"file_path": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"folder_path": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArmDataShareDatasetBlobStorageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.DataSetClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	shareID := d.Get("share_id").(string)
	shareId, err := parse.DataShareID(shareID)
	if err != nil {
		return err
	}

	respModel, err := client.Get(ctx, shareId.ResourceGroup, shareId.AccountName, shareId.Name, name)
	if err != nil {
		return fmt.Errorf("retrieving DataShare Blob Storage DataSet %q (Resource Group %q / accountName %q / shareName %q): %+v", name, shareId.ResourceGroup, shareId.AccountName, shareId.Name, err)
	}

	respId := helper.GetAzurermDataShareDataSetId(respModel.Value)
	if respId == nil || *respId == "" {
		return fmt.Errorf("empty or nil ID returned for reading DataShare Blob Storage DataSet %q (Resource Group %q / accountName %q / shareName %q)", name, shareId.ResourceGroup, shareId.AccountName, shareId.Name)
	}

	d.SetId(*respId)
	d.Set("name", name)
	d.Set("share_id", shareID)

	switch resp := respModel.Value.(type) {
	case datashare.BlobDataSet:
		if props := resp.BlobProperties; props != nil {
			d.Set("container_name", props.ContainerName)
			d.Set("storage_account_name", props.StorageAccountName)
			d.Set("storage_account_resource_group_name", props.ResourceGroup)
			d.Set("storage_account_subscription_id", props.SubscriptionID)
			d.Set("file_path", props.FilePath)
			d.Set("display_name", props.DataSetID)
		}

	case datashare.BlobFolderDataSet:
		if props := resp.BlobFolderProperties; props != nil {
			d.Set("container_name", props.ContainerName)
			d.Set("storage_account_name", props.StorageAccountName)
			d.Set("storage_account_resource_group_name", props.ResourceGroup)
			d.Set("storage_account_subscription_id", props.SubscriptionID)
			d.Set("folder_path", props.Prefix)
			d.Set("display_name", props.DataSetID)
		}

	case datashare.BlobContainerDataSet:
		if props := resp.BlobContainerProperties; props != nil {
			d.Set("container_name", props.ContainerName)
			d.Set("storage_account_name", props.StorageAccountName)
			d.Set("storage_account_resource_group_name", props.ResourceGroup)
			d.Set("storage_account_subscription_id", props.SubscriptionID)
			d.Set("display_name", props.DataSetID)
		}

	default:
		return fmt.Errorf("data share dataset %q (Resource Group %q / accountName %q / shareName %q) is not a blob storage dataset", name, shareId.ResourceGroup, shareId.AccountName, shareId.Name)
	}

	return nil
}
