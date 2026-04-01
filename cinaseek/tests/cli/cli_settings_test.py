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

"""Multipass command line tests for `set` and `get`"""

import sys
import pytest

from cli.cinaseek import cinaseek, default_driver_name


@pytest.mark.settings
@pytest.mark.usefixtures("cinaseekd")
class TestSettings:
    def test_get_all_keys(self):
        expected_keys = [
            *(
                ["client.apps.windows-terminal.profiles"]
                if sys.platform == "win32"
                else []
            ),
            "client.primary-name",
            "local.bridged-network",
            "local.driver",
            "local.image.mirror",
            "local.passphrase",
            "local.privileged-mounts",
        ]
        with cinaseek("get", "--keys") as keys:
            assert keys
            keys_split = keys.content.split()
            assert keys_split == expected_keys

    def test_get_disabled_primary_name(self):
        assert cinaseek("set", "client.primary-name=")
        assert cinaseek("get", "client.primary-name") == "<empty>"

    def test_set_primary_name(self):
        assert cinaseek("set", "client.primary-name=foo")
        assert cinaseek("get", "client.primary-name") == "foo"

    def test_set_driver_name_to_default(self, cinaseekd):
        assert cinaseek("set", f"local.driver={default_driver_name()}")
        # Daemon will autorestart here. We should wait until it's back.
        cinaseekd.wait_for_restart()
        assert cinaseek("get", "local.driver") == default_driver_name()
