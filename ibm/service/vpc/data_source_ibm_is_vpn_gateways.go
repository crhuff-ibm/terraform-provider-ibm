// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isvpnGateways            = "vpn_gateways"
	isVPNGatewayResourceType = "resource_type"
	isVPNGatewayCrn          = "crn"
)

func DataSourceIBMISVPNGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMVPNGatewaysRead,

		Schema: map[string]*schema.Schema{

			isvpnGateways: {
				Type:        schema.TypeList,
				Description: "Collection of VPN Gateways",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isVPNGatewayName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPN Gateway instance name",
						},
						isVPNGatewayCreatedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that this VPN gateway was created",
						},
						isVPNGatewayCrn: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VPN gateway's CRN",
						},
						isVPNGatewayMembers: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Collection of VPN gateway members",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The public IP address assigned to the VPN gateway member",
									},

									"private_address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The private IP address assigned to the VPN gateway member",
									},
									"role": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The high availability role assigned to the VPN gateway member",
									},

									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The status of the VPN gateway member",
									},
								},
							},
						},

						isVPNGatewayResourceType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},

						isVPNGatewayStatus: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the VPN gateway",
						},

						isVPNGatewaySubnet: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPNGateway subnet info",
						},
						isVPNGatewayResourceGroup: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "resource group identifiers ",
						},
						isVPNGatewayMode: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: " VPN gateway mode(policy/route) ",
						},
						"vpc": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "VPC for the VPN Gateway",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this VPC.",
									},
									"deleted": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"more_info": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link to documentation about deleted resources.",
												},
											},
										},
									},
									"href": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this VPC.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this VPC.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique user-defined name for this VPC.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMVPNGatewaysRead(d *schema.ResourceData, meta interface{}) error {

	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	listvpnGWOptions := sess.NewListVPNGatewaysOptions()

	start := ""
	allrecs := []vpcv1.VPNGatewayIntf{}
	for {
		if start != "" {
			listvpnGWOptions.Start = &start
		}
		availableVPNGateways, detail, err := sess.ListVPNGateways(listvpnGWOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error reading list of VPN Gateways:%s\n%s", err, detail)
		}
		start = flex.GetNext(availableVPNGateways.Next)
		allrecs = append(allrecs, availableVPNGateways.VPNGateways...)
		if start == "" {
			break
		}
	}

	vpngateways := make([]map[string]interface{}, 0)
	for _, instance := range allrecs {
		gateway := map[string]interface{}{}
		data := instance.(*vpcv1.VPNGateway)
		gateway[isVPNGatewayName] = *data.Name
		gateway[isVPNGatewayCreatedAt] = data.CreatedAt.String()
		gateway[isVPNGatewayResourceType] = *data.ResourceType
		gateway[isVPNGatewayStatus] = *data.Status
		gateway[isVPNGatewayMode] = *data.Mode
		gateway[isVPNGatewayResourceGroup] = *data.ResourceGroup.ID
		gateway[isVPNGatewaySubnet] = *data.Subnet.ID
		gateway[isVPNGatewayCrn] = *data.CRN

		if data.Members != nil {
			vpcMembersIpsList := make([]map[string]interface{}, 0)
			for _, memberIP := range data.Members {
				currentMemberIP := map[string]interface{}{}
				if memberIP.PublicIP != nil {
					currentMemberIP["address"] = *memberIP.PublicIP.Address
					currentMemberIP["role"] = *memberIP.Role
					currentMemberIP["status"] = *memberIP.Status
					vpcMembersIpsList = append(vpcMembersIpsList, currentMemberIP)
				}
				if memberIP.PrivateIP != nil && memberIP.PrivateIP.Address != nil {
					currentMemberIP["private_address"] = *memberIP.PrivateIP.Address
				}
			}
			gateway[isVPNGatewayMembers] = vpcMembersIpsList
		}

		if data.VPC != nil {
			vpcList := []map[string]interface{}{}
			vpcList = append(vpcList, dataSourceVPNServerCollectionVPNGatewayVpcReferenceToMap(data.VPC))
			gateway["vpc"] = vpcList
		}

		vpngateways = append(vpngateways, gateway)
	}

	d.SetId(dataSourceIBMVPNGatewaysID(d))
	d.Set(isvpnGateways, vpngateways)
	return nil
}

// dataSourceIBMVPNGatewaysID returns a reasonable ID  list.
func dataSourceIBMVPNGatewaysID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceVPNServerCollectionVPNGatewayVpcReferenceToMap(vpcsItem *vpcv1.VPCReference) (vpcsMap map[string]interface{}) {
	vpcsMap = map[string]interface{}{}

	if vpcsItem.CRN != nil {
		vpcsMap["crn"] = vpcsItem.CRN
	}
	if vpcsItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceVPNGatewayCollectionVpcsDeletedToMap(*vpcsItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		vpcsMap["deleted"] = deletedList
	}
	if vpcsItem.Href != nil {
		vpcsMap["href"] = vpcsItem.Href
	}
	if vpcsItem.ID != nil {
		vpcsMap["id"] = vpcsItem.ID
	}
	if vpcsItem.Name != nil {
		vpcsMap["name"] = vpcsItem.Name
	}

	return vpcsMap
}

func dataSourceVPNGatewayCollectionVpcsDeletedToMap(deletedItem vpcv1.VPCReferenceDeleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}
