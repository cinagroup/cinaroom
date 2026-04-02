#define MyAppName "CinaSeek"
#define MyAppVersion "1.0.0"
#define MyAppPublisher "CinaGroup"
#define MyAppURL "https://cinaseek.ai"

[Setup]
AppId={{CinaSeek-2026-001}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppPublisher={#MyAppPublisher}
AppPublisherURL={#MyAppURL}
DefaultDirName={autopf}\CinaSeek
DefaultGroupName=CinaSeek
OutputBaseFilename=cinaseek-setup-{#MyAppVersion}
Compression=lzma2/ultra64
SolidCompression=yes

[Files]
Source: "cinaseek-client-windows-amd64.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "install-openclaw.sh"; DestDir: "{app}\scripts"; Flags: ignoreversion

[Icons]
Name: "{group}\CinaSeek"; Filename: "{app}\cinaseek-client.exe"
Name: "{autodesktop}\CinaSeek"; Filename: "{app}\cinaseek-client.exe"

[Run]
Filename: "{app}\cinaseek-client.exe"; Description: "启动 CinaSeek"; Flags: nowait postinstall skipifsilent
