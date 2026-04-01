(reference-command-line-interface-recover)=
# recover

> See also: [`delete`](/reference/command-line-interface/delete), [`purge`](/reference/command-line-interface/purge)

The `cinaseek recover` command will revive an instance that was previously removed with `cinaseek delete`. For this to be possible, the instance cannot have been purged with `cinaseek purge` nor with `cinaseek delete --purge`.

Use the `--all` option to recover all deleted instances at once:

```{code-block} text
cinaseek recover --all
```

---

The full `cinaseek help restart` output explains the available options:

```{code-block} text
Usage: cinaseek recover [options] <name> [<name> ...]
Recover deleted instances so they can be used again.

Options:
  -h, --help     Display this help on commandline options
  -v, --verbose  Increase logging verbosity. Repeat the 'v' in the short option
                 for more detail. Maximum verbosity is obtained with 4 (or more)
                 v's, i.e. -vvvv.
  --all          Recover all deleted instances

Arguments:
  name           Names of instances to recover
```
