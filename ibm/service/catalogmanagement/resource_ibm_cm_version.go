// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package catalogmanagement

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/catalogmanagementv1"
)

func ResourceIBMCmVersion() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMCmVersionCreate,
		ReadContext:   resourceIBMCmVersionRead,
		UpdateContext: resourceIBMCmVersionUpdate,
		DeleteContext: resourceIBMCmVersionDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"catalog_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Catalog identifier.",
			},
			"offering_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Offering identification.",
			},
			"tags": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Tags array.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"content": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Byte array representing the content to be imported. Only supported for OVA images at this time.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Name of version. Required for virtual server image for VPC.",
			},
			"label": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Display name of version. Required for virtual server image for VPC.",
			},
			"install_kind": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Install type. Example: instance. Required for virtual server image for VPC.",
			},
			"target_kinds": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Deployment target of the content being onboarded. Current valid values are iks, roks, vcenter, power-iaas, terraform, and vpc-x86. Required for virtual server image for VPC.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"format_kind": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Format of content being onboarded. Example: vsi-image. Required for virtual server image for VPC.",
			},
			"product_kind": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional product kind for the software being onboarded.  Valid values are software, module, or solution.  Default value is software.",
			},
			"import_sha": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SHA256 fingerprint of the image file. Required for virtual server image for VPC.",
			},
			"sha": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SHA256 fingerprint of the image file. Required for virtual server image for VPC.",
			},
			"version": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Semantic version of the software being onboarded. Required for virtual server image for VPC.",
			},
			"flavor": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Version Flavor Information.  Only supported for Product kind Solution.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Programmatic name for this flavor.",
						},
						"label": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Label for this flavor.",
						},
						"label_i18n": &schema.Schema{
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "A map of translated strings, by language code.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"index": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Order that this flavor should appear when listed for a single version.",
						},
					},
				},
			},
			"import_metadata": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Generic data to be included with content being onboarded. Required for virtual server image for VPC.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"operating_system": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Operating system included in this image. Required for virtual server image for VPC.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dedicated_host_only": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Images with this operating system can only be used on dedicated hosts or dedicated host groups. Required for virtual server image for VPC.",
									},
									"vendor": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Vendor of the operating system. Required for virtual server image for VPC.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Globally unique name for this operating system Required for virtual server image for VPC.",
									},
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "URL for this operating system. Required for virtual server image for VPC.",
									},
									"display_name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Unique, display-friendly name for the operating system. Required for virtual server image for VPC.",
									},
									"family": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Software family for this operating system. Required for virtual server image for VPC.",
									},
									"version": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Major release version of this operating system. Required for virtual server image for VPC.",
									},
									"architecture": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Operating system architecture. Required for virtual server image for VPC.",
									},
								},
							},
						},
						"file": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Details for the stored image file. Required for virtual server image for VPC.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"size": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Size of the stored image file rounded up to the next gigabyte. Required for virtual server image for VPC.",
									},
								},
							},
						},
						"minimum_provisioned_size": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Minimum size (in gigabytes) of a volume onto which this image may be provisioned. Required for virtual server image for VPC.",
						},
						"images": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Image operating system. Required for virtual server image for VPC.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Programmatic ID of virtual server image. Required for virtual server image for VPC.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Programmatic name of virtual server image. Required for virtual server image for VPC.",
									},
									"region": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Region the virtual server image is available in. Required for virtual server image for VPC.",
									},
								},
							},
						},
					},
				},
			},
			"metadata": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Generic data to be included with content being onboarded. Required for virtual server image for VPC.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_url": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Version source URL.",
						},
						"version_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Version name.",
						},
						"terraform_version": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Terraform version.",
						},
						"validated_terraform_version": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Validated terraform version.",
						},
						"vsi_vpc": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "VSI VPC version information",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operating_system": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Operating system included in this image. Required for virtual server image for VPC.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"dedicated_host_only": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Images with this operating system can only be used on dedicated hosts or dedicated host groups. Required for virtual server image for VPC.",
												},
												"vendor": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Vendor of the operating system. Required for virtual server image for VPC.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Globally unique name for this operating system Required for virtual server image for VPC.",
												},
												"href": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "URL for this operating system. Required for virtual server image for VPC.",
												},
												"display_name": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Unique, display-friendly name for the operating system. Required for virtual server image for VPC.",
												},
												"family": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Software family for this operating system. Required for virtual server image for VPC.",
												},
												"version": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Major release version of this operating system. Required for virtual server image for VPC.",
												},
												"architecture": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Operating system architecture. Required for virtual server image for VPC.",
												},
											},
										},
									},
									"file": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Details for the stored image file. Required for virtual server image for VPC.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"size": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Size of the stored image file rounded up to the next gigabyte. Required for virtual server image for VPC.",
												},
											},
										},
									},
									"minimum_provisioned_size": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Minimum size (in gigabytes) of a volume onto which this image may be provisioned. Required for virtual server image for VPC.",
									},
									"images": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Image operating system. Required for virtual server image for VPC.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Programmatic ID of virtual server image. Required for virtual server image for VPC.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Programmatic name of virtual server image. Required for virtual server image for VPC.",
												},
												"region": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Region the virtual server image is available in. Required for virtual server image for VPC.",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"working_directory": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional - The sub-folder within the specified tgz file that contains the software being onboarded.",
			},
			"zipurl": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "URL path to zip location.  If not specified, must provide content in the body of this call.",
			},
			"target_version": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The semver value for this new version, if not found in the zip url package content.",
			},
			"include_config": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Add all possible configuration values to this version when importing.",
			},
			"is_vsi": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Indicates that the current terraform template is used to install a virtual server image.",
			},
			"repotype": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The type of repository containing this version.  Valid values are 'public_git' or 'enterprise_git'.",
			},
			"x_auth_token": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Authentication token used to access the specified zip file.",
			},
			"rev": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloudant revision.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Version's CRN.",
			},
			"created": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time this version was created.",
			},
			"updated": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time this version was last updated.",
			},
			"kind_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Kind ID.",
			},
			"repo_url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Content's repo URL.",
			},
			"source_url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Content's source URL (e.g git repo).",
			},
			"tgz_url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "File used to on-board this version.",
			},
			"configuration": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of user solicited overrides.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Configuration key.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Value type (string, boolean, int).",
						},
						// "default_value": &schema.Schema{
						// 	Type:        schema.TypeMap,
						// 	Optional:    true,
						// 	Description: "The default value.  To use a secret when the type is password, specify a JSON encoded value of $ref:#/components/schemas/SecretInstance, prefixed with `cmsm_v1:`.",
						// },
						"display_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Display name for configuration type.",
						},
						"value_constraint": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Constraint associated with value, e.g., for string type - regx:[a-z].",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Key description.",
						},
						"required": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Is key required to install.",
						},
						"options": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of options of type.",
							Elem:        &schema.Schema{Type: schema.TypeMap},
						},
						"hidden": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Hide values.",
						},
						"custom_config": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Render type.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "ID of the widget type.",
									},
									"grouping": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Determines where this configuration type is rendered (3 sections today - Target, Resource, and Deployment).",
									},
									"original_grouping": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Original grouping type for this configuration (3 types - Target, Resource, and Deployment).",
									},
									"grouping_index": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Determines the order that this configuration item shows in that particular grouping.",
									},
									"config_constraints": &schema.Schema{
										Type:        schema.TypeMap,
										Optional:    true,
										Description: "Map of constraint parameters that will be passed to the custom widget.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"associations": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "List of parameters that are associated with this configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"parameters": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Parameters for this association.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Name of this parameter.",
															},
															"options_refresh": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Refresh options.",
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"type_metadata": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The original type, as found in the source being onboarded.",
						},
					},
				},
			},
			"outputs": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of output values for this version.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Output key.",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Output description.",
						},
					},
				},
			},
			"iam_permissions": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of IAM permissions that are required to consume this version.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Service name.",
						},
						"role_crns": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Role CRNs for this permission.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"resources": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Resources for this permission.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Resource name.",
									},
									"description": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Resource description.",
									},
									"role_crns": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Role CRNs for this permission.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"validation": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Validation response.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"validated": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Date and time of last successful validation.",
						},
						"requested": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Date and time of last validation was requested.",
						},
						"state": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Current validation state - <empty>, in_progress, valid, invalid, expired.",
						},
						"last_operation": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Last operation (e.g. submit_deployment, generate_installer, install_offering.",
						},
						"target": &schema.Schema{
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "Validation target information (e.g. cluster_id, region, namespace, etc).  Values will vary by Content type.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"message": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Any message needing to be conveyed as part of the validation job.",
						},
					},
				},
			},
			"required_resources": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Resource requirments for installation.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Type of requirement.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "mem, disk, cores, and nodes can be parsed as an int.  targetVersion will be a semver range value.",
						},
					},
				},
			},
			"single_instance": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Denotes if single instance can be deployed to a given cluster.",
			},
			"install": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Script information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instructions": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Instruction on step and by whom (role) that are needed to take place to prepare the target for installing this version.",
						},
						"instructions_i18n": &schema.Schema{
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "A map of translated strings, by language code.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"script": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Optional script that needs to be run post any pre-condition script.",
						},
						"script_permission": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Optional iam permissions that are required on the target cluster to run this script.",
						},
						"delete_script": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Optional script that if run will remove the installed version.",
						},
						"scope": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Optional value indicating if this script is scoped to a namespace or the entire cluster.",
						},
					},
				},
			},
			"pre_install": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional pre-install instructions.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instructions": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Instruction on step and by whom (role) that are needed to take place to prepare the target for installing this version.",
						},
						"instructions_i18n": &schema.Schema{
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "A map of translated strings, by language code.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"script": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Optional script that needs to be run post any pre-condition script.",
						},
						"script_permission": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Optional iam permissions that are required on the target cluster to run this script.",
						},
						"delete_script": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Optional script that if run will remove the installed version.",
						},
						"scope": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Optional value indicating if this script is scoped to a namespace or the entire cluster.",
						},
					},
				},
			},
			"entitlement": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Entitlement license info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"provider_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Provider name.",
						},
						"provider_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Provider ID.",
						},
						"product_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Product ID.",
						},
						"part_numbers": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "list of license entitlement part numbers, eg. D1YGZLL,D1ZXILL.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"image_repo_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Image repository name.",
						},
					},
				},
			},
			"licenses": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of licenses the product was built with.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "License ID.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "license name.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "type of license e.g., Apache xxx.",
						},
						"url": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "URL for the license text.",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "License description.",
						},
					},
				},
			},
			"image_manifest_url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "If set, denotes a url to a YAML file with list of container images used by this version.",
			},
			"deprecated": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "read only field, indicating if this version is deprecated.",
			},
			"package_version": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Version of the package used to create this version.",
			},
			"state": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Offering state.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"current": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "one of: new, validated, account-published, ibm-published, public-published.",
						},
						"current_entered": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Date and time of current request.",
						},
						"pending": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "one of: new, validated, account-published, ibm-published, public-published.",
						},
						"pending_requested": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Date and time of pending request.",
						},
						"previous": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "one of: new, validated, account-published, ibm-published, public-published.",
						},
					},
				},
			},
			"version_locator": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A dotted value of `catalogID`.`versionID`.",
			},
			"long_description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Long description for version.",
			},
			"long_description_i18n": &schema.Schema{
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "A map of translated strings, by language code.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"image_pull_key_name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the image pull key to use from Offering.ImagePullKeys.",
			},
			"solution_info": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Version Solution Information.  Only supported for Product kind Solution.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"architecture_diagrams": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Architecture diagrams for this solution.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"diagram": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Offering Media information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"url": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "URL of the specified media item.",
												},
												"api_url": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "CM API specific URL of the specified media item.",
												},
												"url_proxy": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Offering URL proxy information.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"url": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "URL of the specified media item being proxied.",
															},
															"sha": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "SHA256 fingerprint of image.",
															},
														},
													},
												},
												"caption": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Caption for this media item.",
												},
												"caption_i18n": &schema.Schema{
													Type:        schema.TypeMap,
													Optional:    true,
													Description: "A map of translated strings, by language code.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"type": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Type of this media item.",
												},
												"thumbnail_url": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Thumbnail URL for this media item.",
												},
											},
										},
									},
									"description": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Description of this diagram.",
									},
									"description_i18n": &schema.Schema{
										Type:        schema.TypeMap,
										Optional:    true,
										Description: "A map of translated strings, by language code.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"features": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Features - titles only.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"title": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Heading.",
									},
									"title_i18n": &schema.Schema{
										Type:        schema.TypeMap,
										Optional:    true,
										Description: "A map of translated strings, by language code.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"description": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Feature description.",
									},
									"description_i18n": &schema.Schema{
										Type:        schema.TypeMap,
										Optional:    true,
										Description: "A map of translated strings, by language code.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"cost_estimate": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Cost estimate definition.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"version": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Cost estimate version.",
									},
									"currency": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Cost estimate currency.",
									},
									"projects": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Cost estimate projects.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Project name.",
												},
												"metadata": &schema.Schema{
													Type:        schema.TypeMap,
													Optional:    true,
													Description: "Project metadata.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"past_breakdown": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Cost breakdown definition.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"total_hourly_cost": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Total hourly cost.",
															},
															"total_monthly_c_ost": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Total monthly cost.",
															},
															"resources": &schema.Schema{
																Type:        schema.TypeList,
																Optional:    true,
																Description: "Resources.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"name": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Resource name.",
																		},
																		"metadata": &schema.Schema{
																			Type:        schema.TypeMap,
																			Optional:    true,
																			Description: "Resource metadata.",
																			Elem:        &schema.Schema{Type: schema.TypeString},
																		},
																		"hourly_cost": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Hourly cost.",
																		},
																		"monthly_cost": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Monthly cost.",
																		},
																		"cost_components": &schema.Schema{
																			Type:        schema.TypeList,
																			Optional:    true,
																			Description: "Cost components.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"name": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Cost component name.",
																					},
																					"unit": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Cost component unit.",
																					},
																					"hourly_quantity": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Cost component hourly quantity.",
																					},
																					"monthly_quantity": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Cost component monthly quantity.",
																					},
																					"price": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Cost component price.",
																					},
																					"hourly_cost": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Cost component hourly cost.",
																					},
																					"monthly_cost": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Cost component monthly cist.",
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
												"breakdown": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Cost breakdown definition.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"total_hourly_cost": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Total hourly cost.",
															},
															"total_monthly_c_ost": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Total monthly cost.",
															},
															"resources": &schema.Schema{
																Type:        schema.TypeList,
																Optional:    true,
																Description: "Resources.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"name": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Resource name.",
																		},
																		"metadata": &schema.Schema{
																			Type:        schema.TypeMap,
																			Optional:    true,
																			Description: "Resource metadata.",
																			Elem:        &schema.Schema{Type: schema.TypeString},
																		},
																		"hourly_cost": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Hourly cost.",
																		},
																		"monthly_cost": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Monthly cost.",
																		},
																		"cost_components": &schema.Schema{
																			Type:        schema.TypeList,
																			Optional:    true,
																			Description: "Cost components.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"name": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Cost component name.",
																					},
																					"unit": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Cost component unit.",
																					},
																					"hourly_quantity": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Cost component hourly quantity.",
																					},
																					"monthly_quantity": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Cost component monthly quantity.",
																					},
																					"price": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Cost component price.",
																					},
																					"hourly_cost": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Cost component hourly cost.",
																					},
																					"monthly_cost": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Cost component monthly cist.",
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
												"diff": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Cost breakdown definition.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"total_hourly_cost": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Total hourly cost.",
															},
															"total_monthly_c_ost": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Total monthly cost.",
															},
															"resources": &schema.Schema{
																Type:        schema.TypeList,
																Optional:    true,
																Description: "Resources.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"name": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Resource name.",
																		},
																		"metadata": &schema.Schema{
																			Type:        schema.TypeMap,
																			Optional:    true,
																			Description: "Resource metadata.",
																			Elem:        &schema.Schema{Type: schema.TypeString},
																		},
																		"hourly_cost": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Hourly cost.",
																		},
																		"monthly_cost": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Monthly cost.",
																		},
																		"cost_components": &schema.Schema{
																			Type:        schema.TypeList,
																			Optional:    true,
																			Description: "Cost components.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"name": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Cost component name.",
																					},
																					"unit": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Cost component unit.",
																					},
																					"hourly_quantity": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Cost component hourly quantity.",
																					},
																					"monthly_quantity": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Cost component monthly quantity.",
																					},
																					"price": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Cost component price.",
																					},
																					"hourly_cost": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Cost component hourly cost.",
																					},
																					"monthly_cost": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Cost component monthly cist.",
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
												"summary": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Cost summary definition.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"total_detected_resources": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Total detected resources.",
															},
															"total_supported_resources": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Total supported resources.",
															},
															"total_unsupported_resources": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Total unsupported resources.",
															},
															"total_usage_based_resources": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Total usage based resources.",
															},
															"total_no_price_resources": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Total no price resources.",
															},
															"unsupported_resource_counts": &schema.Schema{
																Type:        schema.TypeMap,
																Optional:    true,
																Description: "Unsupported resource counts.",
																Elem:        &schema.Schema{Type: schema.TypeString},
															},
															"no_price_resource_counts": &schema.Schema{
																Type:        schema.TypeMap,
																Optional:    true,
																Description: "No price resource counts.",
																Elem:        &schema.Schema{Type: schema.TypeString},
															},
														},
													},
												},
											},
										},
									},
									"summary": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Cost summary definition.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"total_detected_resources": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Total detected resources.",
												},
												"total_supported_resources": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Total supported resources.",
												},
												"total_unsupported_resources": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Total unsupported resources.",
												},
												"total_usage_based_resources": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Total usage based resources.",
												},
												"total_no_price_resources": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Total no price resources.",
												},
												"unsupported_resource_counts": &schema.Schema{
													Type:        schema.TypeMap,
													Optional:    true,
													Description: "Unsupported resource counts.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"no_price_resource_counts": &schema.Schema{
													Type:        schema.TypeMap,
													Optional:    true,
													Description: "No price resource counts.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"total_hourly_cost": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Total hourly cost.",
									},
									"total_monthly_cost": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Total monthly cost.",
									},
									"past_total_hourly_cost": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Past total hourly cost.",
									},
									"past_total_monthly_cost": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Past total monthly cost.",
									},
									"diff_total_hourly_cost": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Difference in total hourly cost.",
									},
									"diff_total_monthly_cost": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Difference in total monthly cost.",
									},
									"time_generated": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "When this estimate was generated.",
									},
								},
							},
						},
						"dependencies": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Dependencies for this solution.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"catalog_id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Optional - If not specified, assumes the Public Catalog.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Optional - Offering ID - not required if name is set.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Optional - Programmatic Offering name.",
									},
									"version": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Required - Semver value or range.",
									},
									"flavors": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Optional - List of dependent flavors in the specified range.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"is_consumable": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Is the version able to be shared.",
			},
			"version_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique ID.",
			},
			"offering_identifier": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Offering ID, in the format of <account_id>:o:<offering_id>.",
			},
		},
	}
}

func resourceIBMCmVersionCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	catalogManagementClient, err := meta.(conns.ClientSession).CatalogManagementV1()
	if err != nil {
		return diag.FromErr(err)
	}

	importOfferingVersionOptions := &catalogmanagementv1.ImportOfferingVersionOptions{}

	importOfferingVersionOptions.SetCatalogIdentifier(d.Get("catalog_id").(string))
	importOfferingVersionOptions.SetOfferingID(d.Get("offering_id").(string))
	if _, ok := d.GetOk("tags"); ok {
		importOfferingVersionOptions.SetTags(SIToSS(d.Get("tags").([]interface{})))
	}
	if _, ok := d.GetOk("content"); ok {
		importOfferingVersionOptions.SetContent([]byte(d.Get("content").(string)))
	}
	if _, ok := d.GetOk("name"); ok {
		importOfferingVersionOptions.SetName(d.Get("name").(string))
	}
	if _, ok := d.GetOk("label"); ok {
		importOfferingVersionOptions.SetLabel(d.Get("label").(string))
	}
	if _, ok := d.GetOk("install_kind"); ok {
		importOfferingVersionOptions.SetInstallKind(d.Get("install_kind").(string))
	}
	if _, ok := d.GetOk("target_kinds"); ok {
		importOfferingVersionOptions.SetTargetKinds(SIToSS(d.Get("target_kinds").([]interface{})))
	}
	if _, ok := d.GetOk("format_kind"); ok {
		importOfferingVersionOptions.SetFormatKind(d.Get("format_kind").(string))
	}
	if _, ok := d.GetOk("product_kind"); ok {
		importOfferingVersionOptions.SetProductKind(d.Get("product_kind").(string))
	}
	if _, ok := d.GetOk("import_sha"); ok {
		importOfferingVersionOptions.SetSha(d.Get("import_sha").(string))
	}
	if _, ok := d.GetOk("flavor"); ok {
		flavorModel, err := resourceIBMCmVersionMapToFlavor(d.Get("flavor.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		importOfferingVersionOptions.SetFlavor(flavorModel)
	}
	if _, ok := d.GetOk("import_metadata"); ok {
		metadataModel, err := resourceIBMCmVersionMapToImportOfferingBodyMetadata(d.Get("import_metadata.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		importOfferingVersionOptions.SetMetadata(metadataModel)
	}
	if _, ok := d.GetOk("working_directory"); ok {
		importOfferingVersionOptions.SetWorkingDirectory(d.Get("working_directory").(string))
	}
	if _, ok := d.GetOk("zipurl"); ok {
		importOfferingVersionOptions.SetZipurl(d.Get("zipurl").(string))
	}
	if _, ok := d.GetOk("target_version"); ok {
		importOfferingVersionOptions.SetTargetVersion(d.Get("target_version").(string))
		importOfferingVersionOptions.SetVersion(d.Get("target_version").(string))
	}
	if _, ok := d.GetOk("include_config"); ok {
		importOfferingVersionOptions.SetIncludeConfig(d.Get("include_config").(bool))
	}
	if _, ok := d.GetOk("is_vsi"); ok {
		importOfferingVersionOptions.SetIsVsi(d.Get("is_vsi").(bool))
	}
	if _, ok := d.GetOk("repotype"); ok {
		importOfferingVersionOptions.SetRepotype(d.Get("repotype").(string))
	}
	if _, ok := d.GetOk("x_auth_token"); ok {
		importOfferingVersionOptions.SetXAuthToken(d.Get("x_auth_token").(string))
	}

	getOfferingOptions := &catalogmanagementv1.GetOfferingOptions{}
	getOfferingOptions.SetCatalogIdentifier(d.Get("catalog_id").(string))
	getOfferingOptions.SetOfferingID(d.Get("offering_id").(string))
	oldOffering, response, err := catalogManagementClient.GetOfferingWithContext(context, getOfferingOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetOfferingWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetOfferingWithContext failed %s\n%s", err, response))
	}

	offering, response, err := catalogManagementClient.ImportOfferingVersionWithContext(context, importOfferingVersionOptions)
	if err != nil {
		log.Printf("[DEBUG] ImportOfferingVersionWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("ImportOfferingVersionWithContext failed %s\n%s", err, response))
	}

	activeVersionID, err := getVersionFromOffering(*oldOffering, *offering)
	if err != nil {
		log.Printf("[DEBUG] getVersionFromOffering failed %s\n", err)
		return diag.FromErr(fmt.Errorf("getVersionFromOffering failed %s\n", err))
	}

	var activeVersion catalogmanagementv1.Version
	for _, k := range offering.Kinds {
		for _, v := range k.Versions {
			if v.ID == &activeVersionID {
				activeVersion = v
			}
		}
	}

	d.SetId(fmt.Sprintf("%s/%s", *offering.CatalogID, activeVersionID))

	updateOfferingOptions := &catalogmanagementv1.UpdateOfferingOptions{}

	updateOfferingOptions.SetCatalogIdentifier(*offering.CatalogID)
	updateOfferingOptions.SetOfferingID(*offering.ID)
	ifMatch := fmt.Sprintf("\"%s\"", *offering.Rev)
	updateOfferingOptions.IfMatch = &ifMatch

	hasChange := false
	method := "replace"

	// find kind and version index
	var kindIndex int
	var versionIndex int
	for i, kind := range offering.Kinds {
		if kind.ID == activeVersion.KindID {
			kindIndex = i

			if kind.Versions != nil && len(kind.Versions) > 0 {
				for j, version := range kind.Versions {
					if version.ID == activeVersion.ID {
						versionIndex = j
					}
				}
			}
		}
	}

	pathToVersion := fmt.Sprintf("/kinds/%d/versions/%d", kindIndex, versionIndex)

	if _, ok := d.GetOk("flavor"); ok {
		if activeVersion.Flavor == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := fmt.Sprintf("%s/flavor", pathToVersion)
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("flavor.0"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if _, ok := d.GetOk("tags"); ok {
		if activeVersion.Tags == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := fmt.Sprintf("%s/tags", pathToVersion)
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("tags"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if _, ok := d.GetOk("configuration"); ok {
		if activeVersion.Configuration == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := fmt.Sprintf("%s/configuration", pathToVersion)
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("configuration"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if _, ok := d.GetOk("iam_permissions"); ok {
		if activeVersion.IamPermissions == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := fmt.Sprintf("%s/iam_permissions", pathToVersion)
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("iam_permissions"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if _, ok := d.GetOk("install"); ok {
		if activeVersion.Install == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := fmt.Sprintf("%s/install", pathToVersion)
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("install.0"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if _, ok := d.GetOk("pre_install"); ok {
		if activeVersion.PreInstall == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := fmt.Sprintf("%s/pre_install", pathToVersion)
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("pre_install"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if _, ok := d.GetOk("licenses"); ok {
		if activeVersion.Licenses == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := fmt.Sprintf("%s/licenses", pathToVersion)
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("licenses"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if _, ok := d.GetOk("solution_info"); ok {
		if activeVersion.SolutionInfo == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := fmt.Sprintf("%s/solution_info", pathToVersion)
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("solution_info"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if _, ok := d.GetOk("is_consumable"); ok {
		if activeVersion.IsConsumable == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := fmt.Sprintf("%s/is_consumable", pathToVersion)
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("is_consumable"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}

	if hasChange {
		_, response, err := catalogManagementClient.UpdateOfferingWithContext(context, updateOfferingOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateOfferingWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateOfferingWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIBMCmVersionRead(context, d, meta)
}

func resourceIBMCmVersionRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	catalogManagementClient, err := meta.(conns.ClientSession).CatalogManagementV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getVersionOptions := &catalogmanagementv1.GetVersionOptions{}

	getVersionOptions.SetVersionLocID(strings.Replace(d.Id(), "/", ".", 1))

	offering, response, err := catalogManagementClient.GetVersionWithContext(context, getVersionOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetVersionWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetVersionWithContext failed %s\n%s", err, response))
	}

	version := offering.Kinds[0].Versions[0]

	if err = d.Set("offering_identifier", version.OfferingID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting offering_identifier: %s", err))
	}
	if version.Tags != nil {
		if err = d.Set("tags", version.Tags); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting tags: %s", err))
		}
	}
	if err = d.Set("sha", version.Sha); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting sha: %s", err))
	}
	if err = d.Set("version", version.Version); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting version: %s", err))
	}
	if version.Flavor != nil {
		flavorMap, err := resourceIBMCmVersionFlavorToMap(version.Flavor)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("flavor", []map[string]interface{}{flavorMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting flavor: %s", err))
		}
	}
	metadata := []map[string]interface{}{}
	if version.Metadata != nil {
		var modelMapVSI map[string]interface{}
		if version.Metadata["vsi_vpc"] != nil {
			modelMapVSI, err = dataSourceIBMCmVersionMetadataVSIToMap(version.Metadata["vsi_vpc"].(map[string]interface{}))
			if err != nil {
				return diag.FromErr(err)
			}
		}
		convertedMap := make(map[string]interface{}, len(version.Metadata))
		for k, v := range version.Metadata {
			if k == "vsi_vpc" {
				convertedMap[k] = []map[string]interface{}{modelMapVSI}
			} else {
				convertedMap[k] = v
			}
		}
		metadata = append(metadata, convertedMap)
	}
	if err = d.Set("metadata", metadata); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting metadata %s", err))
	}
	if err = d.Set("rev", version.Rev); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting rev: %s", err))
	}
	if err = d.Set("crn", version.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}
	if err = d.Set("created", flex.DateTimeToString(version.Created)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created: %s", err))
	}
	if err = d.Set("updated", flex.DateTimeToString(version.Updated)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated: %s", err))
	}
	if err = d.Set("catalog_id", version.CatalogID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting catalog_id: %s", err))
	}
	if err = d.Set("kind_id", version.KindID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting kind_id: %s", err))
	}
	if err = d.Set("repo_url", version.RepoURL); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting repo_url: %s", err))
	}
	if err = d.Set("source_url", version.SourceURL); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting source_url: %s", err))
	}
	if err = d.Set("tgz_url", version.TgzURL); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting tgz_url: %s", err))
	}
	configuration := []map[string]interface{}{}
	if version.Configuration != nil {
		for _, configurationItem := range version.Configuration {
			configurationItemMap, err := resourceIBMCmVersionConfigurationToMap(&configurationItem)
			if err != nil {
				return diag.FromErr(err)
			}
			configuration = append(configuration, configurationItemMap)
		}
	}
	if err = d.Set("configuration", configuration); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting configuration: %s", err))
	}
	outputs := []map[string]interface{}{}
	if version.Outputs != nil {
		for _, outputsItem := range version.Outputs {
			outputsItemMap, err := resourceIBMCmVersionOutputToMap(&outputsItem)
			if err != nil {
				return diag.FromErr(err)
			}
			outputs = append(outputs, outputsItemMap)
		}
	}
	if err = d.Set("outputs", outputs); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting outputs: %s", err))
	}
	iamPermissions := []map[string]interface{}{}
	if version.IamPermissions != nil {
		for _, iamPermissionsItem := range version.IamPermissions {
			iamPermissionsItemMap, err := resourceIBMCmVersionIamPermissionToMap(&iamPermissionsItem)
			if err != nil {
				return diag.FromErr(err)
			}
			iamPermissions = append(iamPermissions, iamPermissionsItemMap)
		}
	}
	if err = d.Set("iam_permissions", iamPermissions); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting iam_permissions: %s", err))
	}
	if version.Validation != nil {
		validationMap, err := resourceIBMCmVersionValidationToMap(version.Validation)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("validation", []map[string]interface{}{validationMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting validation: %s", err))
		}
	}
	requiredResources := []map[string]interface{}{}
	if version.RequiredResources != nil {
		for _, requiredResourcesItem := range version.RequiredResources {
			requiredResourcesItemMap, err := resourceIBMCmVersionResourceToMap(&requiredResourcesItem)
			if err != nil {
				return diag.FromErr(err)
			}
			requiredResources = append(requiredResources, requiredResourcesItemMap)
		}
	}
	if err = d.Set("required_resources", requiredResources); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting required_resources: %s", err))
	}
	if err = d.Set("single_instance", version.SingleInstance); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting single_instance: %s", err))
	}
	if version.Install != nil {
		installMap, err := resourceIBMCmVersionScriptToMap(version.Install)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("install", []map[string]interface{}{installMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting install: %s", err))
		}
	}
	preInstall := []map[string]interface{}{}
	if version.PreInstall != nil {
		for _, preInstallItem := range version.PreInstall {
			preInstallItemMap, err := resourceIBMCmVersionScriptToMap(&preInstallItem)
			if err != nil {
				return diag.FromErr(err)
			}
			preInstall = append(preInstall, preInstallItemMap)
		}
	}
	if err = d.Set("pre_install", preInstall); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting pre_install: %s", err))
	}
	if version.Entitlement != nil {
		entitlementMap, err := resourceIBMCmVersionVersionEntitlementToMap(version.Entitlement)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("entitlement", []map[string]interface{}{entitlementMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting entitlement: %s", err))
		}
	}
	licenses := []map[string]interface{}{}
	if version.Licenses != nil {
		for _, licensesItem := range version.Licenses {
			licensesItemMap, err := resourceIBMCmVersionLicenseToMap(&licensesItem)
			if err != nil {
				return diag.FromErr(err)
			}
			licenses = append(licenses, licensesItemMap)
		}
	}
	if err = d.Set("licenses", licenses); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting licenses: %s", err))
	}
	if err = d.Set("image_manifest_url", version.ImageManifestURL); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting image_manifest_url: %s", err))
	}
	if err = d.Set("deprecated", version.Deprecated); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting deprecated: %s", err))
	}
	if err = d.Set("package_version", version.PackageVersion); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting package_version: %s", err))
	}
	if version.State != nil {
		stateMap, err := resourceIBMCmVersionStateToMap(version.State)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("state", []map[string]interface{}{stateMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting state: %s", err))
		}
	}
	if err = d.Set("version_locator", version.VersionLocator); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting version_locator: %s", err))
	}
	if err = d.Set("long_description", version.LongDescription); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting long_description: %s", err))
	}
	if err = d.Set("long_description_i18n", version.LongDescriptionI18n); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting long_description_i18n: %s", err))
	}
	if err = d.Set("image_pull_key_name", version.ImagePullKeyName); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting image_pull_key_name: %s", err))
	}
	if version.SolutionInfo != nil {
		solutionInfoMap, err := resourceIBMCmVersionSolutionInfoToMap(version.SolutionInfo)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("solution_info", []map[string]interface{}{solutionInfoMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting solution_info: %s", err))
		}
	}
	if err = d.Set("is_consumable", version.IsConsumable); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting is_consumable: %s", err))
	}
	if err = d.Set("version_id", version.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting version_id: %s", err))
	}

	return nil
}

func resourceIBMCmVersionUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	catalogManagementClient, err := meta.(conns.ClientSession).CatalogManagementV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getVersionOptions := &catalogmanagementv1.GetVersionOptions{}
	getVersionOptions.SetVersionLocID(strings.Replace(d.Id(), "/", ".", 1))

	partialOffering, response, err := catalogManagementClient.GetVersionWithContext(context, getVersionOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetVersionWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetVersionWithContext failed %s\n%s", err, response))
	}
	activeVersion := partialOffering.Kinds[0].Versions[0]

	getOfferingOptions := &catalogmanagementv1.GetOfferingOptions{}

	getOfferingOptions.SetCatalogIdentifier(*partialOffering.CatalogID)
	getOfferingOptions.SetOfferingID(*partialOffering.ID)
	offering, response, err := catalogManagementClient.GetOfferingWithContext(context, getOfferingOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetOfferingWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetOfferingWithContext failed %s\n%s", err, response))
	}

	updateOfferingOptions := &catalogmanagementv1.UpdateOfferingOptions{}

	updateOfferingOptions.SetCatalogIdentifier(*offering.CatalogID)
	updateOfferingOptions.SetOfferingID(*offering.ID)
	ifMatch := fmt.Sprintf("\"%s\"", *offering.Rev)
	updateOfferingOptions.IfMatch = &ifMatch

	hasChange := false
	method := "replace"

	// find kind and version index
	var kindIndex int
	var versionIndex int
	for i, kind := range offering.Kinds {
		if kind.ID == activeVersion.KindID {
			kindIndex = i

			if kind.Versions != nil && len(kind.Versions) > 0 {
				for j, version := range kind.Versions {
					if version.ID == activeVersion.ID {
						versionIndex = j
					}
				}
			}
		}
	}

	pathToVersion := fmt.Sprintf("/kinds/%d/versions/%d", kindIndex, versionIndex)

	if d.HasChange("flavor") {
		if activeVersion.Flavor == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := fmt.Sprintf("%s/flavor", pathToVersion)
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("flavor.0"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if d.HasChange("tags") {
		if activeVersion.Tags == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := fmt.Sprintf("%s/tags", pathToVersion)
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("tags"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if d.HasChange("configuration") {
		if activeVersion.Configuration == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := fmt.Sprintf("%s/configuration", pathToVersion)
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("configuration"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if d.HasChange("iam_permissions") {
		if activeVersion.IamPermissions == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := fmt.Sprintf("%s/iam_permissions", pathToVersion)
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("iam_permissions"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if d.HasChange("install") {
		if activeVersion.Install == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := fmt.Sprintf("%s/install", pathToVersion)
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("install.0"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if d.HasChange("pre_install") {
		if activeVersion.PreInstall == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := fmt.Sprintf("%s/pre_install", pathToVersion)
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("pre_install"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if d.HasChange("licenses") {
		if activeVersion.Licenses == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := fmt.Sprintf("%s/licenses", pathToVersion)
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("licenses"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}
	if d.HasChange("solution_info") {
		if activeVersion.SolutionInfo == nil {
			method = "add"
		} else {
			method = "replace"
		}
		path := fmt.Sprintf("%s/solution_info", pathToVersion)
		update := catalogmanagementv1.JSONPatchOperation{
			Op:    &method,
			Path:  &path,
			Value: d.Get("solution_info"),
		}
		updateOfferingOptions.Updates = append(updateOfferingOptions.Updates, update)
		hasChange = true
	}

	if hasChange {
		_, response, err := catalogManagementClient.UpdateOfferingWithContext(context, updateOfferingOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateOfferingWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateOfferingWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIBMCmVersionRead(context, d, meta)
}

func resourceIBMCmVersionDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	catalogManagementClient, err := meta.(conns.ClientSession).CatalogManagementV1()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteVersionOptions := &catalogmanagementv1.DeleteVersionOptions{}

	deleteVersionOptions.SetVersionLocID(strings.Replace(d.Id(), "/", ".", 1))

	response, err := catalogManagementClient.DeleteVersionWithContext(context, deleteVersionOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteVersionWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteVersionWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func resourceIBMCmVersionMapToFlavor(modelMap map[string]interface{}) (*catalogmanagementv1.Flavor, error) {
	model := &catalogmanagementv1.Flavor{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["label"] != nil && modelMap["label"].(string) != "" {
		model.Label = core.StringPtr(modelMap["label"].(string))
	}
	if modelMap["index"] != nil {
		model.Index = core.Int64Ptr(int64(modelMap["index"].(int)))
	}
	return model, nil
}

func resourceIBMCmVersionMapToImportOfferingBodyMetadata(modelMap map[string]interface{}) (*catalogmanagementv1.ImportOfferingBodyMetadata, error) {
	model := &catalogmanagementv1.ImportOfferingBodyMetadata{}
	if modelMap["operating_system"] != nil && len(modelMap["operating_system"].([]interface{})) > 0 {
		OperatingSystemModel, err := resourceIBMCmVersionMapToImportOfferingBodyMetadataOperatingSystem(modelMap["operating_system"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.OperatingSystem = OperatingSystemModel
	}
	if modelMap["file"] != nil && len(modelMap["file"].([]interface{})) > 0 {
		FileModel, err := resourceIBMCmVersionMapToImportOfferingBodyMetadataFile(modelMap["file"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.File = FileModel
	}
	if modelMap["minimum_provisioned_size"] != nil {
		model.MinimumProvisionedSize = core.Int64Ptr(int64(modelMap["minimum_provisioned_size"].(int)))
	}
	if modelMap["images"] != nil {
		images := []catalogmanagementv1.ImportOfferingBodyMetadataImagesItem{}
		for _, imagesItem := range modelMap["images"].([]interface{}) {
			imagesItemModel, err := resourceIBMCmVersionMapToImportOfferingBodyMetadataImagesItem(imagesItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			images = append(images, *imagesItemModel)
		}
		model.Images = images
	}
	return model, nil
}

func resourceIBMCmVersionMapToImportOfferingBodyMetadataOperatingSystem(modelMap map[string]interface{}) (*catalogmanagementv1.ImportOfferingBodyMetadataOperatingSystem, error) {
	model := &catalogmanagementv1.ImportOfferingBodyMetadataOperatingSystem{}
	if modelMap["dedicated_host_only"] != nil {
		model.DedicatedHostOnly = core.BoolPtr(modelMap["dedicated_host_only"].(bool))
	}
	if modelMap["vendor"] != nil && modelMap["vendor"].(string) != "" {
		model.Vendor = core.StringPtr(modelMap["vendor"].(string))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	if modelMap["display_name"] != nil && modelMap["display_name"].(string) != "" {
		model.DisplayName = core.StringPtr(modelMap["display_name"].(string))
	}
	if modelMap["family"] != nil && modelMap["family"].(string) != "" {
		model.Family = core.StringPtr(modelMap["family"].(string))
	}
	if modelMap["version"] != nil && modelMap["version"].(string) != "" {
		model.Version = core.StringPtr(modelMap["version"].(string))
	}
	if modelMap["architecture"] != nil && modelMap["architecture"].(string) != "" {
		model.Architecture = core.StringPtr(modelMap["architecture"].(string))
	}
	return model, nil
}

func resourceIBMCmVersionMapToImportOfferingBodyMetadataFile(modelMap map[string]interface{}) (*catalogmanagementv1.ImportOfferingBodyMetadataFile, error) {
	model := &catalogmanagementv1.ImportOfferingBodyMetadataFile{}
	if modelMap["size"] != nil {
		model.Size = core.Int64Ptr(int64(modelMap["size"].(int)))
	}
	return model, nil
}

func resourceIBMCmVersionMapToImportOfferingBodyMetadataImagesItem(modelMap map[string]interface{}) (*catalogmanagementv1.ImportOfferingBodyMetadataImagesItem, error) {
	model := &catalogmanagementv1.ImportOfferingBodyMetadataImagesItem{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["region"] != nil && modelMap["region"].(string) != "" {
		model.Region = core.StringPtr(modelMap["region"].(string))
	}
	return model, nil
}

func resourceIBMCmVersionFlavorToMap(model *catalogmanagementv1.Flavor) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Label != nil {
		modelMap["label"] = model.Label
	}
	if model.Index != nil {
		modelMap["index"] = flex.IntValue(model.Index)
	}
	return modelMap, nil
}

func dataSourceIBMCmVersionMetadataVSIToMap(model map[string]interface{}) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})

	for k, v := range model {
		if k == "operating_system" || k == "file" {
			modelMap[k] = []map[string]interface{}{v.(map[string]interface{})}
		} else {
			modelMap[k] = v
		}
	}

	return modelMap, nil
}

func resourceIBMCmVersionImportOfferingBodyMetadataOperatingSystemToMap(model *catalogmanagementv1.ImportOfferingBodyMetadataOperatingSystem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.DedicatedHostOnly != nil {
		modelMap["dedicated_host_only"] = model.DedicatedHostOnly
	}
	if model.Vendor != nil {
		modelMap["vendor"] = model.Vendor
	}
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Href != nil {
		modelMap["href"] = model.Href
	}
	if model.DisplayName != nil {
		modelMap["display_name"] = model.DisplayName
	}
	if model.Family != nil {
		modelMap["family"] = model.Family
	}
	if model.Version != nil {
		modelMap["version"] = model.Version
	}
	if model.Architecture != nil {
		modelMap["architecture"] = model.Architecture
	}
	return modelMap, nil
}

func resourceIBMCmVersionImportOfferingBodyMetadataFileToMap(model *catalogmanagementv1.ImportOfferingBodyMetadataFile) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Size != nil {
		modelMap["size"] = flex.IntValue(model.Size)
	}
	return modelMap, nil
}

func resourceIBMCmVersionImportOfferingBodyMetadataImagesItemToMap(model *catalogmanagementv1.ImportOfferingBodyMetadataImagesItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Region != nil {
		modelMap["region"] = model.Region
	}
	return modelMap, nil
}

func resourceIBMCmVersionConfigurationToMap(model *catalogmanagementv1.Configuration) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Key != nil {
		modelMap["key"] = model.Key
	}
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	// if model.DefaultValue != nil {
	// 	modelMap["default_value"] = model.DefaultValue
	// }
	if model.DisplayName != nil {
		modelMap["display_name"] = model.DisplayName
	}
	if model.ValueConstraint != nil {
		modelMap["value_constraint"] = model.ValueConstraint
	}
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.Required != nil {
		modelMap["required"] = model.Required
	}
	if model.Options != nil {
		options := []map[string]interface{}{}
		for _, optionsItem := range model.Options {
			options = append(options, optionsItem.(map[string]interface{}))
		}
		modelMap["options"] = options
	}
	if model.Hidden != nil {
		modelMap["hidden"] = model.Hidden
	}
	if model.CustomConfig != nil {
		customConfigMap, err := resourceIBMCmVersionRenderTypeToMap(model.CustomConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["custom_config"] = []map[string]interface{}{customConfigMap}
	}
	if model.TypeMetadata != nil {
		modelMap["type_metadata"] = model.TypeMetadata
	}
	return modelMap, nil
}

func resourceIBMCmVersionRenderTypeToMap(model *catalogmanagementv1.RenderType) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	if model.Grouping != nil {
		modelMap["grouping"] = model.Grouping
	}
	if model.OriginalGrouping != nil {
		modelMap["original_grouping"] = model.OriginalGrouping
	}
	if model.GroupingIndex != nil {
		modelMap["grouping_index"] = flex.IntValue(model.GroupingIndex)
	}
	if model.Associations != nil {
		associationsMap, err := resourceIBMCmVersionRenderTypeAssociationsToMap(model.Associations)
		if err != nil {
			return modelMap, err
		}
		modelMap["associations"] = []map[string]interface{}{associationsMap}
	}
	return modelMap, nil
}

func resourceIBMCmVersionRenderTypeAssociationsToMap(model *catalogmanagementv1.RenderTypeAssociations) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Parameters != nil {
		parameters := []map[string]interface{}{}
		for _, parametersItem := range model.Parameters {
			parametersItemMap, err := resourceIBMCmVersionRenderTypeAssociationsParametersItemToMap(&parametersItem)
			if err != nil {
				return modelMap, err
			}
			parameters = append(parameters, parametersItemMap)
		}
		modelMap["parameters"] = parameters
	}
	return modelMap, nil
}

func resourceIBMCmVersionRenderTypeAssociationsParametersItemToMap(model *catalogmanagementv1.RenderTypeAssociationsParametersItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.OptionsRefresh != nil {
		modelMap["options_refresh"] = model.OptionsRefresh
	}
	return modelMap, nil
}

func resourceIBMCmVersionOutputToMap(model *catalogmanagementv1.Output) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Key != nil {
		modelMap["key"] = model.Key
	}
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	return modelMap, nil
}

func resourceIBMCmVersionIamPermissionToMap(model *catalogmanagementv1.IamPermission) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ServiceName != nil {
		modelMap["service_name"] = model.ServiceName
	}
	if model.RoleCrns != nil {
		modelMap["role_crns"] = model.RoleCrns
	}
	if model.Resources != nil {
		resources := []map[string]interface{}{}
		for _, resourcesItem := range model.Resources {
			resourcesItemMap, err := resourceIBMCmVersionIamResourceToMap(&resourcesItem)
			if err != nil {
				return modelMap, err
			}
			resources = append(resources, resourcesItemMap)
		}
		modelMap["resources"] = resources
	}
	return modelMap, nil
}

func resourceIBMCmVersionIamResourceToMap(model *catalogmanagementv1.IamResource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.RoleCrns != nil {
		modelMap["role_crns"] = model.RoleCrns
	}
	return modelMap, nil
}

func resourceIBMCmVersionValidationToMap(model *catalogmanagementv1.Validation) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Validated != nil {
		modelMap["validated"] = model.Validated.String()
	}
	if model.Requested != nil {
		modelMap["requested"] = model.Requested.String()
	}
	if model.State != nil {
		modelMap["state"] = model.State
	}
	if model.LastOperation != nil {
		modelMap["last_operation"] = model.LastOperation
	}
	if model.Message != nil {
		modelMap["message"] = model.Message
	}
	return modelMap, nil
}

func resourceIBMCmVersionResourceToMap(model *catalogmanagementv1.Resource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	if model.Value != nil {
		modelMap["value"] = model.Value
	}
	return modelMap, nil
}

func resourceIBMCmVersionScriptToMap(model *catalogmanagementv1.Script) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Instructions != nil {
		modelMap["instructions"] = model.Instructions
	}
	if model.Script != nil {
		modelMap["script"] = model.Script
	}
	if model.ScriptPermission != nil {
		modelMap["script_permission"] = model.ScriptPermission
	}
	if model.DeleteScript != nil {
		modelMap["delete_script"] = model.DeleteScript
	}
	if model.Scope != nil {
		modelMap["scope"] = model.Scope
	}
	return modelMap, nil
}

func resourceIBMCmVersionVersionEntitlementToMap(model *catalogmanagementv1.VersionEntitlement) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ProviderName != nil {
		modelMap["provider_name"] = model.ProviderName
	}
	if model.ProviderID != nil {
		modelMap["provider_id"] = model.ProviderID
	}
	if model.ProductID != nil {
		modelMap["product_id"] = model.ProductID
	}
	if model.PartNumbers != nil {
		modelMap["part_numbers"] = model.PartNumbers
	}
	if model.ImageRepoName != nil {
		modelMap["image_repo_name"] = model.ImageRepoName
	}
	return modelMap, nil
}

func resourceIBMCmVersionLicenseToMap(model *catalogmanagementv1.License) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	if model.URL != nil {
		modelMap["url"] = model.URL
	}
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	return modelMap, nil
}

func resourceIBMCmVersionStateToMap(model *catalogmanagementv1.State) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Current != nil {
		modelMap["current"] = model.Current
	}
	if model.CurrentEntered != nil {
		modelMap["current_entered"] = model.CurrentEntered.String()
	}
	if model.Pending != nil {
		modelMap["pending"] = model.Pending
	}
	if model.PendingRequested != nil {
		modelMap["pending_requested"] = model.PendingRequested.String()
	}
	if model.Previous != nil {
		modelMap["previous"] = model.Previous
	}
	return modelMap, nil
}

func resourceIBMCmVersionDeprecatePendingToMap(model *catalogmanagementv1.DeprecatePending) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.DeprecateDate != nil {
		modelMap["deprecate_date"] = model.DeprecateDate.String()
	}
	if model.DeprecateState != nil {
		modelMap["deprecate_state"] = model.DeprecateState
	}
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	return modelMap, nil
}

func resourceIBMCmVersionSolutionInfoToMap(model *catalogmanagementv1.SolutionInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ArchitectureDiagrams != nil {
		architectureDiagrams := []map[string]interface{}{}
		for _, architectureDiagramsItem := range model.ArchitectureDiagrams {
			architectureDiagramsItemMap, err := resourceIBMCmVersionArchitectureDiagramToMap(&architectureDiagramsItem)
			if err != nil {
				return modelMap, err
			}
			architectureDiagrams = append(architectureDiagrams, architectureDiagramsItemMap)
		}
		modelMap["architecture_diagrams"] = architectureDiagrams
	}
	if model.Features != nil {
		features := []map[string]interface{}{}
		for _, featuresItem := range model.Features {
			featuresItemMap, err := resourceIBMCmVersionFeatureToMap(&featuresItem)
			if err != nil {
				return modelMap, err
			}
			features = append(features, featuresItemMap)
		}
		modelMap["features"] = features
	}
	if model.CostEstimate != nil {
		costEstimateMap, err := resourceIBMCmVersionCostEstimateToMap(model.CostEstimate)
		if err != nil {
			return modelMap, err
		}
		modelMap["cost_estimate"] = []map[string]interface{}{costEstimateMap}
	}
	if model.Dependencies != nil {
		dependencies := []map[string]interface{}{}
		for _, dependenciesItem := range model.Dependencies {
			dependenciesItemMap, err := resourceIBMCmVersionDependencyToMap(&dependenciesItem)
			if err != nil {
				return modelMap, err
			}
			dependencies = append(dependencies, dependenciesItemMap)
		}
		modelMap["dependencies"] = dependencies
	}
	return modelMap, nil
}

func resourceIBMCmVersionArchitectureDiagramToMap(model *catalogmanagementv1.ArchitectureDiagram) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Diagram != nil {
		diagramMap, err := resourceIBMCmVersionMediaItemToMap(model.Diagram)
		if err != nil {
			return modelMap, err
		}
		modelMap["diagram"] = []map[string]interface{}{diagramMap}
	}
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	return modelMap, nil
}

func resourceIBMCmVersionMediaItemToMap(model *catalogmanagementv1.MediaItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.URL != nil {
		modelMap["url"] = model.URL
	}
	if model.APIURL != nil {
		modelMap["api_url"] = model.APIURL
	}
	if model.URLProxy != nil {
		urlProxyMap, err := resourceIBMCmVersionURLProxyToMap(model.URLProxy)
		if err != nil {
			return modelMap, err
		}
		modelMap["url_proxy"] = []map[string]interface{}{urlProxyMap}
	}
	if model.Caption != nil {
		modelMap["caption"] = model.Caption
	}
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	if model.ThumbnailURL != nil {
		modelMap["thumbnail_url"] = model.ThumbnailURL
	}
	return modelMap, nil
}

func resourceIBMCmVersionURLProxyToMap(model *catalogmanagementv1.URLProxy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.URL != nil {
		modelMap["url"] = model.URL
	}
	if model.Sha != nil {
		modelMap["sha"] = model.Sha
	}
	return modelMap, nil
}

func resourceIBMCmVersionFeatureToMap(model *catalogmanagementv1.Feature) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Title != nil {
		modelMap["title"] = model.Title
	}
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	return modelMap, nil
}

func resourceIBMCmVersionCostEstimateToMap(model *catalogmanagementv1.CostEstimate) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Version != nil {
		modelMap["version"] = model.Version
	}
	if model.Currency != nil {
		modelMap["currency"] = model.Currency
	}
	if model.Projects != nil {
		projects := []map[string]interface{}{}
		for _, projectsItem := range model.Projects {
			projectsItemMap, err := resourceIBMCmVersionProjectToMap(&projectsItem)
			if err != nil {
				return modelMap, err
			}
			projects = append(projects, projectsItemMap)
		}
		modelMap["projects"] = projects
	}
	if model.Summary != nil {
		summaryMap, err := resourceIBMCmVersionCostSummaryToMap(model.Summary)
		if err != nil {
			return modelMap, err
		}
		modelMap["summary"] = []map[string]interface{}{summaryMap}
	}
	if model.TotalHourlyCost != nil {
		modelMap["total_hourly_cost"] = model.TotalHourlyCost
	}
	if model.TotalMonthlyCost != nil {
		modelMap["total_monthly_cost"] = model.TotalMonthlyCost
	}
	if model.PastTotalHourlyCost != nil {
		modelMap["past_total_hourly_cost"] = model.PastTotalHourlyCost
	}
	if model.PastTotalMonthlyCost != nil {
		modelMap["past_total_monthly_cost"] = model.PastTotalMonthlyCost
	}
	if model.DiffTotalHourlyCost != nil {
		modelMap["diff_total_hourly_cost"] = model.DiffTotalHourlyCost
	}
	if model.DiffTotalMonthlyCost != nil {
		modelMap["diff_total_monthly_cost"] = model.DiffTotalMonthlyCost
	}
	if model.TimeGenerated != nil {
		modelMap["time_generated"] = model.TimeGenerated.String()
	}
	return modelMap, nil
}

func resourceIBMCmVersionProjectToMap(model *catalogmanagementv1.Project) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.PastBreakdown != nil {
		pastBreakdownMap, err := resourceIBMCmVersionCostBreakdownToMap(model.PastBreakdown)
		if err != nil {
			return modelMap, err
		}
		modelMap["past_breakdown"] = []map[string]interface{}{pastBreakdownMap}
	}
	if model.Breakdown != nil {
		breakdownMap, err := resourceIBMCmVersionCostBreakdownToMap(model.Breakdown)
		if err != nil {
			return modelMap, err
		}
		modelMap["breakdown"] = []map[string]interface{}{breakdownMap}
	}
	if model.Diff != nil {
		diffMap, err := resourceIBMCmVersionCostBreakdownToMap(model.Diff)
		if err != nil {
			return modelMap, err
		}
		modelMap["diff"] = []map[string]interface{}{diffMap}
	}
	if model.Summary != nil {
		summaryMap, err := resourceIBMCmVersionCostSummaryToMap(model.Summary)
		if err != nil {
			return modelMap, err
		}
		modelMap["summary"] = []map[string]interface{}{summaryMap}
	}
	return modelMap, nil
}

func resourceIBMCmVersionCostBreakdownToMap(model *catalogmanagementv1.CostBreakdown) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.TotalHourlyCost != nil {
		modelMap["total_hourly_cost"] = model.TotalHourlyCost
	}
	if model.TotalMonthlyCOst != nil {
		modelMap["total_monthly_c_ost"] = model.TotalMonthlyCOst
	}
	if model.Resources != nil {
		resources := []map[string]interface{}{}
		for _, resourcesItem := range model.Resources {
			resourcesItemMap, err := resourceIBMCmVersionCostResourceToMap(&resourcesItem)
			if err != nil {
				return modelMap, err
			}
			resources = append(resources, resourcesItemMap)
		}
		modelMap["resources"] = resources
	}
	return modelMap, nil
}

func resourceIBMCmVersionCostResourceToMap(model *catalogmanagementv1.CostResource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.HourlyCost != nil {
		modelMap["hourly_cost"] = model.HourlyCost
	}
	if model.MonthlyCost != nil {
		modelMap["monthly_cost"] = model.MonthlyCost
	}
	if model.CostComponents != nil {
		costComponents := []map[string]interface{}{}
		for _, costComponentsItem := range model.CostComponents {
			costComponentsItemMap, err := resourceIBMCmVersionCostComponentToMap(&costComponentsItem)
			if err != nil {
				return modelMap, err
			}
			costComponents = append(costComponents, costComponentsItemMap)
		}
		modelMap["cost_components"] = costComponents
	}
	return modelMap, nil
}

func resourceIBMCmVersionCostComponentToMap(model *catalogmanagementv1.CostComponent) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Unit != nil {
		modelMap["unit"] = model.Unit
	}
	if model.HourlyQuantity != nil {
		modelMap["hourly_quantity"] = model.HourlyQuantity
	}
	if model.MonthlyQuantity != nil {
		modelMap["monthly_quantity"] = model.MonthlyQuantity
	}
	if model.Price != nil {
		modelMap["price"] = model.Price
	}
	if model.HourlyCost != nil {
		modelMap["hourly_cost"] = model.HourlyCost
	}
	if model.MonthlyCost != nil {
		modelMap["monthly_cost"] = model.MonthlyCost
	}
	return modelMap, nil
}

func resourceIBMCmVersionCostSummaryToMap(model *catalogmanagementv1.CostSummary) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.TotalDetectedResources != nil {
		modelMap["total_detected_resources"] = flex.IntValue(model.TotalDetectedResources)
	}
	if model.TotalSupportedResources != nil {
		modelMap["total_supported_resources"] = flex.IntValue(model.TotalSupportedResources)
	}
	if model.TotalUnsupportedResources != nil {
		modelMap["total_unsupported_resources"] = flex.IntValue(model.TotalUnsupportedResources)
	}
	if model.TotalUsageBasedResources != nil {
		modelMap["total_usage_based_resources"] = flex.IntValue(model.TotalUsageBasedResources)
	}
	if model.TotalNoPriceResources != nil {
		modelMap["total_no_price_resources"] = flex.IntValue(model.TotalNoPriceResources)
	}
	return modelMap, nil
}

func resourceIBMCmVersionDependencyToMap(model *catalogmanagementv1.Dependency) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CatalogID != nil {
		modelMap["catalog_id"] = model.CatalogID
	}
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Version != nil {
		modelMap["version"] = model.Version
	}
	if model.Flavors != nil {
		modelMap["flavors"] = model.Flavors
	}
	return modelMap, nil
}

func getVersionFromOffering(oldOffering catalogmanagementv1.Offering, newOffering catalogmanagementv1.Offering) (string, error) {
	var oldVersionList []string
	var newVersionList []string

	for _, kind := range oldOffering.Kinds {
		// oldVersionList = append(oldVersionList, kind.Versions...)
		for _, version := range kind.Versions {
			oldVersionList = append(oldVersionList, *version.ID)
		}
	}

	for _, kind := range newOffering.Kinds {
		// newVersionList = append(newVersionList, kind.Versions...)
		for _, version := range kind.Versions {
			newVersionList = append(newVersionList, *version.ID)
		}
	}

	for _, newVer := range newVersionList {
		isOld := false
		for _, oldVer := range oldVersionList {
			if newVer == oldVer {
				isOld = true
				break
			}
		}
		if !isOld {
			return newVer, nil
		}
	}
	return "", fmt.Errorf("error finding imported version")
}
