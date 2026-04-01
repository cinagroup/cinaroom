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

#include "process_factory.h"
#include <cinaseek/logging/log.h>
#include <cinaseek/process/basic_process.h>
#include <cinaseek/process/process_spec.h>
#include <cinaseek/process/simple_process_spec.h>
#include <cinaseek/utils.h>

#include <unistd.h>

namespace mp = cinaseek;
namespace mpl = cinaseek::logging;

mp::ProcessFactory::ProcessFactory(const Singleton<ProcessFactory>::PrivatePass& pass)
    : Singleton<ProcessFactory>::Singleton{pass}
{
    [[maybe_unused]] static bool run_once = []() {
        ::setpgid(
            0,
            0); // create own process group. On MacOS, children of the parent are reaped if it dies.
        return true;
    }();
}

// This is the default ProcessFactory that creates a Process with no security mechanisms enabled
std::unique_ptr<mp::Process> mp::ProcessFactory::create_process(
    std::unique_ptr<mp::ProcessSpec>&& process_spec) const
{
    return std::make_unique<BasicProcess>(std::move(process_spec));
}

std::unique_ptr<mp::Process> mp::ProcessFactory::create_process(const QString& command,
                                                                const QStringList& arguments) const
{
    return create_process(simple_process_spec(command, arguments));
}
