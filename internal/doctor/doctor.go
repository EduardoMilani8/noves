// Package doctor confere se o noves consegue conversar com os bancos: conexão,
// versões, charset e permissão de leitura.
package doctor

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/EduardoMilani8/noves/internal/core"
	"github.com/EduardoMilani8/noves/internal/dsn"
	"github.com/EduardoMilani8/noves/internal/i18n"
	"github.com/EduardoMilani8/noves/internal/ports"
)

// Opener abre a conexão com um banco.
type Opener func(ctx context.Context, dsn string) (ports.Inspector, error)

// Target é um banco a conferir.
type Target struct {
	Name   string // "Firebird", "PostgreSQL"
	DSN    string // vazio quando não configurado
	EnvVar string // variável de ambiente que define o DSN
	Flag   string // flag equivalente
	Open   Opener
}

// Result é o que o doctor descobriu sobre um banco.
type Result struct {
	Target  Target
	Skipped bool
	Info    core.DatabaseInfo
	Read    core.ReadCheck
	ConnErr error // falha ao conectar ou ler o catálogo
	ReadErr error // falha ao conferir a leitura
}

// Issues conta os problemas do resultado.
func (r Result) Issues() int {
	switch {
	case r.Skipped:
		return 0
	case r.ConnErr != nil:
		return 1
	case r.ReadErr != nil || !r.Read.OK():
		return 1
	}
	return 0
}

// Run confere cada banco configurado. Nunca escreve nada nos bancos.
func Run(ctx context.Context, targets []Target) []Result {
	results := make([]Result, 0, len(targets))
	for _, t := range targets {
		results = append(results, check(ctx, t))
	}
	return results
}

func check(ctx context.Context, t Target) Result {
	r := Result{Target: t}
	if t.DSN == "" {
		r.Skipped = true
		return r
	}
	insp, err := t.Open(ctx, t.DSN)
	if err != nil {
		r.ConnErr = err
		return r
	}
	defer insp.Close()

	if r.Info, err = insp.Info(ctx); err != nil {
		r.ConnErr = err
		return r
	}
	r.Read, r.ReadErr = insp.CheckRead(ctx)
	return r
}

// ExitCode segue a convenção do noves: 0 ok, 1 achados, 2 erro de uso ou de
// conexão.
func ExitCode(results []Result) int {
	code := 0
	configured := 0
	for _, r := range results {
		if r.Skipped {
			continue
		}
		configured++
		if r.ConnErr != nil {
			return 2
		}
		if r.Issues() > 0 {
			code = 1
		}
	}
	if configured == 0 {
		return 2
	}
	return code
}

// Render escreve o relatório em texto.
func Render(w io.Writer, m i18n.Messages, results []Result) {
	configured, issues := 0, 0
	for _, r := range results {
		if !r.Skipped {
			configured++
		}
		issues += r.Issues()
	}
	if configured == 0 {
		fmt.Fprintln(w, m.T("doctor.nothing"))
		return
	}

	for _, r := range results {
		fmt.Fprintf(w, "%s", r.Target.Name)
		if r.Skipped {
			fmt.Fprintf(w, ": %s\n\n", m.T("doctor.skipped", r.Target.EnvVar, r.Target.Flag))
			continue
		}
		fmt.Fprintf(w, "  %s\n", dsn.Mask(r.Target.DSN))
		if r.ConnErr != nil {
			fmt.Fprintf(w, "  ✗ %s: %s\n\n", m.T("doctor.failed"), dsn.Scrub(r.ConnErr.Error(), r.Target.DSN))
			continue
		}
		fmt.Fprintf(w, "  ✓ %s\n", m.T("doctor.connected"))
		for _, f := range fields(m, r.Info) {
			fmt.Fprintf(w, "    %-10s %s\n", f[0], f[1])
		}
		switch {
		case r.ReadErr != nil:
			fmt.Fprintf(w, "  ✗ %s: %s\n", m.T("doctor.read_error"), dsn.Scrub(r.ReadErr.Error(), r.Target.DSN))
		case r.Read.OK():
			fmt.Fprintf(w, "  ✓ %s\n", m.T("doctor.read_ok", r.Read.Tables))
		default:
			fmt.Fprintf(w, "  ✗ %s\n", m.T("doctor.read_denied",
				len(r.Read.Unreadable), r.Read.Tables, strings.Join(r.Read.Unreadable, ", ")))
		}
		fmt.Fprintln(w)
	}

	if issues == 0 {
		fmt.Fprintln(w, m.T("doctor.summary_ok"))
	} else {
		fmt.Fprintln(w, m.T("doctor.summary_issues", issues))
	}
}

func fields(m i18n.Messages, in core.DatabaseInfo) [][2]string {
	var out [][2]string
	add := func(key, v string) {
		if v != "" && v != "0" {
			out = append(out, [2]string{m.T(key), v})
		}
	}
	add("doctor.version", in.Version)
	add("doctor.ods", in.ODS)
	add("doctor.dialect", strconv.Itoa(in.Dialect))
	add("doctor.charset", in.Charset)
	add("doctor.page_size", strconv.Itoa(in.PageSize))
	add("doctor.database", in.Database)
	add("doctor.user", in.User)
	return out
}
