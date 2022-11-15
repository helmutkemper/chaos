package docker

//User memory constraints🔗
//We have four ways to set user memory usage:
//
//Option	Result
//memory=inf, memory-swap=inf (default)	There is no memory limit for the container. The container can use as much memory as needed.
//memory=L<inf, memory-swap=inf	(specify memory and set memory-swap as -1) The container is not allowed to use more than L bytes of memory, but can use as much swap as is needed (if the host supports swap memory).
//memory=L<inf, memory-swap=2*L	(specify memory without memory-swap) The container is not allowed to use more than L bytes of memory, swap plus memory usage is double of that.
//memory=L<inf, memory-swap=S<inf, L<=S	(specify both memory and memory-swap) The container is not allowed to use more than L bytes of memory, swap plus memory usage is limited by S.
//Examples:
//
//$ docker run -it ubuntu:14.04 /bin/bash
//We set nothing about memory, this means the processes in the container can use as much memory and swap memory as they need.
//
//$ docker run -it -m 300M --memory-swap -1 ubuntu:14.04 /bin/bash
//We set memory limit and disabled swap memory limit, this means the processes in the container can use 300M memory and as much swap memory as they need (if the host supports swap memory).
//
//$ docker run -it -m 300M ubuntu:14.04 /bin/bash
//We set memory limit only, this means the processes in the container can use 300M memory and 300M swap memory, by default, the total virtual memory size (--memory-swap) will be set as double of memory, in this case, memory + swap would be 2*300M, so processes can use 300M swap memory as well.
//
//$ docker run -it -m 300M --memory-swap 1G ubuntu:14.04 /bin/bash
//We set both memory and swap memory, so the processes in the container can use 300M memory and 700M swap memory.
//
//Memory reservation is a kind of memory soft limit that allows for greater sharing of memory. Under normal circumstances, containers can use as much of the memory as needed and are constrained only by the hard limits set with the -m/--memory option. When memory reservation is set, Docker detects memory contention or low memory and forces containers to restrict their consumption to a reservation limit.
//
//Always set the memory reservation value below the hard limit, otherwise the hard limit takes precedence. A reservation of 0 is the same as setting no reservation. By default (without reservation set), memory reservation is the same as the hard memory limit.
//
//Memory reservation is a soft-limit feature and does not guarantee the limit won’t be exceeded. Instead, the feature attempts to ensure that, when memory is heavily contended for, memory is allocated based on the reservation hints/setup.
//
//The following example limits the memory (-m) to 500M and sets the memory reservation to 200M.
//
//$ docker run -it -m 500M --memory-reservation 200M ubuntu:14.04 /bin/bash
//Under this configuration, when the container consumes memory more than 200M and less than 500M, the next system memory reclaim attempts to shrink container memory below 200M.
//
//The following example set memory reservation to 1G without a hard memory limit.
//
//$ docker run -it --memory-reservation 1G ubuntu:14.04 /bin/bash
//The container can use as much memory as it needs. The memory reservation setting ensures the container doesn’t consume too much memory for long time, because every memory reclaim shrinks the container’s consumption to the reservation.
//
//By default, kernel kills processes in a container if an out-of-memory (OOM) error occurs. To change this behaviour, use the --oom-kill-disable option. Only disable the OOM killer on containers where you have also set the -m/--memory option. If the -m flag is not set, this can result in the host running out of memory and require killing the host’s system processes to free memory.
//
//The following example limits the memory to 100M and disables the OOM killer for this container:
//
//$ docker run -it -m 100M --oom-kill-disable ubuntu:14.04 /bin/bash
//The following example, illustrates a dangerous way to use the flag:
//
//$ docker run -it --oom-kill-disable ubuntu:14.04 /bin/bash
//The container has unlimited memory which can cause the host to run out memory and require killing system processes to free memory. The --oom-score-adj parameter can be changed to select the priority of which containers will be killed when the system is out of memory, with negative scores making them less likely to be killed, and positive scores more likely.

// SetImageBuildOptionsMemory
//
// English:
//
//	The maximum amount of memory the container can use.
//
//	 Input:
//	   value: amount of memory in bytes
//
// Note:
//
//   - If you set this option, the minimum allowed value is 4 * 1024 * 1024 (4 megabyte);
//   - Use value * KKiloByte, value * KMegaByte and value * KGigaByte
//     See https://docs.docker.com/engine/reference/run/#user-memory-constraints
//
// Português:
//
//	Memória máxima total que o container pode usar.
//
//	 Entrada:
//	   value: Quantidade de memória em bytes
//
// Nota:
//
//   - Se você vai usar esta opção, o máximo permitido é 4 * 1024 * 1024 (4 megabyte)
//   - Use value * KKiloByte, value * KMegaByte e value * KGigaByte
//     See https://docs.docker.com/engine/reference/run/#user-memory-constraints
func (e *ContainerBuilder) SetImageBuildOptionsMemory(value int64) {
	e.buildOptions.Memory = value

	e.addProblem("The SetImageBuildOptionsMemory() function can generate an error when building the image.")
}
