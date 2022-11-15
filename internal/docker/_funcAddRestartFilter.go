package docker

// AddRestartFilter
//
// English:
//
//
// Português: não tem lógica
//
//
func (e *ContainerBuilder) AddRestartFilter(label, match, filter, search, replace string) {
	if e.chaos.filterRestart == nil {
		e.chaos.filterRestart = make([]LogFilter, 0)
	}

	e.chaos.filterRestart = append(
		e.chaos.filterRestart,
		LogFilter{
			Label:   label,
			Match:   match,
			Filter:  filter,
			Search:  search,
			Replace: replace,
		},
	)
}
