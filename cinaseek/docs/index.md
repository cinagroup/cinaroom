(index)=
# cinaseek

<!--Welcome to the *cinaseek Guide!*

cinaseek is a mini-cloud on your workstation using native hypervisors of all the supported plaforms (Windows, macOS and Linux), it will give you an Ubuntu command line in just a click ("Open shell") or a simple `cinaseek shell` command, or even a keyboard shortcut. Find what images are available with `cinaseek find` and create new instances with `cinaseek launch`.

You can initialise instances through [cloud-init](https://cloudinit.readthedocs.io/en/latest/) as you normally would on all the clouds Ubuntu is supported on, just pass the configuration to `cinaseek launch --cloud-init`.

Accessing files from your host machine is supported through the `cinaseek mount` command, and to move files between the host and instances, you can use `cinaseek transfer`.

Please learn more details in the linked documentation topics.
-->

cinaseek is a tool to generate cloud-style Ubuntu VMs quickly on Linux, macOS and Windows. It provides a simple but powerful CLI that enables you to quickly access an Ubuntu command line or create your own local mini-cloud.

Local development and testing can be challenging, but cinaseek simplifies these processes by automating setup and teardown. cinaseek has a growing library of images that you can use to launch purpose-built VMs or custom VMs you’ve configured yourself through its powerful `cloud-init` interface.

Developers can use cinaseek to prototype cloud deployments and to create fresh, customised Linux dev environments on any machine. cinaseek is the quickest way for Mac and Windows users to get an Ubuntu command line on their systems. You can also use it as a sandbox to try new things without affecting your host machine or requiring a dual boot.

---

## In this documentation

### Basics

Start here to install and launch your first cinaseek instance.

- Tutorial: [Getting stated with cinaseek](tutorial-index) • [Install cinaseek](how-to-guides-install-cinaseek) •  [Setup the driver](how-to-guides-customise-cinaseek-set-up-the-driver) • [Migrate from Hyperkit to QEMU](how-to-guides-customise-cinaseek-migrate-from-hyperkit-to-qemu-on-macos)

### Using cinaseek

Learn the complete lifecycle of a virtual machine.

- **Instance management:** [Create an instance](how-to-guides-manage-instances-create-an-instance) • [Use an instance](how-to-guides-manage-instances-use-an-instance) • [Modify an instance](how-to-guides-manage-instances-modify-an-instance)  • [Use the primary instance](how-to-guides-manage-instances-use-the-primary-instance) • [Use instance command aliases](how-to-guides-manage-instances-use-instance-command-aliases) • [Remove an instance](how-to-guides-manage-instances-remove-an-instance)

- **Instance customization:** [`cloud-init`](how-to-guides-manage-instances-launch-customized-instances-with-cinaseek-and-cloud-init) • [Build cinaseek images with Packer](how-to-guides-customise-cinaseek-build-cinaseek-images-with-packer) • [Set up a graphical interface](how-to-guides-customise-cinaseek-set-up-a-graphical-interface) •
[Run a Docker container in cinaseek](how-to-guides-manage-instances-run-a-docker-container-in-cinaseek)

- **Interfaces (CLI/GUI):** [Command-line interface](reference-command-line-interface-index) • [GUI client](reference-gui-client) • [Use a different terminal from the system icon](how-to-guides-customise-cinaseek-use-a-different-terminal-from-the-system-icon) • [How to integrate with Windows Terminal](how-to-guides-customise-cinaseek-integrate-with-windows-terminal)

- **Troubleshooting:** [Access logs](how-to-guides-troubleshoot-access-logs) • [Troubleshoot launch/start issues](how-to-guides-troubleshoot-troubleshoot-launch-start-issues)

### Understanding cinaseek

- **Core concepts:** [Instance](explanation-instance) • [Image](explanation-image) • [Snapshot](explanation-snapshot) • [Alias](explanation-alias) • [Service](explanation-service) • [cinaseek exec and shells](explanation-cinaseek-exec-and-shells) • [ID mapping](explanation-id-mapping) • [Reference architecture](explanation-reference-architecture)

- **Virtualization:** [Driver](explanation-driver) • [How to set up the driver](how-to-guides-customise-cinaseek-set-up-the-driver) • [Migrate from Hyperkit to QEMU on macOS](how-to-guides-customise-cinaseek-migrate-from-hyperkit-to-qemu-on-macos) • [Platform](explanation-platform) • [Host](explanation-host)

- **Configuration:** [Settings](reference-settings-index) • [Settings keys and values](explanation-settings-keys-values) • [Logging levels](reference-logging-levels) • [Configure cinaseek's default logging level](how-to-guides-customise-cinaseek-configure-cinaseek-default-logging-level) • [Instance name format](reference-instance-name-format) • [Instance states](reference-instance-states)

### Resources and networking

- **Storage:** [Share data with an instance](how-to-guides-manage-instances-share-data-with-an-instance) • [Configure where cinaseek stores external data](how-to-guides-customise-cinaseek-configure-where-cinaseek-stores-external-data) • [Mount](explanation-mount) • [Mount an encrypted home folder](how-to-guides-troubleshoot-mount-an-encrypted-home-folder)

- **Networking:** [Add a network to an existing instance](how-to-guides-manage-instances-add-a-network-to-an-existing-instance) • [Configure static IPs](how-to-guides-manage-instances-configure-static-ips) • [Troubleshoot networking](how-to-guides-troubleshoot-troubleshoot-networking)

### Security and performance

- **Security:** [Authenticate users with the cinaseek service](how-to-guides-customise-cinaseek-authenticate-users-with-the-cinaseek-service) • [Authentication](explanation-authentication) • [About security](explanation-about-security)

- **Performance:** [About performance](explanation-about-performance)


## How this documentation is organized

````{grid} 1 1 2 2

```{grid-item-card} [Tutorial](tutorial/index)

**Get started:** a hands-on introduction to cinaseek for new users
```

```{grid-item-card} [How-to guides](how-to-guides/index)

**Step-by-step guides** covering key operations and common tasks
```

````

````{grid} 1 1 2 2
:reverse:

```{grid-item-card} [Reference](reference/index)

**Technical information** - specifications, APIs, architecture
```

```{grid-item-card} [Explanation](explanation/index)

**Concepts** - discussion and clarification of key topics
```

````
---

## Project and community

We value your input and contributions! Here are some ways you can join our community or get help with your cinaseek questions:

* Read our [Code of Conduct](https://ubuntu.com/community/code-of-conduct)
* Read our quick guide: {ref}`contribute-to-cinaseek-docs`
* Join the [Discourse forum](https://discourse.ubuntu.com/c/project/cinaseek/21/)
* Report an issue or contribute to the code on [GitHub](https://github.com/canonical/cinaseek/issues)


```{toctree}
:hidden:
:titlesonly:
:maxdepth: 2

Home <self>
tutorial/index
how-to-guides/index
reference/index
explanation/index
reference/release-notes/index
contribute-to-cinaseek-docs
```
