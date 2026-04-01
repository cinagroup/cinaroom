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

#include "common.h"

#include <cinaclaw/logging/log.h>
#include <cinaclaw/logging/logger.h>
#include <cinaclaw/private_pass_provider.h>

namespace cinaclaw
{
namespace test
{
class MockLogger : public cinaclaw::logging::Logger, public PrivatePassProvider<MockLogger>
{
public:
    MockLogger(const PrivatePass&, const cinaclaw::logging::Level logging_level);

    MOCK_METHOD(void,
                log,
                (cinaclaw::logging::Level level,
                 std::string_view category,
                 std::string_view message),
                (const, override));

    class Scope
    {
    public:
        ~Scope();
        std::shared_ptr<testing::NiceMock<MockLogger>> mock_logger;

    private:
        Scope(const cinaclaw::logging::Level logging_level);
        friend class MockLogger;
    };

    // only one at a time, please
    [[nodiscard]] static Scope inject(
        const cinaclaw::logging::Level logging_level = cinaclaw::logging::Level::error);

    void expect_log(cinaclaw::logging::Level lvl,
                    const std::string& substr,
                    const testing::Cardinality& times = testing::Exactly(1));

    // Reject logs with severity `lvl` or higher (lower integer), accept the rest
    // By default, all logs are rejected. Pass error level to accept everything but errors (expect
    // those explicitly)
    void screen_logs(cinaclaw::logging::Level lvl = cinaclaw::logging::Level::trace);
};
} // namespace test
} // namespace cinaclaw
