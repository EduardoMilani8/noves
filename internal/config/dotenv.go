// Package config lê a configuração do noves de flags, variáveis de ambiente e
// do arquivo .env do diretório atual.
package config

import (
	"bufio"
	"errors"
	"io"
	"io/fs"
	"os"
	"strings"
)

// Variáveis de ambiente reconhecidas.
const (
	EnvFirebirdDSN = "NOVES_FB_DSN"
	EnvPostgresDSN = "NOVES_PG_DSN"
	EnvLang        = "NOVES_LANG"
)

// LoadDotEnv carrega um arquivo .env sem sobrescrever variáveis já definidas
// no ambiente. Arquivo inexistente não é erro.
func LoadDotEnv(path string) error {
	f, err := os.Open(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	defer f.Close()

	vars, err := ParseDotEnv(f)
	if err != nil {
		return err
	}
	for k, v := range vars {
		if _, ok := os.LookupEnv(k); !ok {
			if err := os.Setenv(k, v); err != nil {
				return err
			}
		}
	}
	return nil
}

// ParseDotEnv lê linhas CHAVE=valor, ignorando comentários e linhas vazias.
// Aceita o prefixo "export " e valores entre aspas simples ou duplas.
func ParseDotEnv(r io.Reader) (map[string]string, error) {
	vars := map[string]string{}
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		k, v = strings.TrimSpace(k), strings.TrimSpace(v)
		if len(v) >= 2 && (v[0] == '"' || v[0] == '\'') && v[len(v)-1] == v[0] {
			v = v[1 : len(v)-1]
		}
		vars[k] = v
	}
	return vars, sc.Err()
}
