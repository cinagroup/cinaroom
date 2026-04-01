(reference-command-line-interface-purge)=
# purge

> See also: [`delete`](/reference/command-line-interface/delete), [`recover`](/reference/command-line-interface/recover)

The `cinaseek purge` command will permanently remove all instances deleted with the `cinaseek delete` command. This will destroy all the traces of the instance, and cannot be undone.

---

The full `cinaseek help purge` output explains the available options:

```{code-block} text
Usage: cinaseek purge [options]
Purge all deleted instances permanently, including all their data.

Options:
  -h, --help     Displays help on commandline options
  -v, --verbose  Increase logging verbosity. Repeat the 'v' in the short option
                 for more detail. Maximum verbosity is obtained with 4 (or more)
                 v's, i.e. -vvvv.
```
