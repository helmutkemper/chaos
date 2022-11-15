package docker

// BlkioStats
//
// English:
//
//	Stores All IO service stats for data read and write.
//
//	This is a Linux specific structure as the differences between expressing block I/O on Windows and Linux are sufficiently significant to make little sense attempting to morph into a combined structure.
//
// Português:
//
//	Armazena todos os estatísticas de serviço de IO para leitura e escrita de dados.
//
//	Este é um estrutura Linux específica devido às diferenças entre expressar o IO de bloco no Windows e Linux são suficientemente significativas para fazer pouco sentido tentar morfar em uma combinação de estrutura.
type BlkioStats struct {
	// number of bytes transferred to and from the block device
	IoServiceBytesRecursive []BlkioStatEntry `json:"io_service_bytes_recursive"`
	IoServicedRecursive     []BlkioStatEntry `json:"io_serviced_recursive"`
	IoQueuedRecursive       []BlkioStatEntry `json:"io_queue_recursive"`
	IoServiceTimeRecursive  []BlkioStatEntry `json:"io_service_time_recursive"`
	IoWaitTimeRecursive     []BlkioStatEntry `json:"io_wait_time_recursive"`
	IoMergedRecursive       []BlkioStatEntry `json:"io_merged_recursive"`
	IoTimeRecursive         []BlkioStatEntry `json:"io_time_recursive"`
	SectorsRecursive        []BlkioStatEntry `json:"sectors_recursive"`
}
