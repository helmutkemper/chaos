package docker

// BlkioStatEntry
//
// Português:
//
//	Estrutura para armazenar uma peça de estatísticas de Blkio
//
//	Não usado no windows.
//
// English:
//
//	Structure to store a piece of Blkio stats
//
//	Not used on Windows.
type BlkioStatEntry struct {
	Major uint64 `json:"major"`
	Minor uint64 `json:"minor"`
	Op    string `json:"op"`
	Value uint64 `json:"value"`
}
