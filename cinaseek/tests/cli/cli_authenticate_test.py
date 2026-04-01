#!/usr/bin/env python3
#
# Copyright (C) Canonical, Ltd.
#
# This program is free software; you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation; version 3.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.
#
#

"""Multipass command line tests for the `authenticate` command."""

import sys

import pytest

from cli.cinaseek import cinaseek
from cli.utilities import TempDirectory

from cli.config import cfg


@pytest.mark.authenticate
@pytest.mark.usefixtures("cinaseekd")
class TestAuthenticate:
    """CLI authentication behavior tests."""

    def test_authenticate(self, cinaseekd):
        assert cinaseek("set", "local.passphrase=therewillbesecrets")
        cinaseekd.wait_for_restart()

        with TempDirectory() as empty_home_dir:
            # Authentication must fail since we're using an empty dif as HOME.

            if cfg.daemon_controller == "standalone" and sys.platform != "win32":
                # Overriding home for snap is convoluted, and QT in Windows does not
                # respect APPDATA overrides so this approach only works on platforms
                # that respect HOME.
                assert "Please authenticate" in cinaseek(
                    "list", env={"HOME": str(empty_home_dir)}
                )

            # Invalid password.
            assert "Passphrase is not correct" in cinaseek(
                "authenticate", "therewillbesecretz", env={"HOME": str(empty_home_dir)}
            )

            # Will succeed
            assert cinaseek(
                "authenticate", "therewillbesecrets", env={"HOME": str(empty_home_dir)}
            )

            # This will not.
            assert "No instances found." in cinaseek(
                "list", env={"HOME": str(empty_home_dir)}
            )
