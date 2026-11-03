package ports

import (
	"context"

	"github.com/EduardoMilani8/noves/internal/core"
)

// Inspector dá ao núcleo informações básicas de um banco conectado. Cada
// banco de origem ou destino tem um adaptador que implementa esta porta.
type Inspector interface {
	// Info devolve versão, charset e demais dados do servidor e do banco.
	Info(ctx context.Context) (core.DatabaseInfo, error)
	// CheckRead confere se o usuário consegue ler todas as tabelas de usuário.
	CheckRead(ctx context.Context) (core.ReadCheck, error)
	// Close encerra a conexão.
	Close() error
}
