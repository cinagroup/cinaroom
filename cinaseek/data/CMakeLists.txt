if(MULTIPASS_ENABLE_FLUTTER_GUI)
  install(
    FILES cinaseek.gui.desktop
    DESTINATION ${CMAKE_INSTALL_DATADIR}/applications
  )
endif()

install(
  FILES cinaseek.gui.svg
  DESTINATION ${CMAKE_INSTALL_DATADIR}/icons/hicolor/scalable/apps
)
