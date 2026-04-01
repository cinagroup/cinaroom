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

#include <memory>

#include <cinaseek/process/process_spec.h>
#include <cinaseek/singleton.h>

#define MP_PROCFACTORY cinaseek::ProcessFactory::instance()

namespace cinaseek
{
class Process;

class ProcessFactory : public Singleton<ProcessFactory>
{
public:
    ProcessFactory(const Singleton<ProcessFactory>::PrivatePass&);

    virtual std::unique_ptr<Process> create_process(
        std::unique_ptr<ProcessSpec>&& process_spec) const;
    std::unique_ptr<Process> create_process(const QString& command,
                                            const QStringList& = QStringList()) const;
};

} // namespace cinaseek
