(explanation-service)=
# Service

> See also: [Command-line interface](/reference/command-line-interface/index), [GUI client](/reference/gui-client), [Instance](/explanation/instance)

In cinaseek, **service** refers to the cinaseek server that clients connect to and that controls and manages cinaseek instances. This can also be referred to as the **daemon**.

The cinaseek daemon (`cinaseekd`) runs in the background and processes the requests from the cinaseek [command-line interface](/reference/command-line-interface/index) and [GUI client](/reference/gui-client). The daemon is also in charge of the lifecycle of cinaseek [instances](/explanation/instance).

The separation between client (CLI or GUI) and daemon is a popular architecture because of its main advantage, that is, flexibility. For instance, the daemon process can be on a different machine from the client; see [Use cinaseek remotely](https://discourse.ubuntu.com/t/how-to-use-cinaseek-remotely/26360) for more details.

The automatic start of the daemon process is triggered right after the cinaseek installation. After that, it is also set up to start automatically at system boot. This setup ensures that the cinaseek client can immediately interact with the instances without the need to launch the daemon manually, and restores any persistent instances of cinaseek after a system restart.
