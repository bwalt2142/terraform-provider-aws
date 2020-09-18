package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAWSRDSEngineVersionDataSource_basic(t *testing.T) {
	dataSourceName := "data.aws_rds_engine_version.test"
	engine := "mysql"
	version := "5.7.17"
	paramGroup := "mysql5.7"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSRDSEngineVersion(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSRDSEngineVersionDataSourceBasicConfig(engine, version, paramGroup),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "engine", engine),
					resource.TestCheckResourceAttr(dataSourceName, "version", version),
					resource.TestCheckResourceAttr(dataSourceName, "parameter_group_name", paramGroup),
				),
			},
		},
	})
}

func TestAccAWSRDSEngineVersionDataSource_preferred(t *testing.T) {
	dataSourceName := "data.aws_rds_engine_version.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSRDSEngineVersion(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSRDSEngineVersionDataSourcePreferredConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "version", "5.7.19"),
				),
			},
		},
	})
}

func TestAccAWSRDSEngineVersionDataSource_defaultOnly(t *testing.T) {
	dataSourceName := "data.aws_rds_engine_version.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSRDSEngineVersion(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSRDSEngineVersionDataSourceDefaultOnlyConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "version"),
				),
			},
		},
	})
}

func testAccPreCheckAWSRDSEngineVersion(t *testing.T) {
	conn := testAccProvider.Meta().(*AWSClient).rdsconn

	input := &rds.DescribeDBEngineVersionsInput{
		Engine:      aws.String("mysql"),
		DefaultOnly: aws.Bool(true),
	}

	_, err := conn.DescribeDBEngineVersions(input)

	if testAccPreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}

	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}

func testAccAWSRDSEngineVersionDataSourceBasicConfig(engine, version, paramGroup string) string {
	return fmt.Sprintf(`
data "aws_rds_engine_version" "test" {
  engine               = %q
  version              = %q
  parameter_group_name = %q
}
`, engine, version, paramGroup)
}

func testAccAWSRDSEngineVersionDataSourcePreferredConfig() string {
	return fmt.Sprintf(`
data "aws_rds_engine_version" "test" {
  engine             = "mysql"
  preferred_versions = ["85.9.12", "5.7.19", "5.7.17"]
}
`)
}

func testAccAWSRDSEngineVersionDataSourceDefaultOnlyConfig() string {
	return fmt.Sprintf(`
data "aws_rds_engine_version" "test" {
  engine = "mysql"
}
`)
}
