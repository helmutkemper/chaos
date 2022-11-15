package docker

type LogFilter struct {
	Label string

	// Texto contido na linha (tudo ou nada)
	Match string

	// expressão regular contendo o filtro para capturar o elemento
	// Ex.: ^(.*?)(?P<valueToGet>\\d+)(.*)
	Filter string

	// texto usado em replaceAll
	// Ex.: search: "." replace: "," para compatibilizar número com o excel
	Search  string
	Replace string

	// path to save container default output into file format
	LogPath string

	// tamanho do log quando o último evento ocorreu, para que o mesmo evento não seja capturado indefinidamente
	size int
}
