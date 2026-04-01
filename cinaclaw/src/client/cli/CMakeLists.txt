# Copyright (C) Canonical, Ltd.
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License version 3 as
# published by the Free Software Foundation.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.

add_library(client STATIC
  argparser.cpp
  client.cpp)

target_link_libraries(client
  client_platform
  commands
  fmt::fmt-header-only
  formatter
  rpc
  scope_guard
  shared
  Qt6::Core)

add_executable(cinaclaw
  main.cpp)

target_link_libraries(cinaclaw
  client)

install(TARGETS cinaclaw
  DESTINATION bin
  COMPONENT cinaclaw)

if(APPLE)
  set_target_properties( cinaclaw PROPERTIES INSTALL_RPATH "@executable_path/../lib" )

  # Copy any framework depenencies into the install directory for packaging
  install(CODE "
    include(${CMAKE_SOURCE_DIR}/packaging/FrameworkUtilities.cmake)
    install_frameworks(\"${CMAKE_RUNTIME_OUTPUT_DIRECTORY}/cinaclaw\" \"\$ENV{DESTDIR}\${CMAKE_INSTALL_PREFIX}/lib\")
    " COMPONENT cinaclaw)
endif(APPLE)

if(MSVC)
  install(CODE "
    include(BundleUtilities)
    fixup_bundle(\"\${CMAKE_INSTALL_PREFIX}/bin/cinaclaw.exe\"  \"\"  \"${CMAKE_RUNTIME_OUTPUT_DIRECTORY}\")
    " COMPONENT cinaclaw)
endif()

add_subdirectory(cmd)
add_subdirectory(formatter)
