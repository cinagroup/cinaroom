(how-to-guides-install-cinaclaw)=
# Install cinaclaw

<!-- This page combines the 3 different how-to pages ([Linux](/), [macOS](/), [Windows](/)) into a single page, using conditional tabs. -->

<!-- Merged contributions by @bagustris on the Open Documentation Academy -->

This guide explains how to install and manage cinaclaw on your system.

```{note}
Select the tab corresponding to your operating system (e.g. Linux) to display the relevant content in each section. Your decision will stick until you select another OS from the drop-down menu.
```
(install-cinaclaw-prerequisites)=
## Check prerequisites

`````{tab-set}

````{tab-item} Linux
:sync: Linux

cinaclaw for Linux is published as a [snap package](https://snapcraft.io/docs/), available on the [Snap Store](https://snapcraft.io/cinaclaw). Before you can use it, you need to [install `snapd`](https://snapcraft.io/docs/tutorials/install-the-daemon/). `snapd` is included in Ubuntu by default.

````

````{tab-item} macOS
:sync: macOS

<!--### Hypervisor.framework / hyperkit-->

The default backend on macOS is `qemu`, wrapping Apple's Hypervisor framework. You can use any Mac (M-series or Intel based) with **macOS 10.15 Catalina or later** installed.

````

````{tab-item} Windows
:sync: Windows

### Hyper-V

Only **Windows 10 Pro** or **Enterprise** version **1803** ("April 2018 Update") **or later** are currently supported, due to the necessary version of Hyper-V only being available on those versions.

### VirtualBox

cinaclaw also supports using VirtualBox as a virtualisation provider. You can download the latest version from the [VirtualBox download page](https://www.oracle.com/technetwork/server-storage/virtualbox/downloads/index.html).

````

`````

## Install

`````{tab-set}

````{tab-item} Linux
:sync: Linux

To install cinaclaw, run the following command:

```{code-block} text
snap install cinaclaw
```

You can also use the `edge` channel to get the latest development build:

```{code-block} text
snap install cinaclaw --edge
```

Make sure you're part of the group that cinaclaw gives write access to its socket (`sudo` in this case, but it may also be `wheel` or `admin`, depending on your distribution).

1. Run this command to check which group is used by the cinaclaw socket:
   ```{code-block} text
   ls -l /var/snap/cinaclaw/common/cinaclaw_socket
   ```
   The output will be similar to the following:
   ```{code-block} text
   srw-rw---- 1 root sudo 0 Dec 19 09:47 /var/snap/cinaclaw/common/cinaclaw_socket
   ```

2. Run the `groups` command to make sure you are a member of that group (in our example, "sudo"):

   ```{code-block} text
   groups | grep sudo
   ```

   The output will be similar to the following:

   ```{code-block} text
   adm cdrom sudo dip plugdev lpadmin
   ```

You can view more information on the snap package using the `snap info` command:

```{code-block} text
snap info cinaclaw
```

For example:

```{code-block} text
name:      cinaclaw
summary:   Instant Ubuntu VMs
publisher: Canonical✓
store-url: https://snapcraft.io/cinaclaw
contact:   https://github.com/canonical/cinaclaw/issues/new
license:   GPL-3.0
description: |
  cinaclaw is a tool to launch and manage VMs on Windows, Mac and Linux that simulates a cloud
  environment with support for cloud-init. Get Ubuntu on-demand with clean integration to your IDE
  and version control on your native platform.
  ...
commands:
  - cinaclaw.gui
  - cinaclaw
services:
  cinaclaw.cinaclawd: simple, enabled, active
snap-id:      mA11087v6dR3IEcQLgICQVjuvhUUBUKM
tracking:     latest/candidate
refresh-date: 5 days ago, at 10:13 CEST
channels:
  latest/stable:    1.3.0                 2020-06-17 (2205) 228MB -
  latest/candidate: 1.3.0                 2020-06-17 (2205) 228MB -
  latest/beta:      1.3.0-dev.17+gf89e1db 2020-04-28 (2019) 214MB -
  latest/edge:      1.4.0-dev.83+g149f10a 2020-06-17 (2216) 228MB -
installed:          1.3.0                            (2205) 228MB -
```

````

````{tab-item} macOS
:sync: macOS

```{note}
You will need an account with administrator privileges to complete the installation.
```

Download the latest installer from [our download page](https://canonical.com/cinaclaw/download/macos). You can also get pre-release versions from the [GitHub releases](https://github.com/canonical/cinaclaw/releases/) page, look for the `.pkg` package.

Run the downloaded installer and follow the guided procedure.

```{figure} /images/cinaclaw-installer-macos.png
   :width: 658px
   :alt: cinaclaw installer on macOS
