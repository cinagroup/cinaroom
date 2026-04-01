(reference-command-line-interface-authenticate)=
# authenticate

> See also: [Authentication](/explanation/authentication), [How to authenticate users with the cinaseek service](how-to-guides-customise-cinaseek-authenticate-users-with-the-cinaseek-service), [`local.passhprase`](/reference/settings/local-passphrase)

The `authenticate` command is used to authenticate a user with the cinaseek service. Once authenticated, the user can issue commands such as `list`, `launch`, etc.

To help reduce the amount of typing for `authenticate`, one can also use `cinaseek auth` as an alias:

```{code-block} text
cinaseek auth foo
```

If no passphrase is specified in the `cinaseek authenticate` command line, you will be prompted to enter it.

---

The full `cinaseek help authenticate` output explains the available options:

```{code-block} text
Usage: cinaseek authenticate [options] [<passphrase>]
Authenticate with the cinaseek service.
A system administrator should provide you with a passphrase
to allow use of the cinaseek service.

Options:
  -h, --help     Displays help on commandline options
  -v, --verbose  Increase logging verbosity. Repeat the 'v' in the short option
                 for more detail. Maximum verbosity is obtained with 4 (or more)
                 v's, i.e. -vvvv.

Arguments:
  passphrase     Passphrase to register with the cinaseek service. If omitted,
                 a prompt will be displayed for entering the passphrase.
```
