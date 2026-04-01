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

#include <cinaclaw/image_host/vm_image_host.h>
#include <cinaclaw/vm_image.h>

namespace cinaclaw
{
namespace test
{
struct StubVMImageHost final : public cinaclaw::VMImageHost
{
    std::optional<cinaclaw::VMImageInfo> info_for(const cinaclaw::Query& query) override
    {
        return std::optional<cinaclaw::VMImageInfo>{
            VMImageInfo{{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, -1, {}}};
    };

    std::vector<std::pair<std::string, cinaclaw::VMImageInfo>> all_info_for(
        const cinaclaw::Query& query) override
    {
        return {};
    };

    cinaclaw::VMImageInfo info_for_full_hash(const std::string& full_hash) override
    {
        return {{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, -1, {}};
    };

    std::vector<cinaclaw::VMImageInfo> all_images_for(const std::string& remote_name,
                                                       const bool allow_unsupported) override
    {
        return {};
    };

    void for_each_entry_do(const Action&) override
    {
    }

    std::vector<std::string> supported_remotes() override
    {
        return {};
    }

    void update_manifests(const bool force_update) override
    {
    }
};
} // namespace test
} // namespace cinaclaw
