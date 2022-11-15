package docker

// SetOnBuild (English): The ONBUILD instruction adds to the image a trigger instruction
// to be executed at a later time, when the image is used as the base for another build.
// The trigger will be executed in the context of the downstream build, as if it had been
// inserted immediately after the FROM instruction in the downstream Dockerfile.
//
// Any build instruction can be registered as a trigger.
//
// This is useful if you are building an image which will be used as a base to build other
// images, for example an application build environment or a daemon which may be
// customized with user-specific configuration.
//
// For example, if your image is a reusable Python application builder, it will require
// application source code to be added in a particular directory, and it might require a
// build script to be called after that. You can’t just call ADD and RUN now, because you
// don’t yet have access to the application source code, and it will be different for each
// application build. You could simply provide application developers with a boilerplate
// Dockerfile to copy-paste into their application, but that is inefficient, error-prone
// and difficult to update because it mixes with application-specific code.
//
// The solution is to use ONBUILD to register advance instructions to run later, during
// the next build stage.
//
// Here’s how it works:
//
// When it encounters an ONBUILD instruction, the builder adds a trigger to the metadata
// of the image being built. The instruction does not otherwise affect the current build.
// At the end of the build, a list of all triggers is stored in the image manifest, under
// the key OnBuild. They can be inspected with the docker inspect command.
// Later the image may be used as a base for a new build, using the FROM instruction. As
// part of processing the FROM instruction, the downstream builder looks for ONBUILD
// triggers, and executes them in the same order they were registered. If any of the
// triggers fail, the FROM instruction is aborted which in turn causes the build to fail.
// If all triggers succeed, the FROM instruction completes and the build continues as
// usual.
// Triggers are cleared from the final image after being executed. In other words they are
// not inherited by “grand-children” builds.
// For example you might add something like this:
//
// ONBUILD ADD . /app/src
// ONBUILD RUN /usr/local/bin/python-build --dir /app/src
//
//	Warning:
//	Chaining ONBUILD instructions using ONBUILD ONBUILD isn’t allowed.
//
//	Warning:
//	The ONBUILD instruction may not trigger FROM or MAINTAINER instructions.
//
// https://docs.docker.com/engine/reference/builder/#onbuild
func (el *DockerSystem) SetOnBuild(onBuild []string) {
	el.onBuild = onBuild
}
