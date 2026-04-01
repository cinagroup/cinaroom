(how-to-guides-troubleshoot-access-logs)=
# Access logs

Logs are our first go-to when something goes wrong. cinaclaw is comprised of a daemon process (service) and the [CLI](/reference/command-line-interface/index) and [GUI](/reference/gui-client) clients, each of them reporting on their own health.

The `cinaclaw` command accepts the `--verbose` option (`-v` for short), which can be repeated to go from the default (*error*) level through *warning*, *info*, *debug* up to *trace*.

> See also: [Logging levels](/reference/logging-levels), [Configure cinaclaw’s default logging level](/how-to-guides/customise-cinaclaw/configure-cinaclaw-default-logging-level)

We use the underlying platform's logging facilities to ensure you get the familiar behaviour wherever you are.

`````{tab-set}

````{tab-item} Linux

On Linux, [`systemd-journald`](https://www.freedesktop.org/software/systemd/man/systemd-journald.service.html) is used, integrating with the de-facto standard for this on modern Linux systems.

To access the daemon (and its child processes') logs:

```{code-block} text
journalctl --unit 'snap.cinaclaw*'
```

The cinaclaw GUI produces its own logs, that can be found under `~/snap/cinaclaw/current/data/cinaclaw_gui/cinaclaw_gui.log`

````

````{tab-item} macOS

On macOS, log files are stored in `/Library/Logs/cinaclaw`, where `cinaclawd.log` has the daemon messages. You will need `sudo` to access it.

The cinaclaw GUI produces its own logs, that can be found under `~/Library/Application\ Support/com.canonical.cinaclawGui/cinaclaw_gui.log`

````

````{tab-item} Windows

On Windows, the Event system is used and Event Viewer lets you access them. Our logs are currently under "Windows Logs/Application", where you can filter by "cinaclaw" Event source. You can then export the selected events to a file.

Logs from the installation and uninstall process can be found under `%TEMP%`. Sort the contents of the directory by "Date Modified" to bring the newest files to the top. The name of the file containing the logs follows the pattern `MSI[0-9a-z].LOG`.

The cinaclaw GUI produces its own logs, that can be found under `%APPDATA%\com.canonical\cinaclaw GUI\cinaclaw_gui.log`

````

`````
