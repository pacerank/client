# define installer name
OutFile "install_x64-.exe"

# set desktop as install directory
InstallDir "$PROGRAMFILES64\PaceRank"

# default section start
Section

# create a shortcut named "new shortcut" in the start menu programs directory
# presently, the new shortcut doesn't call anything (the second field is blank)
CreateShortcut "$SMPROGRAMS\PaceRank.lnk" "$INSTDIR\pacerank.exe"

# create desktop shortcut
CreateShortCut "$DESKTOP\PaceRank.lnk" "$INSTDIR\pacerank.exe"

# create autostart shortcut
CreateShortCut "$APPDATA\Microsoft\Windows\Start Menu\Programs\Startup\PaceRank.lnk" "$INSTDIR\pacerank.exe"

CreateDirectory "$PROFILE\.pacerank"

# define output path
SetOutPath $INSTDIR

# specify file to go in output path
File pacerank.exe
File sciter.dll

# define uninstaller name
WriteUninstaller $INSTDIR\uninstall.exe

# default section end
ExecShell "" "$INSTDIR\pacerank.exe"
SectionEnd

# -----

# create a section to define what the uninstaller does.
# the section will always be named "Uninstall"
Section "Uninstall"

# Always delete uninstaller first
Delete $INSTDIR\uninstall.exe

# now delete installed file
Delete $INSTDIR\pacerank.exe
Delete $INSTDIR\sciter.dll

# delete shortcut
Delete $SMPROGRAMS\PaceRank.lnk
Delete $DESKTOP\PaceRank.lnk
Delete "$APPDATA\Roaming\Microsoft\Windows\Start Menu\Programs\Startup\PaceRank.lnk"

# Delete the directory
RMDir $INSTDIR

SectionEnd
