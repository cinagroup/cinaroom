(reference-command-line-interface-authenticate)=
# authenticate

> See also: [Authentication](/explanation/authentication), [How to authenticate users with the cinaclaw service](how-to-guides-customise-cinaclaw-authenticate-users-with-the-cinaclaw-service), [`local.passhprase`](/reference/settings/local-passphrase)

The `authenticate` command is used to authenticate a user with the cinaclaw service. Once authenticated, the user can issue commands such as `list`, `launch`, etc.

To help reduce the amount of typing for `authenticate`, one can also use `cinaclaw auth` as an alias:

```{code-block} text
cinaclaw auth foo
```

If no passphrase is specified in the `cinaclaw authenticate` command line, you will be prompted to enter it.

---

The full `cinaclaw help authenticate` output explains the available options:

```{code-block} text
Usage: cinaclaw authenticate [options] [<passphrase>]
Authenticate with the cinaclaw service.
A system administrator should provide you with a passphrase
to allow use of the cinaclaw service.

Options:
  -h, --help     Displays help on commandline options
  -v, --verbose  Increase logging verbosity. Repeat the 'v' in the short option
                 for more detail. Maximum verbosity is obtained with 4 (or more)
                 v's, i.e. -vvvv.

Arguments:
  passphrase     Passphrase to register with the cinaclaw service. If omitted,
                 a prompt will be displayed for entering the passphrase.
```
