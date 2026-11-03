package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/EduardoMilani8/noves/internal/config"
	"github.com/EduardoMilani8/noves/internal/i18n"
)

// Códigos de saída do noves.
const (
	exitOK       = 0
	exitFindings = 1
	exitError    = 2
)

// exitCodeError carrega um código de saída diferente de 2 até o main.
type exitCodeError struct {
	code int
	err  error
}

func (e *exitCodeError) Error() string {
	if e.err == nil {
		return fmt.Sprintf("código de saída %d", e.code)
	}
	return e.err.Error()
}

type globalOptions struct {
	lang  string
	fbDSN string
	pgDSN string
	msgs  i18n.Messages
}

func execute(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	// .env é conveniência local; variáveis já exportadas têm prioridade.
	if err := config.LoadDotEnv(".env"); err != nil {
		fmt.Fprintln(stderr, "noves:", err)
		return exitError
	}

	root := newRootCmd()
	root.SetArgs(args)
	root.SetOut(stdout)
	root.SetErr(stderr)

	err := root.ExecuteContext(ctx)
	if err == nil {
		return exitOK
	}
	var ec *exitCodeError
	if errors.As(err, &ec) {
		if ec.err != nil {
			fmt.Fprintln(stderr, "noves:", ec.err)
		}
		return ec.code
	}
	fmt.Fprintln(stderr, "noves:", err)
	return exitError
}

func newRootCmd() *cobra.Command {
	opts := &globalOptions{}
	root := &cobra.Command{
		Use:   "noves",
		Short: "A prova dos nove de migrações Firebird → PostgreSQL",
		Long: `noves avalia o risco antes de migrar de Firebird para PostgreSQL e prova,
depois, que os dados e o comportamento são os mesmos.

Converta com o que quiser. Prove com o noves.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(*cobra.Command, []string) error {
			lang, err := i18n.Parse(opts.lang)
			if err != nil {
				return err
			}
			opts.msgs = i18n.New(lang)
			return nil
		},
	}

	f := root.PersistentFlags()
	f.StringVar(&opts.lang, "lang", os.Getenv(config.EnvLang), "idioma das mensagens: pt ou en (env NOVES_LANG)")
	f.StringVar(&opts.fbDSN, "fb", "", "DSN do Firebird, ex.: usuario:senha@host:3050/caminho/banco.fdb (env NOVES_FB_DSN)")
	f.StringVar(&opts.pgDSN, "pg", "", "DSN do PostgreSQL, ex.: postgres://usuario:senha@host/banco (env NOVES_PG_DSN)")

	root.AddCommand(newVersionCmd(), newDoctorCmd(opts))
	return root
}

// dsns resolve os DSNs: flag primeiro, depois variável de ambiente.
func (o *globalOptions) dsns() (fb, pg string) {
	fb, pg = o.fbDSN, o.pgDSN
	if fb == "" {
		fb = os.Getenv(config.EnvFirebirdDSN)
	}
	if pg == "" {
		pg = os.Getenv(config.EnvPostgresDSN)
	}
	return fb, pg
}
