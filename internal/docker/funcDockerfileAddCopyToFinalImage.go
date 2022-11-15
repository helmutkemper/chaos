package docker

// DockerfileAddCopyToFinalImage
//
// English:
//
//	Add one instruction 'COPY --from=builder /app/`dst` `src`' to final image builder.
//
// Português:
//
//	Adiciona uma instrução 'COPY --from=builder /app/`dst` `src`' ao builder da imagem final.
func (e *ContainerBuilder) DockerfileAddCopyToFinalImage(src, dst string) {
	if e.copyFile == nil {
		e.copyFile = make([]CopyFile, 0)
	}

	e.copyFile = append(
		e.copyFile,
		CopyFile{
			Src: src,
			Dst: dst,
		},
	)
}
