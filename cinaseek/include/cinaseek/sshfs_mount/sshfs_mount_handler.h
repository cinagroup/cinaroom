/*
 * Copyright (C) Canonical, Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; version 3.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

#pragma once

#include <cinaseek/mount_handler.h>
#include <cinaseek/process/process.h>
#include <cinaseek/qt_delete_later_unique_ptr.h>
#include <cinaseek/sshfs_server_config.h>

namespace cinaseek
{
class SSHFSMountHandler : public MountHandler
{
public:
    SSHFSMountHandler(VirtualMachine* vm,
                      const SSHKeyProvider* ssh_key_provider,
                      const std::string& target,
                      VMMount mount_spec);
    ~SSHFSMountHandler() override;

    void activate_impl(ServerVariant server, std::chrono::milliseconds timeout) override;
    void deactivate_impl(bool force) override;

private:
    qt_delete_later_unique_ptr<Process> process;
    SSHFSServerConfig config;
};
} // namespace cinaseek
