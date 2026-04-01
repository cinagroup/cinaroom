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

#include <cinaseek/format.h>
#include <cinaseek/logging/standard_logger.h>

#include <iostream>

namespace cinaseek::logging
{

StandardLogger::StandardLogger(Level level) : StandardLogger::StandardLogger(level, std::cerr)
{
}

StandardLogger::StandardLogger(Level level, std::ostream& target_ostream)
    : Logger{level}, target(target_ostream)
{
}

void StandardLogger::log(Level level, std::string_view category, std::string_view message) const
{
    if (level <= logging_level)
    {
        fmt::print(target, "[{}] [{}] [{}] {}\n", timestamp(), as_string(level), category, message);
    }
}
} // namespace cinaseek::logging
