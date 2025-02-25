// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/scc-go-sdk/v4/posturemanagementv2"
)

func ResourceIBMSccPostureCredentials() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMSccPostureCredentialsCreate,
		ReadContext:   resourceIBMSccPostureCredentialsRead,
		UpdateContext: resourceIBMSccPostureCredentialsUpdate,
		DeleteContext: resourceIBMSccPostureCredentialsDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Credentials status enabled/disbaled.",
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_scc_posture_credential", "type"),
				Description:  "Credentials type.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_scc_posture_credential", "name"),
				Description:  "Credentials name.",
			},
			"description": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_scc_posture_credential", "description"),
				Description:  "Credentials description.",
			},
			"display_fields": {
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Description: "Details the fields on the credential. This will change as per credential type selected.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ibm_api_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The IBM Cloud API Key. This is mandatory for IBM Credential Type ie when type=ibm_cloud.",
						},
					},
				},
			},
			"purpose": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_scc_posture_credential", "purpose"),
				Description:  "Purpose for which the credential is created.",
			},
		},
	}
}

func ResourceIBMSccPostureCredentialsValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "aws_cloud, azure_cloud, database, ibm_cloud, kerberos_windows, ms_365, openstack_cloud, username_password",
		},
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-zA-Z0-9-._,\s]*$`,
			MinValueLength:             3,
			MaxValueLength:             30,
		},
		validate.ValidateSchema{
			Identifier:                 "description",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-zA-Z0-9-._,\s]*$`,
			MinValueLength:             1,
			MaxValueLength:             255,
		},
		validate.ValidateSchema{
			Identifier:                 "purpose",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "discovery_collection, discovery_collection_remediation, discovery_fact_collection, discovery_fact_collection_remediation, remediation",
			Regexp:                     `^[a-zA-Z0-9-\\.,_\\s]*$`,
			MinValueLength:             1,
			MaxValueLength:             100,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_scc_posture_credentials", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMSccPostureCredentialsCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	postureManagementClient, err := meta.(conns.ClientSession).PostureManagementV2()
	if err != nil {
		return diag.FromErr(err)
	}

	createCredentialOptions := &posturemanagementv2.CreateCredentialOptions{}

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error getting userDetails %s", err))
	}
	createCredentialOptions.SetAccountID(userDetails.UserAccount)

	createCredentialOptions.SetEnabled(d.Get("enabled").(bool))
	createCredentialOptions.SetType(d.Get("type").(string))
	createCredentialOptions.SetName(d.Get("name").(string))
	createCredentialOptions.SetDescription(d.Get("description").(string))
	displayFields := resourceIBMSccPostureCredentialsMapToNewCredentialDisplayFields(d.Get("display_fields.0").(map[string]interface{}))
	createCredentialOptions.SetDisplayFields(&displayFields)
	createCredentialOptions.SetPurpose(d.Get("purpose").(string))

	credential, response, err := postureManagementClient.CreateCredentialWithContext(context, createCredentialOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateCredentialWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateCredentialWithContext failed %s\n%s", err, response))
	}

	d.SetId(*credential.ID)

	return resourceIBMSccPostureCredentialsRead(context, d, meta)
}

func resourceIBMSccPostureCredentialsMapToNewCredentialDisplayFields(newCredentialDisplayFieldsMap map[string]interface{}) posturemanagementv2.NewCredentialDisplayFields {
	newCredentialDisplayFields := posturemanagementv2.NewCredentialDisplayFields{}

	if newCredentialDisplayFieldsMap["ibm_api_key"] != nil {
		newCredentialDisplayFields.IBMAPIKey = core.StringPtr(newCredentialDisplayFieldsMap["ibm_api_key"].(string))
	}

	return newCredentialDisplayFields
}

func resourceIBMSccPostureCredentialsMapToUpdateCredentialDisplayFields(updateCredentialDisplayFieldsMap map[string]interface{}) posturemanagementv2.UpdateCredentialDisplayFields {
	updateCredentialDisplayFields := posturemanagementv2.UpdateCredentialDisplayFields{}

	if updateCredentialDisplayFieldsMap["ibm_api_key"] != nil {
		updateCredentialDisplayFields.IBMAPIKey = core.StringPtr(updateCredentialDisplayFieldsMap["ibm_api_key"].(string))
	}

	return updateCredentialDisplayFields
}

func resourceIBMSccPostureCredentialsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	postureManagementClient, err := meta.(conns.ClientSession).PostureManagementV2()
	if err != nil {
		return diag.FromErr(err)
	}

	getCredentialsOptions := &posturemanagementv2.GetCredentialOptions{}
	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error getting userDetails %s", err))
	}

	accountID := userDetails.UserAccount
	getCredentialsOptions.SetAccountID(accountID)
	getCredentialsOptions.SetID(d.Id())

	credential, response, err := postureManagementClient.GetCredentialWithContext(context, getCredentialsOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetCredentialWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetCredentialWithContext failed %s\n%s", err, response))
	}
	d.SetId(*(credential.ID))
	return nil
}

func resourceIBMSccPostureCredentialsNewCredentialDisplayFieldsToMap(newCredentialDisplayFields posturemanagementv2.NewCredentialDisplayFields) map[string]interface{} {
	newCredentialDisplayFieldsMap := map[string]interface{}{}

	if newCredentialDisplayFields.IBMAPIKey != nil {
		newCredentialDisplayFieldsMap["ibm_api_key"] = newCredentialDisplayFields.IBMAPIKey
	}

	return newCredentialDisplayFieldsMap
}

func resourceIBMSccPostureCredentialsUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	postureManagementClient, err := meta.(conns.ClientSession).PostureManagementV2()
	if err != nil {
		return diag.FromErr(err)
	}

	updateCredentialOptions := &posturemanagementv2.UpdateCredentialOptions{}

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error getting userDetails %s", err))
	}
	updateCredentialOptions.SetAccountID(userDetails.UserAccount)

	updateCredentialOptions.SetID(d.Id())

	updateCredentialOptions.SetEnabled(d.Get("enabled").(bool))

	updateCredentialOptions.SetType(d.Get("type").(string))

	updateCredentialOptions.SetName(d.Get("name").(string))

	updateCredentialOptions.SetDescription(d.Get("description").(string))

	updateCredentialDisplayFieldsModel := &posturemanagementv2.UpdateCredentialDisplayFields{
		IBMAPIKey: core.StringPtr("sample_api_key"),
	}
	//displayFields := resourceIBMSccPostureV2CredentialsMapToUpdateCredentialDisplayFields(d.Get("display_fields.0").(map[string]interface{}))
	updateCredentialOptions.SetDisplayFields(updateCredentialDisplayFieldsModel)

	updateCredentialOptions.SetPurpose(d.Get("purpose").(string))

	_, response, err := postureManagementClient.UpdateCredentialWithContext(context, updateCredentialOptions)
	if err != nil {
		log.Printf("[DEBUG] UpdateCredentialWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("UpdateCredentialWithContext failed %s\n%s", err, response))
	}

	return resourceIBMSccPostureCredentialsRead(context, d, meta)
}

func resourceIBMSccPostureCredentialsDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	postureManagementClient, err := meta.(conns.ClientSession).PostureManagementV2()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteCredentialOptions := &posturemanagementv2.DeleteCredentialOptions{}

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error getting userDetails %s", err))
	}
	deleteCredentialOptions.SetAccountID(userDetails.UserAccount)

	deleteCredentialOptions.SetID(d.Id())

	response, err := postureManagementClient.DeleteCredentialWithContext(context, deleteCredentialOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteCredentialWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteCredentialWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}
