(reference-command-line-interface-version)=
# version

The `cinaseek version` command without an argument will display the client and daemon versions of cinaseek; for example:

```{code-block} text
cinaseek  1.0.0
cinaseekd 1.0.0
```

If there is an update to cinaseek available, it will be printed out in addition to the standard output; for example:

```{code-block} text
cinaseek  1.0.0
cinaseekd 1.0.0

########################################################################################
cinaseek 1.0.1 release
Bugfix release to address a crash

Go here for more information: https://github.com/canonical/cinaseek/releases/tag/v1.0.1
########################################################################################
```

---

The full `cinaseek help version` output explains the available options:

```{code-block} text
Usage: cinaseek version [options]
Display version information about the cinaseek command
and daemon.

Options:
  -h, --help         Displays help on commandline options
  -v, --verbose      Increase logging verbosity. Repeat the 'v' in the short
                     option for more detail. Maximum verbosity is obtained with
                     4 (or more) v's, i.e. -vvvv.
  --format <format>  Output version information in the requested format.
                     Valid formats are: table (default), json, csv and yaml
```
