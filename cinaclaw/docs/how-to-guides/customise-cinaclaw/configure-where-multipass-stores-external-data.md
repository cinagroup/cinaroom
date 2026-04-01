(how-to-guides-customise-cinaclaw-configure-where-cinaclaw-stores-external-data)=
# Configure where cinaclaw stores external data

This document demonstrates how to configure the location where cinaclaw stores instances, caches images, and other data. Configuring a new storage location can be useful, for example, if you need to free up storage space on your boot partition.

(configuring-a-new-storage-location)=
## Configuring a new storage location

```{caution}
**Caveats:**
- cinaclaw will not migrate your existing data; this article explains how to do it manually. If you do not transfer the data, you will have to re-download any Ubuntu images and reinitialise any instances that you need.
- When uninstalling cinaclaw, the uninstaller will not remove data stored in custom locations, so you'll have to delete it manually.
```

`````{tab-set}

````{tab-item} Linux
:sync: Linux

First, stop the cinaclaw daemon:

```{code-block} text
sudo snap stop cinaclaw
```

Since cinaclaw is installed using a strictly confined snap, it is limited on what it can do or access on your host. Depending on where the new storage directory is located, you will need to connect the respective interface to the cinaclaw snap. Because of [snap confinement](https://snapcraft.io/docs/snap-confinement), this directory needs to be located in either your home directory (`~`, e.g. `/home/username/`, which is connected by default) or one of the removable mounts points (`/mnt` or `/media`). To connect the removable mount points, use the command:

  ```{code-block} text
  sudo snap connect cinaclaw:removable-media
  ```

Create the new directory in which cinaclaw will store its data:

```{code-block} text
mkdir -p <path>
sudo chown root <path>
```

After that, create the override config file, replacing `<path>` with the absolute path of the directory created above.

```{code-block} text
sudo mkdir /etc/systemd/system/snap.cinaclaw.cinaclawd.service.d/
sudo tee /etc/systemd/system/snap.cinaclaw.cinaclawd.service.d/override.conf <<EOF
[Service]
Environment=MULTIPASS_STORAGE=<path>
EOF
```

The output at this point will be:
```{code-block} text
[Service]
Environment=MULTIPASS_STORAGE=<path>
```

Then, instruct `systemd` to reload the daemon configuration files:

```{code-block} text
sudo systemctl daemon-reload
```

Now you can transfer the data from its original location to the new location:

```{code-block} text
sudo cp -r /var/snap/cinaclaw/common/data/cinaclawd <path>/data
sudo cp -r /var/snap/cinaclaw/common/cache/cinaclawd <path>/cache
```

<!-- The following step was added to address GitHub issue https://github.com/canonical/cinaclaw/issues/3254 -->

You also need to edit the following configuration files so that the specified paths point to the new cinaclaw storage directory, otherwise your instances will fail to start:

* `cinaclaw-vm-instances.json`: Update the absolute path of the instance images in the "arguments" key for each instance.
* `vault/cinaclawd-instance-image-records.json`: Update the "path" key for each instance.

Finally, start the cinaclaw daemon:

```{code-block} text
sudo snap start cinaclaw
```

You can delete the original data at your discretion, to free up space:

```{code-block} text
sudo rm -rf /var/snap/cinaclaw/common/data/cinaclawd/vault
sudo rm -rf /var/snap/cinaclaw/common/cache/cinaclawd
```

````

````{tab-item} macOS
:sync: macOS

First, become `root`:

```{code-block} text
sudo su
```

Stop the cinaclaw daemon:

```{code-block} text
launchctl unload /Library/LaunchDaemons/com.canonical.cinaclawd.plist
```

Move your current data from its original location to `<path>`, replacing `<path>` with your custom location of choice:

```{code-block} text
mv /var/root/Library/Application\ Support/cinaclawd <path>
```

```{caution}
Make sure the `cinaclawd` directory is moved to `<path>`, and not inside the  `<path>` folder.
```

Define a symbolic link from the original location to the absolute path of new location:

```{code-block} text
ln -s <path> /var/root/Library/Application\ Support/cinaclawd
```

Finally, start the cinaclaw daemon:

```{code-block} text
launchctl load /Library/LaunchDaemons/com.canonical.cinaclawd.plist
```

````

````{tab-item} Windows
:sync: Windows

First, open a PowerShell prompt with administration privileges.

Stop cinaclaw instances:

```{code-block} powershell
cinaclaw stop --all
```

Stop the cinaclaw daemon:

```{code-block} powershell
Stop-Service cinaclaw
```

Create and set the new storage location, replacing `<path>` with the absolute path of your choice:

```{code-block} powershell
New-Item -ItemType Directory -Path "<path>"
Set-ItemProperty -Path "HKLM:System\CurrentControlSet\Control\Session Manager\Environment" -Name MULTIPASS_STORAGE -Value "<path>"
```

Now you can transfer the data from its original location to the new location:

```{code-block} powershell
Copy-Item -Path "C:\ProgramData\cinaclaw\*" -Recurse -Force -Destination "<path>"
```

```{caution}
It is important to copy any existing data to the new location. This avoids unauthenticated user issues, permission issues, and in general, to have any previously created instances available.
```

You also need to edit several settings so that the specified paths point to the new cinaclaw storage directory, otherwise your instances will fail to start:

* `<path>/data/vault/cinaclawd-instance-image-records.json`: Update the "path" key for each instance.
* Open Hyper-V Manager > For each instance right-click the instance and open the settings. Navigate to SCSI Controller > Hard Drive and update the Media path. Do the same for SCSI Controller > DVD Drive > Media Image file.

Finally, start the cinaclaw daemon:

```{code-block} powershell
Start-Service cinaclaw
```

You can delete the original data at your discretion, to free up space:

```{code-block} powershell
Remove-Item -Path "C:\ProgramData\cinaclaw\cache\*" -Recurse
Remove-Item -Path "C:\ProgramData\cinaclaw\data\vault\*" -Recurse
```

````

`````

