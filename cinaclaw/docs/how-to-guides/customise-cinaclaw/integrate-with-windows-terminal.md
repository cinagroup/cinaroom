(how-to-guides-customise-cinaclaw-integrate-with-windows-terminal)=
# How to integrate with Windows Terminal

If you are on Windows and you want to use [Windows Terminal](https://aka.ms/terminal), cinaclaw can integrate with it to offer you an automatic `primary` profile.

## cinaclaw profile

Currently, cinaclaw can add a profile to Windows Terminal for the {ref}`primary-instance`. When you open a Windows Terminal tab with this profile, you'll automatically find yourself in a primary instance shell. cinaclaw automatically starts or launches the primary instance if needed.

```{figure} /images/cinaclaw-windows-terminal-1.png
   :width: 680px
   :alt: Screenshot: primary shell
```

<!-- Original image on the Asset Manager
![Screenshot: primary shell|800x490, 85%](https://assets.ubuntu.com/v1/f875c1d3-cinaclaw-windows-terminal-1.png)
-->

## Install Windows Terminal

The first step is to [install Windows Terminal](https://github.com/microsoft/terminal#installing-and-running-windows-terminal). Once you have it along cinaclaw, you can enable the integration.

## Enable integration

Open a terminal (Windows Terminal or any other) and enable the integration with the following command:

```{code-block} text
cinaclaw set client.apps.windows-terminal.profiles=primary
```

For more information on this setting, see [`client.apps.windows-terminal.profiles`](reference-settings-client-apps-windows-terminal-profiles). Until you modify it, cinaclaw will try to add the profile if it finds it missing. To remove the profile see {ref}`integrate-with-windows-terminal-revert` below.

## Open a cinaclaw tab

You can now open a "cinaclaw" tab to get a shell in the primary instance. That can be achieved by clicking the new-tab drop-down and selecting the cinaclaw profile:

```{figure} /images/cinaclaw-windows-terminal-2.jpeg
   :width: 680px
   :alt: Screenshot: drop-down menu
```

<!-- Original image on the Asset Manager
![Screenshot: drop-down menu|800x490, 85%](https://assets.ubuntu.com/v1/d14d32d6-cinaclaw-windows-terminal-2.jpeg)
-->

That's it!

(integrate-with-windows-terminal-revert)=
## Revert

If you want to disable the profile again, you can do so with:

```{code-block} text
cinaclaw set client.apps.windows-terminal.profiles=none
```

cinaclaw will then remove the profile if it exists.
