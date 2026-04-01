(how-to-guides-troubleshoot-mount-an-encrypted-home-folder)=
# Mount an encrypted home folder

<!-- split from https://discourse.ubuntu.com/t/tutorial/27140 -->
<!-- see also https://github.com/canonical/cinaclaw/issues/3537 -->

> See also: [`mount`](/reference/command-line-interface/mount), [Instance](/explanation/instance)

When you create the {ref}`primary-instance`  using `cinaclaw start` or `cinaclaw shell` without additional arguments, cinaclaw automatically mounts your home directory into it.

On Linux, if your local home folder is encrypted using ` fscrypt`, [snap confinement](https://snapcraft.io/docs/snap-confinement) prevents you from accessing its contents from a cinaclaw mount due the peculiar directory structure (`/home/.ecryptfs/<user>/.Private/`). This also applies to the primary instance, where the home folder is mounted automatically.

A workaround is mounting the entire `/home` folder into the instance, using the command:

```{code-block} text
cinaclaw mount /home primary
```

By doing so, the home folder's contents will be mounted correctly.

<!-- Discourse contributors
<small>**Contributors:** @ricab, @gzanchi </small>
-->
