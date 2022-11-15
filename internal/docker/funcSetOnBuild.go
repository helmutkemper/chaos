package docker

// SetOnBuild
//
// English:
//
//	Adds to the image a trigger instruction to be executed at a later time, when the image is used as
//	the base for another build.
//
//	 Input:
//	   onBuild: List of trigger instruction to be executed at a later time, when the image is used as
//	    the base for another build
//
//	The trigger will be executed in the context of the downstream build, as if it had been
//	inserted immediately after the FROM instruction in the downstream Dockerfile.
//
//	Any build instruction can be registered as a trigger.
//
//	This is useful if you are building an image which will be used as a base to build other
//	images, for example an application build environment or a daemon which may be
//	customized with user-specific configuration.
//
//	For example, if your image is a reusable Python application builder, it will require
//	application source code to be added in a particular directory, and it might require a
//	build script to be called after that. You can’t just call ADD and RUN now, because you
//	don’t yet have access to the application source code, and it will be different for each
//	application build. You could simply provide application developers with a boilerplate
//	Dockerfile to copy-paste into their application, but that is inefficient, error-prone
//	and difficult to update because it mixes with application-specific code.
//
//	The solution is to use ONBUILD to register advance instructions to run later, during
//	the next build stage.
//
//	Here’s how it works:
//
//	When it encounters an OnBuild instruction, the builder adds a trigger to the metadata
//	of the image being built. The instruction does not otherwise affect the current build.
//	At the end of the build, a list of all triggers is stored in the image manifest, under
//	the key OnBuild. They can be inspected with the docker inspect command.
//
//	Later the image may be used as a base for a new build, using the FROM instruction.
//	As part of processing the FROM instruction, the downstream builder looks for OnBuild
//	triggers, and executes them in the same order they were registered. If any of the
//	triggers fail, the FROM instruction is aborted which in turn causes the build to fail.
//
//	If all triggers succeed, the FROM instruction completes and the build continues as
//	usual.
//
//	Triggers are cleared from the final image after being executed. In other words they are
//	not inherited by “grand-children” builds.
//
//	For example you might add something like this:
//
//	 []string{
//	   "ADD . /app/src",
//	   "RUN /usr/local/bin/python-build --dir /app/src",
//	 }
//
// Warning:
//
//	The ONBUILD instruction may not trigger FROM or MAINTAINER instructions.
//
// Note:
//
//	See https://docs.docker.com/engine/reference/builder/#onbuild
//
// Português:
//
//	Adiciona à imagem uma instrução de gatilho a ser executada posteriormente, quando a imagem for
//	usada como base para outra construção.
//
//	 Entrada:
//	   onBuild: Lista de instruções de gatilho a serem executadas posteriormente, quando a imagem for
//	     usada como base para outra construção
//
//	O gatilho será executado no contexto do downstream build , como se tivesse sido
//	inserido imediatamente após a instrução FROM no downstream Dockerfile.
//
//	Qualquer instrução de construção pode ser registrada como um gatilho.
//
//	Isso é útil se você estiver construindo uma imagem que será usada como base para
//	construir outras imagens, por exemplo, um ambiente de construção de aplicativo ou um
//	daemon que pode ser personalizado com configuração específica do usuário.
//
//	Por exemplo, se sua imagem for um construtor de aplicativo Python reutilizável, ela
//	exigirá que o código-fonte do aplicativo seja adicionado em um diretório específico e
//	pode exigir que um script de construção seja chamado depois disso. Você não pode
//	simplesmente chamar ADD e RUN agora, porque você ainda não tem acesso ao código-fonte
//	do aplicativo e será diferente para cada construção de aplicativo. Você poderia
//	simplesmente fornecer aos desenvolvedores de aplicativos um Dockerfile padrão para
//	copiar e colar em seus aplicativos, mas isso é ineficiente, sujeito a erros e difícil
//	de atualizar porque se mistura com o código específico do aplicativo.
//
//	A solução é usar o OnBuild para registrar instruções antecipadas para executar mais
//	tarde, durante o próximo estágio de compilação.
//
//	Funciona assim:
//
//	Ao encontrar uma instrução OnBuild, o construtor adiciona um gatilho aos metadados da
//	imagem que está sendo construída. A instrução não afeta de outra forma a construção
//	atual.
//
//	No final da construção, uma lista de todos os gatilhos é armazenada no manifesto da
//	imagem, sob a chave OnBuild. Eles podem ser inspecionados com o comando docker inspect.
//	Posteriormente, a imagem pode ser usada como base para uma nova construção, usando a
//	instrução FROM. Como parte do processamento da instrução FROM, o downstream builder
//	procura gatilhos OnBuild e os executa na mesma ordem em que foram registrados.
//	Se qualquer um dos gatilhos falhar, a instrução FROM é abortada, o que, por sua vez,
//	faz com que o build falhe. Se todos os gatilhos forem bem-sucedidos, a instrução FROM
//	será concluída e a construção continuará normalmente.
//
//	Os gatilhos são apagados da imagem final após serem executados. Em outras palavras,
//	eles não são herdados por construções de "netos".
//
//	Por exemplo, você pode adicionar algo assim:
//
//	 []string{
//	   "ADD . /app/src",
//	   "RUN /usr/local/bin/python-build --dir /app/src",
//	 }
//
// Atenção:
//
//	A instrução ONBUILD não pode disparar as instruções FROM ou MAINTAINER.
//
// Nota:
//
//	https://docs.docker.com/engine/reference/builder/#onbuild
func (e *ContainerBuilder) SetOnBuild(onBuild []string) {
	e.containerConfig.OnBuild = onBuild
}
