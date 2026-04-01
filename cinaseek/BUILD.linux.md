# Build instructions for Linux

## Environment Setup

### Build dependencies

```
cd <cinaseek>
sudo apt install devscripts equivs
mk-build-deps -s sudo -i
```

## Building

First, go into the repository root and get all the submodules:

```
cd <cinaseek>
git submodule update --init --recursive
```

If building on arm, s390x, ppc64le, or riscv, you will need to set the `VCPKG_FORCE_SYSTEM_BINARIES`  environment
variable:

```
export VCPKG_FORCE_SYSTEM_BINARIES=1
```

Then create a build directory and run CMake.

```
mkdir build
cd build
cmake ../
```

This will fetch all necessary content, build vcpkg dependencies, and initialize the build system. You can also specify
the `-DCMAKE_BUILD_TYPE` option to set the build type (e.g., `Debug`, `Release`, `Coverage`, etc.).

To use a different vcpkg, pass `-DMULTIPASS_VCPKG_LOCATION="path/to/vcpkg"` to CMake.
It should point to the root vcpkg location, where the top bootstrap scripts are located.

Finally, to build the project, run:

```
cmake --build . --parallel
```

Please note that if you're working on a forked repository that you created using the "Copy the main branch only" option,
the repository will not include the necessary git tags to determine the cinaseek version during CMake configuration. In
this case, you need to manually fetch the tags from the upstream by running
`git fetch --tags https://github.com/canonical/cinaseek.git` in the `<cinaseek>` source code directory.

## Run the cinaseek daemon and client

First, install cinaseek's runtime dependencies. On AMD64 architecture, you can do this with:

```
sudo apt update
sudo apt install libgl1 libpng16-16 libxml2 dnsmasq-base \
    dnsmasq-utils qemu-utils libslang2 iproute2 iptables \
    iputils-ping libatm1 libxtables12 xterm
```

On ARM64 architecture, you can do this by running:

```
sudo apt update
sudo apt install libgl1 libpng16-16 libxml2 dnsmasq-base \
    dnsmasq-utils qemu-efi-aarch64 qemu-utils libslang2 \
    iproute2 iptables iputils-ping libatm1 libxtables12 \
    xterm
```

You will also need to install your CPU architecture's variant of `qemu-system`. For example, you will need

```
sudo apt install qemu-system-x86
```

on x86_64 machines.

Additionally, on ARM64 architecture, there is an extra step to set up the `QEMU_EFI.fd` file:

```
sudo cp /usr/share/qemu-efi-aarch64/QEMU_EFI.fd /usr/share/qemu/QEMU_EFI.fd
```

Then run the cinaseek daemon:

```
sudo <cinaseek>/build/bin/cinaseekd &
```

Copy the desktop file that cinaseek clients expect to find in your home:

```
mkdir -p ~/.local/share/cinaseek/
cp <cinaseek>/src/client/gui/assets/cinaseek.gui.autostart.desktop ~/.local/share/cinaseek/
```

Optionally, enable auto-complete in Bash:

```
source <cinaseek>/completions/bash/cinaseek
```

To be able to use the binaries without specifying their path:

```
export PATH=<cinaseek>/build/bin
```

Now you can use the `cinaseek` command from your terminal (for example
`<cinaseek>/build/bin/cinaseek launch --name foo`) or launch the GUI client with the command
`<cinaseek>/build/bin/cinaseek.gui`.
