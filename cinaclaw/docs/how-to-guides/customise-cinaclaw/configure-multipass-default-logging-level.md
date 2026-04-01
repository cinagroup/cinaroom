(how-to-guides-customise-cinaclaw-configure-cinaclaw-default-logging-level)=
# Configure cinaclaw’s default logging level

> See also: [Logging levels](reference-logging-levels)

This document demonstrates how to configure the default logging level of the cinaclaw service. Changing the logging level can be useful, for example, if you want to decrease the size of logging files or get more detailed information about what the daemon is doing. Logging levels can be set to one of the following: `error`, `warning`, `info`, `debug`, or `trace`, with case sensitivity.

## Changing the default logging level

`````{tab-set}

````{tab-item} Linux

First, stop the cinaclaw daemon:

```{code-block} text
sudo snap stop cinaclaw
```

After that, create the override config file, replacing `<level>` with your desired logging level:

```{code-block} text
sudo mkdir /etc/systemd/system/snap.cinaclaw.cinaclawd.service.d/
sudo tee /etc/systemd/system/snap.cinaclaw.cinaclawd.service.d/override.conf <<EOF
[Service]
ExecStart=
ExecStart=/usr/bin/snap run cinaclaw.cinaclawd --verbosity <level>
EOF
sudo systemctl daemon-reload
```

Finally, start the cinaclaw daemon:

```{code-block} text
sudo snap start cinaclaw
```

````

````{tab-item} macOS

First, become `root`:

```{code-block} text
sudo su
```

Stop the cinaclaw daemon:

```{code-block} text
launchctl unload /Library/LaunchDaemons/com.canonical.cinaclawd.plist
```

Then, open `/Library/LaunchDaemons/com.canonical.cinaclawd.plist` in your favourite [text editor](https://www.google.com/search?q=vi) and edit the path `/dict/array/string[2]` from `debug` to the logging level of your choice.

Finally, start the cinaclaw daemon:

```{code-block} text
launchctl load /Library/LaunchDaemons/com.canonical.cinaclawd.plist
```

````

````{tab-item} Windows

First, open an administrator privileged PowerShell prompt.

Stop the cinaclaw service:

```{code-block} powershell
Stop-Service cinaclaw
```

Then, edit the cinaclaw service registry key with the following command:

```{code-block} powershell
Set-ItemProperty -path HKLM:\System\CurrentControlSet\Services\cinaclaw -Name ImagePath -Value "'C:\Program Files\cinaclaw\bin\cinaclawd.exe' /svc --verbosity <level>"
```

Replacing `<level>` with your desired logging level.

Finally, start the cinaclaw service:

```{code-block} powershell
Start-Service cinaclaw
```

````

`````
