if(MULTIPASS_ENABLE_FLUTTER_GUI)
  install(
    FILES cinaclaw.gui.desktop
    DESTINATION ${CMAKE_INSTALL_DATADIR}/applications
  )
endif()

install(
  FILES cinaclaw.gui.svg
  DESTINATION ${CMAKE_INSTALL_DATADIR}/icons/hicolor/scalable/apps
)
