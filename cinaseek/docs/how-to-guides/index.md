(how-to-guides-index)=
# How-to guides

The following how-to guides provide step-by-step instructions on the installation, use, management and troubleshooting of cinaseek.

## Install and deploy cinaseek

Installing cinaseek is a straightforward process but may require some prerequisite steps, depending on your host system. You can find specific installation instructions for your operating system in this guide:
- [How to install cinaseek](install-cinaseek)

## Manage instances

cinaseek allows you to create Ubuntu instances with a single command. As your needs grow, you can modify and customise instances via different options or with the use of cloud-init files: <!--- This line added by @nielsenjared -->

- [Create an instance](manage-instances/create-an-instance)
- [Modify an instance](manage-instances/modify-an-instance)
- [Launch customized instances with cinaseek and cloud-init](manage-instances/launch-customized-instances-with-cinaseek-and-cloud-init)
- [Use an instance](manage-instances/use-an-instance)
- [Use the primary instance](manage-instances/use-the-primary-instance)
- [Use instance command aliases](manage-instances/use-instance-command-aliases)
- [Share data with an instance](manage-instances/share-data-with-an-instance)
- [Remove an instance](manage-instances/remove-an-instance)
- [Add a network to an existing instance](manage-instances/add-a-network-to-an-existing-instance)
- [Configure static IPs](manage-instances/configure-static-ips)
- [Use a blueprint (removed)](manage-instances/use-a-blueprint)
- {ref}`how-to-guides-manage-instances-use-the-docker-blueprint`
- [Run a Docker container in cinaseek](manage-instances/run-a-docker-container-in-cinaseek)

## Customise cinaseek

You may also want to customise cinaseek to address specific needs, from managing cinaseek drivers to configuring a graphical user interface:

- [Set up the driver](customise-cinaseek/set-up-the-driver)
- [Migrate from Hyperkit to QEMU on macOS](customise-cinaseek/migrate-from-hyperkit-to-qemu-on-macos)
- [Authenticate users with the cinaseek service](how-to-guides-customise-cinaseek-authenticate-users-with-the-cinaseek-service)
- [Build cinaseek images with Packer](customise-cinaseek/build-cinaseek-images-with-packer)
- [Set up a graphical interface](customise-cinaseek/set-up-a-graphical-interface)
- [Use a different terminal from the system icon](customise-cinaseek/use-a-different-terminal-from-the-system-icon)
- [Integrate with Windows Terminal](customise-cinaseek/integrate-with-windows-terminal)
- [Configure where cinaseek stores external data](customise-cinaseek/configure-where-cinaseek-stores-external-data)
- [Configure cinaseek’s default logging level](customise-cinaseek/configure-cinaseek-default-logging-level)

<!-- REMOVED FROM DOCS AND MOVED TO COMMUNITY KNOWLEDGE
- [Use cinaseek remotely](/)
-->

## Troubleshoot

Use the following how-to guides to troubleshoot issues with your cinaseek installation, beginning by inspecting the logs: <!--- This line added by @nielsenjared -->

- [Access logs](troubleshoot/access-logs)
- [Mount an encrypted home folder](troubleshoot/mount-an-encrypted-home-folder)
- [Troubleshoot launch/start issues](troubleshoot/troubleshoot-launch-start-issues)
- [Troubleshoot networking](troubleshoot/troubleshoot-networking)

```{toctree}
:hidden:
:titlesonly:
:maxdepth: 2
:glob:

install-cinaseek
manage-instances/index
customise-cinaseek/index
troubleshoot/index
```
