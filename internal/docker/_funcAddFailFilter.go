package docker

func (e *ContainerBuilder) AddFailFilter(label, match, filter, search, replace string) {
	if e.chaos.filterFail == nil {
		e.chaos.filterFail = make([]LogFilter, 0)
	}

	e.chaos.filterFail = append(
		e.chaos.filterFail,
		LogFilter{
			Label:   label,
			Match:   match,
			Filter:  filter,
			Search:  search,
			Replace: replace,
		},
	)
}
