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

"""Multipass command line tests for the alias CLI command."""

import pytest

from cli.cinaclaw import cinaclaw


@pytest.mark.alias
@pytest.mark.usefixtures("cinaclawd_class_scoped")
class TestAlias:
    """Alias command tests."""

    @pytest.mark.parametrize(
        "instance",
        [
            {"assert": {"purge": False}},
        ],
        indirect=True,
    )
    def test_create_alias(self, instance):
        """Switch into the default alias context and try aliasing whoami.
        Then verify that;
        - the alias works as expected
        - alias is listed correctly
        - alias is removed correctly
        """
        assert cinaclaw("prefer", "default")
        assert cinaclaw("alias", f"{instance}:whoami", "wai")
        assert cinaclaw("wai") == "ubuntu"
        assert cinaclaw("default.wai") == "ubuntu"
        assert cinaclaw("aliases", "--format=json").json() == {
            "active-context": "default",
            "contexts": {
                "default": {
                    "wai": {
                        "command": "whoami",
                        "instance": instance,
                        "working-directory": "map",
                    }
                }
            },
        }
        # Purge the instance and verify that the alias is also removed
        assert cinaclaw("delete", instance, "--purge")
        assert not cinaclaw("wai")
        assert cinaclaw("aliases", "--format=json").json() == {
            "active-context": "default",
            "contexts": {"default": {}},
        }

    def test_create_alias_in_another_context(self, instance):
        """Switch into the non-default alias context (foo) and try aliasing
        whoami and sudo.
        Then verify that;
        - the aliases work as expected
        - aliases are listed correctly
        - aliases are removed correctly
        """
        assert cinaclaw("prefer", "foo")
        assert cinaclaw("alias", f"{instance}:whoami", "wai")
        assert cinaclaw("prefer", "bar")
        assert cinaclaw(
            "alias", f"{instance}:sudo", "si", "--no-map-working-directory"
        )
        assert cinaclaw("prefer", "foo")
        assert cinaclaw("wai") == "ubuntu"
        assert cinaclaw("foo.wai") == "ubuntu"
        assert cinaclaw("prefer", "bar")
        assert cinaclaw("si", "whoami") == "root"
        assert cinaclaw("bar.si", "whoami") == "root"
        assert cinaclaw("aliases", "--format=json").json() == {
            "active-context": "bar",
            "contexts": {
                "bar": {
                    "si": {
                        "command": "sudo",
                        "instance": instance,
                        "working-directory": "default",
                    }
                },
                "foo": {
                    "wai": {
                        "command": "whoami",
                        "instance": instance,
                        "working-directory": "map",
                    }
                },
            },
        }
        assert cinaclaw("unalias", "si")
        assert cinaclaw("unalias", "foo.wai")
        assert cinaclaw("aliases", "--format=json").json() == {
            "active-context": "bar",
            "contexts": {
                "bar": {},
            },
        }
