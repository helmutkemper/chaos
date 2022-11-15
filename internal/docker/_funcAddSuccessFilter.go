package docker

func (e *ContainerBuilder) AddSuccessFilter(label, match, filter, search, replace string) {
	if e.chaos.filterSuccess == nil {
		e.chaos.filterSuccess = make([]LogFilter, 0)
	}

	e.chaos.filterSuccess = append(
		e.chaos.filterSuccess,
		LogFilter{
			Label:   label,
			Match:   match,
			Filter:  filter,
			Search:  search,
			Replace: replace,
		},
	)
}
