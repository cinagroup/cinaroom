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

#include <cinaseek/cert_provider.h>
#include <cinaseek/cli/alias_dict.h>
#include <cinaseek/cli/command.h>
#include <cinaseek/rpc/cinaseek.grpc.pb.h>
#include <cinaseek/terminal.h>

#include <memory>

namespace cinaseek
{
struct ClientConfig
{
    const std::string server_address;
    std::unique_ptr<CertProvider> cert_provider;
    Terminal* term;
};

class Client
{
public:
    explicit Client(ClientConfig& context);
    virtual ~Client() = default;
    ReturnCodeVariant run(const QStringList& arguments);

protected:
    template <typename T, typename... Ts>
    void add_command(Ts&&... params);
    void sort_commands();

private:
    std::unique_ptr<cinaseek::Rpc::Stub> stub;

    std::vector<cmd::Command::UPtr> commands;

    Terminal* term;
    AliasDict aliases;
};
} // namespace cinaseek

template <typename T, typename... Ts>
void cinaseek::Client::add_command(Ts&&... params)
{
    auto cmd = std::make_unique<T>(*stub, term, std::forward<Ts>(params)...);
    commands.push_back(std::move(cmd));
}
