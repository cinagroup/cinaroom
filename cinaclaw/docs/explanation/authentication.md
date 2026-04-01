(explanation-authentication)=
# Authentication

> See also: [How to authenticate users with the cinaclaw service](how-to-guides-customise-cinaclaw-authenticate-users-with-the-cinaclaw-service)

Before executing any commands, cinaclaw requires users to authenticate with the service. cinaclaw employs an authentication process based on x509 certificates signed by elliptic curve (EC) keys, powered by OpenSSL, to authenticate users. When a user connects, cinaclaw validates the certificate to ensure only verified users can access the service.

`````{tab-set}

````{tab-item} Linux
Linux and macOS hosts currently use a Unix domain socket for client and daemon communication. Upon first use, this socket only allows a client to connect via a user belonging to the group that owns the socket. For example, this group could be `sudo`, `admin`, or `wheel` and the user needs to belong to this group or else permission will be denied when connecting.

After the first client connects with a user belonging to the socket's admin group, the user's OpenSSL certificate will be accepted by the daemon and the socket will then be open for all users to connect. Any other user trying to connect to the cinaclaw service will need to authenticate with the service using the previously set [`local.passphrase`](/reference/settings/local-passphrase).
````

````{tab-item} macOS
Linux and macOS hosts currently use a Unix domain socket for client and daemon communication. Upon first use, this socket only allows a client to connect via a user belonging to the group that owns the socket. For example, this group could be `sudo`, `admin`, or `wheel` and the user needs to belong to this group or else permission will be denied when connecting.

After the first client connects with a user belonging to the socket's admin group, the user's OpenSSL certificate will be accepted by the daemon and the socket will then be open for all users to connect. Any other user trying to connect to the cinaclaw service will need to authenticate with the service using the previously set [`local.passphrase`](/reference/settings/local-passphrase).
````

````{tab-item} Windows
The Windows host uses a TCP socket listening on port 50051 for client connections. This socket is open for all to use since there is no concept of file ownership for TCP sockets. This is not very secure in that any cinaclaw user can connect to the service and issue any commands.

To close this gap, the user will now need to be authenticated with the cinaclaw service. To ease the burden of having to authenticate, the user who installs the updated version of cinaclaw will automatically have their clients authenticated with the service. Any other users connecting to the service will have to use authenticate using the previously set [`local.passphrase`](/reference/settings/local-passphrase).
````

`````
