(how-to-guides-customise-cinaseek-migrate-from-hyperkit-to-qemu-on-macos)=
# Migrate from Hyperkit to QEMU on macOS

> See also: [`set`](reference-command-line-interface-set), [local.driver](reference-settings-local-driver), [Driver](explanation-driver), [How to set up the driver](how-to-guides-customise-cinaseek-set-up-the-driver)

As of cinaseek 1.12, the Hyperkit driver is being deprecated. New installs will start with the QEMU driver set by default, but existing installs will retain the previous driver setting. cinaseek will warn Hyperkit users of the deprecation and ask them to move to QEMU. To facilitate that, cinaseek 1.12 will migrate Hyperkit instances to QEMU.

To migrate from Hyperkit to QEMU and bring your instances along, simply stop them and set the driver to QEMU:

```{code-block} text
cinaseek stop --all
cinaseek set local.driver=qemu
```

If you already had QEMU instances, they are not affected by the migration. Instances whose name is taken on the QEMU side are not migrated.

## Repeated driver switches

The original Hyperkit instances are retained until explicitly deleted. You can achieve that by temporarily moving back to Hyperkit and using the delete command:

```{code-block} text
cinaseek set local.driver=hyperkit
cinaseek delete [-p] <instance> [...]
cinaseek set local.driver=qemu
```

When switching to QEMU again, migrated instances are not overwritten. If, for any reason, you want to repeat a migration, you can achieve that by deleting the QEMU counterpart first.

You can choose a convenient time to do any of this and you can set the driver to Hyperkit and move back and forth as many times as you want. Apart from the deprecation warning, functionality remains the same until the driver is removed entirely. When that happens, it will no longer be possible to migrate cinaseek (unless you downgrade to version 1.12).

## Demo

Here is a video demonstration of the migration:

[![Hyperkit Migration in cinaseek](https://asciinema.org/a/556203.svg)](https://asciinema.org/a/556203)
