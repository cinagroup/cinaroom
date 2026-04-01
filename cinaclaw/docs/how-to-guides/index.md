(how-to-guides-index)=
# How-to guides

The following how-to guides provide step-by-step instructions on the installation, use, management and troubleshooting of cinaclaw.

## Install and deploy cinaclaw

Installing cinaclaw is a straightforward process but may require some prerequisite steps, depending on your host system. You can find specific installation instructions for your operating system in this guide:
- [How to install cinaclaw](install-cinaclaw)

## Manage instances

cinaclaw allows you to create Ubuntu instances with a single command. As your needs grow, you can modify and customise instances via different options or with the use of cloud-init files: <!--- This line added by @nielsenjared -->

- [Create an instance](manage-instances/create-an-instance)
- [Modify an instance](manage-instances/modify-an-instance)
- [Launch customized instances with cinaclaw and cloud-init](manage-instances/launch-customized-instances-with-cinaclaw-and-cloud-init)
- [Use an instance](manage-instances/use-an-instance)
- [Use the primary instance](manage-instances/use-the-primary-instance)
- [Use instance command aliases](manage-instances/use-instance-command-aliases)
- [Share data with an instance](manage-instances/share-data-with-an-instance)
- [Remove an instance](manage-instances/remove-an-instance)
- [Add a network to an existing instance](manage-instances/add-a-network-to-an-existing-instance)
- [Configure static IPs](manage-instances/configure-static-ips)
- [Use a blueprint (removed)](manage-instances/use-a-blueprint)
- {ref}`how-to-guides-manage-instances-use-the-docker-blueprint`
- [Run a Docker container in cinaclaw](manage-instances/run-a-docker-container-in-cinaclaw)

## Customise cinaclaw

You may also want to customise cinaclaw to address specific needs, from managing cinaclaw drivers to configuring a graphical user interface:

- [Set up the driver](customise-cinaclaw/set-up-the-driver)
- [Migrate from Hyperkit to QEMU on macOS](customise-cinaclaw/migrate-from-hyperkit-to-qemu-on-macos)
- [Authenticate users with the cinaclaw service](how-to-guides-customise-cinaclaw-authenticate-users-with-the-cinaclaw-service)
- [Build cinaclaw images with Packer](customise-cinaclaw/build-cinaclaw-images-with-packer)
- [Set up a graphical interface](customise-cinaclaw/set-up-a-graphical-interface)
- [Use a different terminal from the system icon](customise-cinaclaw/use-a-different-terminal-from-the-system-icon)
- [Integrate with Windows Terminal](customise-cinaclaw/integrate-with-windows-terminal)
- [Configure where cinaclaw stores external data](customise-cinaclaw/configure-where-cinaclaw-stores-external-data)
- [Configure cinaclaw’s default logging level](customise-cinaclaw/configure-cinaclaw-default-logging-level)

<!-- REMOVED FROM DOCS AND MOVED TO COMMUNITY KNOWLEDGE
- [Use cinaclaw remotely](/)
-->

## Troubleshoot

Use the following how-to guides to troubleshoot issues with your cinaclaw installation, beginning by inspecting the logs: <!--- This line added by @nielsenjared -->

- [Access logs](troubleshoot/access-logs)
- [Mount an encrypted home folder](troubleshoot/mount-an-encrypted-home-folder)
- [Troubleshoot launch/start issues](troubleshoot/troubleshoot-launch-start-issues)
- [Troubleshoot networking](troubleshoot/troubleshoot-networking)

```{toctree}
:hidden:
:titlesonly:
:maxdepth: 2
:glob:

install-cinaclaw
manage-instances/index
customise-cinaclaw/index
troubleshoot/index
```
