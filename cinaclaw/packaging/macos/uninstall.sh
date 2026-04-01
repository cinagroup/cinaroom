#!/bin/sh
set -eu

if [ "$EUID" -ne 0 ]; then
    echo "This script needs to run as root"
    exit 1
fi

while true; do
    read -p "Are you sure you want to remove Multipass from your system? [Y/N] " yn
    case $yn in
        [Yy]* ) break;;
        [Nn]* ) echo "Aborted"; exit;;
        * ) echo "Please answer yes or no.";;
    esac
done

DELETE_VMS=0

while true; do
    read -p "Do you want to delete all your Multipass VMs and daemon data too? [Y/N] " yn
    case $yn in
        [Yy]* ) DELETE_VMS=1; break;;
        [Nn]* ) DELETE_VMS=0; break;;
        * ) echo "Please answer yes or no.";;
    esac
done

if [ $DELETE_VMS -eq 1 ]; then
    echo "Removing VMs:"
    sudo -u "$(logname)" cinaclaw delete -vv --purge --all || echo "Failed to delete cinaclaw VMs from underlying driver" >&2

fi

LAUNCH_AGENT_DEST="/Library/LaunchDaemons/com.canonical.cinaclawd.plist"

echo .
echo "Removing the Multipass daemon launch agent:"
launchctl unload -w "$LAUNCH_AGENT_DEST"

if [ $DELETE_VMS -eq 1 ]; then
    echo "Removing daemon data:"
    rm -rfv "/var/root/Library/Application Support/cinaclawd"
    rm -rfv "/var/root/Library/Application Support/cinaclaw-client-certificate"
    rm -rfv "/var/root/Library/Preferences/cinaclawd"
    rm -fv "/Library/Keychains/cinaclaw_root_cert.pem"
fi

echo .
echo "Removing Multipass:"
rm -fv "$LAUNCH_AGENT_DEST"

rm -fv /usr/local/bin/cinaclaw
rm -rfv /Applications/Multipass.app

rm -rfv "/Library/Application Support/com.canonical.cinaclaw"
rm -rfv "/var/root/Library/Caches/cinaclawd"

# GUI Autostart
rm -fv "$HOME/Library/LaunchAgents/com.canonical.cinaclaw.gui.autostart.plist"

# User-specific client certificates and GUI data
rm -rfv "$HOME/Library/Application Support/cinaclaw-client-certificate"
rm -rfv "$HOME/Library/Application Support/com.canonical.cinaclawGui"
rm -rfv "$HOME/Library/Preferences/cinaclaw"

# Bash completions
rm -rfv "/usr/local/etc/bash_completion.d/cinaclaw"
rm -rf "/opt/local/share/bash-completion/completions/cinaclaw"

# Log files
rm -rfv "/Library/Logs/Multipass"

echo .
echo "Removing package installation receipts"
rm -fv "/private/var/db/receipts/com.canonical.cinaclaw.cinaclawd.bom"
rm -fv "/private/var/db/receipts/com.canonical.cinaclaw.cinaclawd.plist"
rm -fv "/private/var/db/receipts/com.canonical.cinaclaw.cinaclaw.bom"
rm -fv "/private/var/db/receipts/com.canonical.cinaclaw.cinaclaw.plist"

echo .
echo "Uninstall complete"

if [ $DELETE_VMS -eq 0 ]; then
    echo "Your Multipass VMs were preserved in /var/root/Library/Application Support/cinaclawd"
fi
