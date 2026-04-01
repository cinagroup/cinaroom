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

#include "animated_spinner.h"

#include <cinaseek/cli/alias_dict.h>
#include <cinaseek/cli/command.h>
#include <cinaseek/timer.h>

#include <QString>

#include <memory>
#include <string>
#include <utility>

namespace cinaseek
{
namespace cmd
{
class Launch final : public Command
{
public:
    using Command::Command;

    Launch(Rpc::StubInterface& stub, Terminal* term, AliasDict& dict)
        : Command(stub, term), aliases(dict)
    {
    }

    ReturnCodeVariant run(ArgParser* parser) override;
    std::string name() const override;
    QString short_help() const override;
    QString description() const override;

private:
    ParseCode parse_args(ArgParser* parser);
    ReturnCodeVariant request_launch(const ArgParser* parser);
    ReturnCodeVariant mount(const ArgParser* parser,
                            const QString& mount_source,
                            const QString& mount_target);
    bool ask_bridge_permission(cinaseek::LaunchReply& reply);

    LaunchRequest request;
    QString petenv_name;
    std::unique_ptr<cinaseek::AnimatedSpinner> spinner;
    std::unique_ptr<cinaseek::utils::Timer> timer;

    std::vector<std::pair<QString, QString>> mount_routes;
    QString instance_name;

    AliasDict aliases;
};
} // namespace cmd
} // namespace cinaseek
