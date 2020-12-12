package ibm

import (
	"fmt"
	"strings"
	"testing"

	"github.com/IBM/go-sdk-core/v3/core"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccIBMCisDNSRecordsImport_Basic(t *testing.T) {
	name := "ibm_cis_dns_records_import." + "test"
	file := "test-fixtures/dns_records_import.txt"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCis(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCisDNSRecordsImportConfigBasic1(file),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "file", file),
					testAccCheckIBMCisDNSRecordsImportRemoveImportedRecords(name),
				),
			},
		},
	})
}

func testAccCheckCisDNSRecordsImportConfigBasic1(file string) string {
	return testAccCheckIBMCisDNSRecordConfigCisDSBasic(
		"test-dns-record", cisDomainStatic) +
		fmt.Sprintf(`
		resource "ibm_cis_dns_records_import" "test" {
			cis_id    = data.ibm_cis.cis.id
			domain_id = data.ibm_cis_domain.cis_domain.domain_id
			file      = "%[1]s"
		}`, file)
}

func testAccCheckIBMCisDNSRecordsImportRemoveImportedRecords(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		cisClient, err := testAccProvider.Meta().(ClientSession).CisDNSRecordClientSession()
		if err != nil {
			return err
		}
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		idSplitStr := strings.SplitN(rs.Primary.ID, ":", 5)
		zoneID := idSplitStr[3]
		crn := idSplitStr[4]
		cisClient.Crn = core.StringPtr(crn)
		cisClient.ZoneIdentifier = core.StringPtr(zoneID)
		listOpt := cisClient.NewListAllDnsRecordsOptions()
		result, _, err := cisClient.ListAllDnsRecords(listOpt)
		if err != nil {
			return err
		}
		for _, record := range result.Result {
			if strings.Contains(*record.Name, "test-import") {
				delOpt := cisClient.NewDeleteDnsRecordOptions(*record.ID)
				_, _, err = cisClient.DeleteDnsRecord(delOpt)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
}
