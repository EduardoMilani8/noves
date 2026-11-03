package main

import (
	"context"
	"time"

	"github.com/spf13/cobra"

	"github.com/EduardoMilani8/noves/internal/adapters/firebird"
	"github.com/EduardoMilani8/noves/internal/adapters/postgres"
	"github.com/EduardoMilani8/noves/internal/config"
	"github.com/EduardoMilani8/noves/internal/doctor"
	"github.com/EduardoMilani8/noves/internal/ports"
)

func newDoctorCmd(opts *globalOptions) *cobra.Command {
	var timeout time.Duration
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Confere conexões, versões, charset e permissão de leitura",
		Long: `Conecta ao Firebird e ao PostgreSQL configurados e mostra versão, ODS,
dialeto e charset do Firebird, versão do PostgreSQL e se o usuário consegue
ler todas as tabelas. Nunca escreve nos bancos.

Saída: 0 tudo certo; 1 falta permissão de leitura; 2 erro de conexão ou nada
configurado.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, cancel := context.WithTimeout(cmd.Context(), timeout)
			defer cancel()

			fb, pg := opts.dsns()
			results := doctor.Run(ctx, []doctor.Target{
				{Name: "Firebird", DSN: fb, EnvVar: config.EnvFirebirdDSN, Flag: "--fb", Open: openFirebird},
				{Name: "PostgreSQL", DSN: pg, EnvVar: config.EnvPostgresDSN, Flag: "--pg", Open: openPostgres},
			})
			doctor.Render(cmd.OutOrStdout(), opts.msgs, results)
			if code := doctor.ExitCode(results); code != exitOK {
				return &exitCodeError{code: code}
			}
			return nil
		},
	}
	cmd.Flags().DurationVar(&timeout, "timeout", 30*time.Second, "tempo máximo para conectar e consultar")
	return cmd
}

func openFirebird(ctx context.Context, dsn string) (ports.Inspector, error) {
	return firebird.Open(ctx, dsn)
}

func openPostgres(ctx context.Context, dsn string) (ports.Inspector, error) {
	return postgres.Open(ctx, dsn)
}
