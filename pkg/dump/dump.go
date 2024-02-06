package dump

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/mylittleboxy/backup/pkg/configType"
)

func Dump(conf configType.Config) (string, error) {
	// Open connection to database
	config := mysql.NewConfig()
	config.User = conf.DB.User
	config.Passwd = conf.DB.Passwd
	config.DBName = "bitnami_prestashop"
	config.Net = "tcp"
	config.Addr = conf.DB.Addr
	config.ParseTime = true

	dumpFilenameFormat := fmt.Sprintf("%s/%s-%d.sql", conf.DB.DumpDir, "bitnami_prestashop", time.Now().Unix()) // accepts time layout string and add .sql at the end of file

	// Build the mysqldump command
	cmd := exec.Command("mysqldump",
		"--host="+conf.DB.Addr,
		"--user="+conf.DB.User,
		"--password="+conf.DB.Passwd,
		"--single-transaction",
		"--create-options",
		"--extended-insert",
		"--hex-blob",
		"--complete-insert",
		"bitnami_prestashop",
	)

	// Redirect the output to a file (dump.sql)
	outfile, err := exec.Command("bash", "-c", fmt.Sprintf("%s > %s", cmd.String(), dumpFilenameFormat)).CombinedOutput()
	if err != nil {
		fmt.Println("Error executing mysqldump:", err)
		fmt.Println("Output:", string(outfile))
		return "", err
	}

	fmt.Println("Dump completed successfully. Output saved to", dumpFilenameFormat)
	return dumpFilenameFormat, nil

	/*db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		fmt.Println("Error opening database: ", err)
		return "", err
	}

	// Register database with mysqldump
	dumper, err := mysqldump.Register(db, conf.DB.DumpDir, dumpFilenameFormat)
	if err != nil {
		fmt.Println("Error registering databse:", err)
		return "", err
	}
	// Dump database to file
	err = dumper.Dump()
	if err != nil {
		fmt.Println("Error dumping:", err)
		return "", err
	}
	file, _ := dumper.Out.(*os.File)

	// Close dumper, connected database and file stream.
	dumper.Close()
	return file.Name(), nil*/
}
