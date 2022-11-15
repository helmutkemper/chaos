package docker

// SetImageBuildOptionsCacheFrom
//
// English:
//
//	Specifies images that are used for matching cache.
//
//	 Entrada:
//	   values: images that are used for matching cache.
//
// Note:
//
//	Images specified here do not need to have a valid parent chain to match cache.
//
// Português:
//
//	Especifica imagens que são usadas para correspondência de cache.
//
//	 Entrada:
//	   values: imagens que são usadas para correspondência de cache.
//
// Note:
//
//	As imagens especificadas aqui não precisam ter uma cadeia pai válida para corresponder a cache.
func (e *ContainerBuilder) SetImageBuildOptionsCacheFrom(values []string) {
	e.buildOptions.CacheFrom = values
}