## Reverting back to the default location

`````{tab-set}

````{tab-item} Linux
:sync: Linux

Stop the cinaclaw daemon:

```{code-block} text
sudo snap stop cinaclaw
```

Although not required, to make sure that cinaclaw does not have access to directories that it shouldn't, you can disconnect the respective interface depending on where the custom storage location was set (see {ref}`configuring-a-new-storage-location` above).
For example, to disconnect the removable mounts points (`/mnt` or `/media`), run:

```{code-block} text
sudo snap disconnect cinaclaw:removable-media
```

Then, remove the override config file:

```{code-block} text
sudo rm /etc/systemd/system/snap.cinaclaw.cinaclawd.service.d/override.conf
sudo systemctl daemon-reload
```

Now you can transfer your data from the custom location back to its original location:

```{code-block} text
sudo cp -r <path>/data /var/snap/cinaclaw/common/data/cinaclawd
sudo cp -r <path>/cache /var/snap/cinaclaw/common/cache/cinaclawd
```

You also need to edit the following configuration files so that the specified paths point to the original cinaclaw storage directory, otherwise your instances will fail to start:

* `cinaclaw-vm-instances.json`: Update the absolute path of the instance images in the "arguments" key for each instance.
* `vault/cinaclawd-instance-image-records.json`: Update the "path" key for each instance.

Finally, start the cinaclaw daemon:

```{code-block} text
sudo snap start cinaclaw
```

You can delete the data from the custom location at your discretion, to free up space:

```{code-block} text
sudo rm -rf <path>
```

````

````{tab-item} macOS
:sync: macOS

First, become `root`:

```{code-block} text
sudo su
```

Stop the cinaclaw daemon:

```{code-block} text
launchctl unload /Library/LaunchDaemons/com.canonical.cinaclawd.plist
```

Remove the link pointing to your custom location:

```{code-block} text
unlink /var/root/Library/Application\ Support/cinaclawd
```

Move the data from your custom location back to its original location:

```{code-block} text
mv <path> /var/root/Library/Application\ Support/cinaclawd
```

Finally, start the cinaclaw daemon:

```{code-block} text
launchctl load /Library/LaunchDaemons/com.canonical.cinaclawd.plist
```

````

````{tab-item} Windows
:sync: Windows

First, open a PowerShell prompt with administrator privileges.

Stop cinaclaw instances:

```{code-block} powershell
cinaclaw stop --all
```

Stop the cinaclaw daemon:

```{code-block} powershell
Stop-Service cinaclaw
```

Remove the setting for the custom storage location:

```{code-block} powershell
Remove-ItemProperty -Path "HKLM:System\CurrentControlSet\Control\Session Manager\Environment" -Name MULTIPASS_STORAGE
```

Now you can transfer the data back to its original location:

```{code-block} powershell
Copy-Item -Path "<path>\*" -Destination "C:\ProgramData\cinaclaw" -Recurse -Force
```

Follow the same instructions from setting up the custom image location to update the paths to their original location.

Finally, start the cinaclaw daemon:

```{code-block} powershell
Start-Service cinaclaw
```

You can delete the data from the custom location at your discretion, to free up space:

```{code-block} powershell
Remove-Item -Path "<path>" -Recurse -Force
```

````

`````
