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

#include <cinaseek/memory_size.h>
#include <cinaseek/network_interface.h>
#include <cinaseek/path.h>
#include <cinaseek/vm_image.h>

#include <yaml-cpp/yaml.h>

#include <string>
#include <vector>

#include <QMetaType>

namespace cinaseek
{
class VirtualMachineDescription
{
public:
    using MBytes = size_t;

    int num_cores;
    MemorySize mem_size;
    MemorySize disk_space;
    std::string vm_name;
    std::string default_mac_address;
    std::vector<NetworkInterface> extra_interfaces;
    std::string ssh_username;
    VMImage image;
    Path cloud_init_iso;
    YAML::Node meta_data_config;
    YAML::Node user_data_config;
    YAML::Node vendor_data_config;
    YAML::Node network_data_config;
};
} // namespace cinaseek

Q_DECLARE_METATYPE(cinaseek::VirtualMachineDescription)
