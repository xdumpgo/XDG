package main

import (
	"fmt"
	"github.com/xdumpgo/XDG/injection"
	"github.com/xdumpgo/XDG/injection/mssql"
	"github.com/xdumpgo/XDG/injection/mysql"
	"github.com/xdumpgo/XDG/injection/oracle"
	"github.com/xdumpgo/XDG/injection/postgres"
	"strings"
)

func main() {
	mods := len(injection.Modulators)

	my_err := len(mysql.E_Vectors) *2
	ms_err := len(mssql.E_Vectors) *2
	po_err := len(postgres.E_Vectors) *2
	or_err := len(oracle.E_Vectors) *2

	gen_bli := len(injection.B_Vectors)
	my_bli := len(mysql.B_Vectors)
	po_bli := len(postgres.B_Vectors)
	or_bli := len(oracle.B_Vectors)
	ms_bli := len(mssql.B_Vectors)

	union := len(injection.U_Vectors)

	total := 0

	total += my_err * mods
	total += ms_err * mods
	total += po_err * mods
	total += or_err * mods
	totalb := gen_bli * 4
	totalb += my_bli + po_bli + or_bli + ms_bli

	var s []string
	for name,val := range mysql.E_Vectors {
		s = append(s, fmt.Sprintf("(\"mysql_v_%s\", \"%s\", 57)", name, val))
	}
	for name,val := range mssql.E_Vectors {
		s = append(s, fmt.Sprintf("(\"mssql_v_%s\", \"%s\", 57)", name, val))
	}
	for name,val := range postgres.E_Vectors {
		s = append(s, fmt.Sprintf("(\"postgres_v_%s\", \"%s\", 57)", name, val))
	}
	for name,val := range oracle.E_Vectors {
		s = append(s, fmt.Sprintf("(\"oracle_v_%s\", \"%s\", 57)", name, val))
	}
	for name,val := range injection.U_Vectors {
		s = append(s, fmt.Sprintf("(\"union_v_%s\", \"%s\", 57)", name, val))
	}

	for name,val := range mysql.Queries {
		s = append(s, fmt.Sprintf("(\"mysql_%s\", \"%s\", 57)", name, val))
	}
	for name,val := range mssql.Queries {
		s = append(s, fmt.Sprintf("(\"mssql_%s\", \"%s\", 57)", name, val))
	}
	for name,val := range oracle.Queries {
		s = append(s, fmt.Sprintf("(\"oracle_%s\", \"%s\", 57)", name, val))
	}
	for name,val := range postgres.Queries {
		s = append(s, fmt.Sprintf("(\"postgres_%s\", \"%s\", 57)", name, val))
	}

	insert := fmt.Sprintf("INSERT OR IGNORE INTO variables (name, value, program_id) VALUES %s", strings.Join(s, ","))

	fmt.Println("XDumpGO Injection statistics")
	fmt.Printf("Active Modulators: [M:%d]\n", mods)
	fmt.Printf("MySQL:      [E:%d]\t[U:%d]\t[B:%d]\n", my_err*mods, union*mods, (gen_bli+my_bli)*mods)
	fmt.Printf("MSSQL:      [E:%d]\t[U:%d]\t[B:%d]\n", ms_err*mods, union*mods, (gen_bli+ms_bli)*mods)
	fmt.Printf("PostGreSQL: [E:%d]\t[U:%d]\t[B:%d]\n", po_err*mods, union*mods, (gen_bli+po_bli)*mods)
	fmt.Printf("Oracle:     [E:%d]\t[U:%d]\t[B:%d]\n", or_err*mods, union*mods, (gen_bli+or_bli)*mods)
	fmt.Printf("Total:      [E:%d]\t[U:%d]\t[B:%d]\n", total, (union * 4) * mods, totalb * mods)
	fmt.Printf("SQL Insert\n%s", insert)
}