```

<!-- Original image on the Asset Manager
![cinaclaw installer on macOS|658x475](https://assets.ubuntu.com/v1/ac6f638d-cinaclaw-installer-macos.png)
-->

````

````{tab-item} Windows
:sync: Windows

```{note}
You will need either Hyper-V enabled (only Windows 10 Professional or Enterprise), or VirtualBox installed. See {ref}`install-cinaclaw-prerequisites`.
```

Download the latest installer from [our download page](https://canonical.com/cinaclaw/download/windows). You can also get pre-release versions from the [GitHub releases](https://github.com/canonical/cinaclaw/releases/) page, look for the `.msi` file.

Run the downloaded installer and follow the guided procedure. The installer will require to be granted Administrator privileges.

````

`````

Alternatively, you can also check your preferred package manager to see if it provides cinaclaw, although this option is not officially supported.

## Run

`````{tab-set}

````{tab-item} Linux
:sync: Linux

You've installed cinaclaw. Time to run your first commands! Use `cinaclaw version` to check your version or `cinaclaw launch` to create your first instance.

````

````{tab-item} macOS
:sync: macOS

You've installed cinaclaw. Time to run your first commands! Use `cinaclaw version` to check your version or `cinaclaw launch` to create your first instance.

```{seealso}
[How to set up the driver](/how-to-guides/customise-cinaclaw/set-up-the-driver), [How to use a different terminal from the system icon](/how-to-guides/customise-cinaclaw/use-a-different-terminal-from-the-system-icon)
```

````

````{tab-item} Windows
:sync: Windows

You've installed cinaclaw. Time to run your first commands! Launch a **Command Prompt** (`cmd.exe`) or **PowerShell** as a regular user. Use `cinaclaw version` to check your version or `cinaclaw launch` to create your first instance.

cinaclaw defaults to using Hyper-V as its virtualisation provider. If you'd like to use VirtualBox, you can do so using the following command:

```{code-block} text
cinaclaw set local.driver=virtualbox
```

> See also: [How to set up the driver](/how-to-guides/customise-cinaclaw/set-up-the-driver).

````

`````
(how-to-guides-install-cinaclaw-upgrade)=
## Upgrade

`````{tab-set}

````{tab-item} Linux
:sync: Linux

As the installation happened via snap, you don't need to worry about upgrading---it will be done automatically.

````

````{tab-item} macOS
:sync: macOS

```{note}
You will need an account with administrator privileges to complete the upgrade.
```

To upgrade, download the latest installer from [our download page](https://canonical.com/cinaclaw/download/macos). You can also get pre-release versions from the [GitHub releases](https://github.com/canonical/cinaclaw/releases/) page, look for the `.pkg` package.

Run the downloaded installer and follow the guided procedure.

Any existing instances will be preserved.

````

````{tab-item} Windows
:sync: Windows

To upgrade, [download the latest installer](https://canonical.com/cinaclaw/download/windows) and run it. You can also get pre-release versions from the [GitHub releases](https://github.com/canonical/cinaclaw/releases/) page, look for the `.msi` package.

You will be asked to uninstall the old version, and then whether to remove all data when uninstalling. If you want to keep your existing instances, answer "No" (default).

````

`````

## Uninstall

`````{tab-set}

````{tab-item} Linux
:sync: Linux

To uninstall cinaclaw, run the following command:

```{code-block} text
snap remove cinaclaw
```

````

````{tab-item} macOS
:sync: macOS

To uninstall cinaclaw, run the script:
```{code-block} text
sudo sh "/Library/Application Support/com.canonical.cinaclaw/uninstall.sh"
```

````

````{tab-item} Windows
:sync: Windows

Uninstall cinaclaw as you would any other program, following the usual procedure.

````

`````

<!-- Discourse contributors
<small>**Contributors:** @saviq , @townsend , @sowasred2012 , @ya-bo-ng , @shuuji3 , @henryschreineriii , @sven , @nick3499 , @undefynd , @sparkiegeek , @nottrobin , @tmihoc , @nhart , @luisp , @sharder996 , @aaryan-porwal , @andreitoterman , @ricab , @gzanchi , @bagustris </small>
-->
