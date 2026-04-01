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

#include <cinaclaw/disabled_copy_move.h>
#include <cinaclaw/rpc/cinaclaw.grpc.pb.h>

#include <cinaclaw/cli/alias_dict.h>
#include <cinaclaw/cli/client_platform.h>

#include <string>

namespace cinaclaw
{
constexpr auto default_id_str = "default";

class Formatter : private DisabledCopyMove
{
public:
    virtual ~Formatter() = default;
    virtual std::string format(const InfoReply& reply) const = 0;
    virtual std::string format(const ListReply& reply) const = 0;
    virtual std::string format(const NetworksReply& reply) const = 0;
    virtual std::string format(const FindReply& reply) const = 0;
    virtual std::string format(const VersionReply& reply,
                               const std::string& client_version) const = 0;
    virtual std::string format(const AliasDict& aliases) const = 0;

protected:
    Formatter() = default;
};
} // namespace cinaclaw
