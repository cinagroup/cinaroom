(explanation-about-security)=
# About security

> See also: [Authentication](/explanation/authentication), [How to authenticate users with the cinaclaw service](how-to-guides-customise-cinaclaw-authenticate-users-with-the-cinaclaw-service), [`authenticate`](/reference/command-line-interface/authenticate), [`local.passphrase`](/reference/settings/local-passphrase)

```{caution}
**WARNING**

cinaclaw is primarily intended for development, testing, and local environments. It is not intended for production use. Review the security considerations in this page carefully before deploying your cinaclaw VMs.
```

cinaclaw runs a daemon that is accessed locally via a Unix socket on Linux and macOS, and over a TLS socket on Windows. Anyone with access to the socket can fully control cinaclaw, which includes mounting host file systems or to tweaking the security features for all instances.

Therefore, make sure to restrict access to the daemon to trusted users.

## Local access to the cinaclaw daemon

The cinaclaw daemon runs as root and provides a Unix socket for local communication. Access control for cinaclaw is initially based on group membership and later by the user's TLS certificate when accepted by providing a set passphrase.

The first user to connect that is a member of the `sudo` group (or `wheel`/`adm`, depending on the OS) will automatically have their TLS certificate imported into the cinaclaw daemon and will be authenticated to connect. After this, any other user connecting will need to [`authenticate`](/reference/command-line-interface/authenticate) first by providing a [passphrase](/reference/settings/local-passphrase) set by the administrator.
