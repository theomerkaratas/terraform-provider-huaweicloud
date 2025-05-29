package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccRdsBackupStop_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstantJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsBackupStop_basic(),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func testDataSourceRdsBackupStop_base() string {
	return fmt.Sprintf(`
%s
resource "huaweicloud_networking_secgroup" "secgroup_test" {
  name = "secgroup_test"
}
resource "huaweicloud_rds_instance" "test" {
  name              = "test"
  flavor            = "rds.mysql.x1.large.2"
  vpc_id            = data.huaweicloud_vpc.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.secgroup_test.id
  charging_mode     = "postPaid"
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]
  db {
    type     = "MySQL"
    version  = "8.0"
    password = "Terraform145!"
  }
  volume {
    type = "CLOUDSSD"
    size = 40
  }
}

resource "huaweicloud_rds_backup" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  name        = "test"
}
`, testAccRdsInstance_base())
}

func testAccRdsBackupStop_basic() string {
	return fmt.Sprintf(`
%[1]s
resource "huaweicloud_rds_backup_stop" "test" {
  instance_id = "%[2]s"
}
`, testDataSourceRdsBackupStop_base(), acceptance.HW_RDS_INSTANCE_ID)
}
