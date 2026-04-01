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

#include <cinaclaw/cert_provider.h>
#include <cinaclaw/cert_store.h>
#include <cinaclaw/days.h>
#include <cinaclaw/image_host/vm_image_host.h>
#include <cinaclaw/logging/logger.h>
#include <cinaclaw/logging/multiplexing_logger.h>
#include <cinaclaw/name_generator.h>
#include <cinaclaw/path.h>
#include <cinaclaw/rpc/cinaclaw.grpc.pb.h>
#include <cinaclaw/ssh/ssh_key_provider.h>
#include <cinaclaw/update_prompt.h>
#include <cinaclaw/url_downloader.h>
#include <cinaclaw/virtual_machine_factory.h>
#include <cinaclaw/vm_image_vault.h>

#include <QNetworkProxy>

#include <memory>
#include <vector>

namespace cinaclaw
{
struct DaemonConfig
{
    ~DaemonConfig();
    const std::unique_ptr<URLDownloader> url_downloader;
    const std::unique_ptr<VirtualMachineFactory> factory;
    const std::vector<std::unique_ptr<VMImageHost>> image_hosts;
    const std::unique_ptr<VMImageVault> vault;
    const std::unique_ptr<NameGenerator> name_generator;
    const std::unique_ptr<SSHKeyProvider> ssh_key_provider;
    const std::unique_ptr<CertProvider> cert_provider;
    const std::unique_ptr<CertStore> client_cert_store;
    const std::unique_ptr<UpdatePrompt> update_prompt;
    const std::shared_ptr<logging::MultiplexingLogger> logger;
    const std::unique_ptr<QNetworkProxy> network_proxy;
    const cinaclaw::Path cache_directory;
    const cinaclaw::Path data_directory;
    const std::string server_address;
    const std::string ssh_username;
    const std::chrono::hours image_refresh_timer;
};

struct DaemonConfigBuilder
{
    std::unique_ptr<URLDownloader> url_downloader;
    std::unique_ptr<VirtualMachineFactory> factory;
    std::vector<std::unique_ptr<VMImageHost>> image_hosts;
    std::unique_ptr<VMImageVault> vault;
    std::unique_ptr<NameGenerator> name_generator;
    std::unique_ptr<SSHKeyProvider> ssh_key_provider;
    std::unique_ptr<CertProvider> cert_provider;
    std::unique_ptr<CertStore> client_cert_store;
    std::unique_ptr<UpdatePrompt> update_prompt;
    std::unique_ptr<logging::Logger> logger;
    std::unique_ptr<QNetworkProxy> network_proxy;
    cinaclaw::Path cache_directory;
    cinaclaw::Path data_directory;
    std::string server_address;
    std::string ssh_username;
    cinaclaw::days days_to_expire{14};
    std::chrono::hours image_refresh_timer{6};
    cinaclaw::logging::Level verbosity_level{cinaclaw::logging::Level::info};

    std::unique_ptr<const DaemonConfig> build();
};
} // namespace cinaclaw
