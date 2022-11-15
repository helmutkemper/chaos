package docker

// SetImageBuildOptionsSecurityOpt
//
// English:
//
//	Set the container security options
//
//	 Input:
//	   values: container security options
//
// Examples:
//
//	label=user:USER        — Set the label user for the container
//	label=role:ROLE        — Set the label role for the container
//	label=type:TYPE        — Set the label type for the container
//	label=level:LEVEL      — Set the label level for the container
//	label=disable          — Turn off label confinement for the container
//	apparmor=PROFILE       — Set the apparmor profile to be applied to the container
//	no-new-privileges:true — Disable container processes from gaining new privileges
//	seccomp=unconfined     — Turn off seccomp confinement for the container
//	seccomp=profile.json   — White-listed syscalls seccomp Json file to be used as a seccomp filter
//
// Português:
//
//	Modifica as opções de segurança do container
//
//	 Entrada:
//	   values: opções de segurança do container
//
// Exemplos:
//
//	label=user:USER        — Determina o rótulo user para o container
//	label=role:ROLE        — Determina o rótulo role para o container
//	label=type:TYPE        — Determina o rótulo type para o container
//	label=level:LEVEL      — Determina o rótulo level para o container
//	label=disable          — Desliga o confinamento do rótulo para o container
//	apparmor=PROFILE       — Habilita o perfil definido pelo apparmor do linux para ser definido ao container
//	no-new-privileges:true — Impede o processo do container a ganhar novos privilégios
//	seccomp=unconfined     — Desliga o confinamento causado pelo seccomp do linux ao container
//	seccomp=profile.json   — White-listed syscalls seccomp Json file to be used as a seccomp filter
func (e *ContainerBuilder) SetImageBuildOptionsSecurityOpt(value []string) {
	e.buildOptions.SecurityOpt = value
}
