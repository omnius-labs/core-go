package migration_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/omnius-labs/core-go/migration"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestMigration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "migration Spec")
}

var _ = Describe("Migration Test", func() {
	var mysqlC testcontainers.Container
	var url string

	BeforeEach(func() {
		ctx := context.Background()
		req := testcontainers.ContainerRequest{
			Image:        "mysql:8.2",
			ExposedPorts: []string{"3306/tcp"},
			Env: map[string]string{
				"MYSQL_ROOT_PASSWORD": "password",
			},
			WaitingFor: wait.ForAll(wait.ForListeningPort("3306/tcp"), wait.ForLog("mysqld: ready for connections").WithOccurrence(2)),
		}
		var err error
		mysqlC, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
		Expect(err).NotTo(HaveOccurred())

		ip, err := mysqlC.Host(ctx)
		Expect(err).NotTo(HaveOccurred())

		port, err := mysqlC.MappedPort(ctx, "3306")
		Expect(err).NotTo(HaveOccurred())

		url = fmt.Sprintf("root:password@tcp(%s:%s)/mysql?parseTime=true", ip, port.Port())
	})

	AfterEach(func() {
		ctx := context.Background()
		mysqlC.Terminate(ctx)
	})

	It("simple create table test", Serial, func() {
		path := "./case/simple_create_table"
		username := "your_username"
		description := "your_description"
		migrator, err := migration.NewMySQLMigrator(url, path, username, description)
		Expect(err).NotTo(HaveOccurred())

		err = migrator.Migrate()
		Expect(err).NotTo(HaveOccurred())
	})

	It("create table syntax error test", Serial, func() {
		path := "./case/create_table_syntax_error_test"
		username := "your_username"
		description := "your_description"
		migrator, err := migration.NewMySQLMigrator(url, path, username, description)
		Expect(err).NotTo(HaveOccurred())

		err = migrator.Migrate()
		Expect(err).Error()
	})

	It("migrate twice test", Serial, func() {
		{
			path := "./case/simple_create_table"
			username := "your_username"
			description := "your_description"
			migrator, err := migration.NewMySQLMigrator(url, path, username, description)
			Expect(err).NotTo(HaveOccurred())

			err = migrator.Migrate()
			Expect(err).NotTo(HaveOccurred())
		}

		{
			path := "./case/simple_create_table"
			username := "your_username"
			description := "your_description"
			migrator, err := migration.NewMySQLMigrator(url, path, username, description)
			Expect(err).NotTo(HaveOccurred())

			err = migrator.Migrate()
			Expect(err).NotTo(HaveOccurred())
		}
	})
})
