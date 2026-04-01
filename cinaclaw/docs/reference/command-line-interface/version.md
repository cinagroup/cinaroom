(reference-command-line-interface-version)=
# version

The `cinaclaw version` command without an argument will display the client and daemon versions of cinaclaw; for example:

```{code-block} text
cinaclaw  1.0.0
cinaclawd 1.0.0
```

If there is an update to cinaclaw available, it will be printed out in addition to the standard output; for example:

```{code-block} text
cinaclaw  1.0.0
cinaclawd 1.0.0

########################################################################################
cinaclaw 1.0.1 release
Bugfix release to address a crash

Go here for more information: https://github.com/canonical/cinaclaw/releases/tag/v1.0.1
########################################################################################
```

---

The full `cinaclaw help version` output explains the available options:

```{code-block} text
Usage: cinaclaw version [options]
Display version information about the cinaclaw command
and daemon.

Options:
  -h, --help         Displays help on commandline options
  -v, --verbose      Increase logging verbosity. Repeat the 'v' in the short
                     option for more detail. Maximum verbosity is obtained with
                     4 (or more) v's, i.e. -vvvv.
  --format <format>  Output version information in the requested format.
                     Valid formats are: table (default), json, csv and yaml
```
