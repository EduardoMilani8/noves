// Command seed popula o PostgreSQL de teste a partir do ERP sintético do
// Firebird e injeta as divergências conhecidas. Usado por make up e pelo CI.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/EduardoMilani8/noves/internal/testkit"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "seed:", err)
		os.Exit(2)
	}
}

func run() error {
	divergencias := flag.String("divergencias", testkit.ArquivoDivergente, "script SQL com as divergências (vazio para não aplicar)")
	flag.Parse()

	var sqlText string
	if *divergencias != "" {
		b, err := os.ReadFile(*divergencias)
		if err != nil {
			return err
		}
		sqlText = string(b)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	fb, pg := testkit.DSNs()
	logf := func(format string, args ...any) { fmt.Fprintf(os.Stderr, format+"\n", args...) }
	return testkit.Seed(ctx, fb, pg, sqlText, logf)
}
