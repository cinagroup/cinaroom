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

#include <cinaclaw/rpc/cinaclaw.grpc.pb.h>

using namespace testing;

namespace cinaclaw::test
{
template <class W, class R>
class MockClientReaderWriter : public grpc::ClientReaderWriterInterface<W, R>
{
public:
    MockClientReaderWriter()
    {
        EXPECT_CALL(*this, Read(_)).WillRepeatedly(Return(false));
        EXPECT_CALL(*this, Finish()).WillRepeatedly(Return(grpc::Status()));
    }

    MOCK_METHOD(grpc::Status, Finish, (), (override));
    MOCK_METHOD(bool, NextMessageSize, (uint32_t * sz), (override));
    MOCK_METHOD(bool, Read, (R * msg), (override));
    MOCK_METHOD(void, WaitForInitialMetadata, (), (override));
    MOCK_METHOD(bool, Write, (const W& msg, grpc::WriteOptions options), (override));
    MOCK_METHOD(bool, WritesDone, (), (override));
};

class MockRpcStub : public cinaclaw::Rpc::StubInterface
{
public:
    MockRpcStub() = default;

    MOCK_METHOD(grpc::Status,
                ping,
                (grpc::ClientContext * context,
                 const cinaclaw::PingRequest& request,
                 cinaclaw::PingReply* response),
                (override));

    // originally private
    // NB: we're sort of relying on gRPC implementation here... but it's only for tests and we can
    // update as needed
    MOCK_METHOD(
        (grpc::ClientReaderWriterInterface<cinaclaw::LaunchRequest, cinaclaw::LaunchReply>*),
        createRaw,
        (grpc::ClientContext * context),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaclaw::LaunchRequest, cinaclaw::LaunchReply>*),
        AsynccreateRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaclaw::LaunchRequest, cinaclaw::LaunchReply>*),
        PrepareAsynccreateRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq),
        (override));
    MOCK_METHOD(
        (grpc::ClientReaderWriterInterface<cinaclaw::LaunchRequest, cinaclaw::LaunchReply>*),
        launchRaw,
        (grpc::ClientContext * context),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaclaw::LaunchRequest, cinaclaw::LaunchReply>*),
        AsynclaunchRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaclaw::LaunchRequest, cinaclaw::LaunchReply>*),
        PrepareAsynclaunchRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq),
        (override));
    MOCK_METHOD(
        (grpc::ClientReaderWriterInterface<cinaclaw::PurgeRequest, cinaclaw::PurgeReply>*),
        purgeRaw,
        (grpc::ClientContext * context),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaclaw::PurgeRequest, cinaclaw::PurgeReply>*),
        AsyncpurgeRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaclaw::PurgeRequest, cinaclaw::PurgeReply>*),
        PrepareAsyncpurgeRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq),
        (override));
    MOCK_METHOD((grpc::ClientReaderWriterInterface<cinaclaw::FindRequest, cinaclaw::FindReply>*),
                findRaw,
                (grpc::ClientContext * context),
                (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaclaw::FindRequest, cinaclaw::FindReply>*),
        AsyncfindRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaclaw::FindRequest, cinaclaw::FindReply>*),
        PrepareAsyncfindRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq),
        (override));
    MOCK_METHOD((grpc::ClientReaderWriterInterface<cinaclaw::InfoRequest, cinaclaw::InfoReply>*),
                infoRaw,
                (grpc::ClientContext * context),
                (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaclaw::InfoRequest, cinaclaw::InfoReply>*),
        AsyncinfoRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaclaw::InfoRequest, cinaclaw::InfoReply>*),
        PrepareAsyncinfoRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq),
        (override));
    MOCK_METHOD((grpc::ClientReaderWriterInterface<cinaclaw::ListRequest, cinaclaw::ListReply>*),
                listRaw,
                (grpc::ClientContext * context),
                (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaclaw::ListRequest, cinaclaw::ListReply>*),
        AsynclistRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaclaw::ListRequest, cinaclaw::ListReply>*),
        PrepareAsynclistRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq),
        (override));
    MOCK_METHOD(
        (grpc::ClientReaderWriterInterface<cinaclaw::NetworksRequest, cinaclaw::NetworksReply>*),
        networksRaw,
        (grpc::ClientContext * context),
        (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaclaw::NetworksRequest,
                                                        cinaclaw::NetworksReply>*),
                AsyncnetworksRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
                (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaclaw::NetworksRequest,
                                                        cinaclaw::NetworksReply>*),
                PrepareAsyncnetworksRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq),
                (override));
    MOCK_METHOD(
        (grpc::ClientReaderWriterInterface<cinaclaw::MountRequest, cinaclaw::MountReply>*),
        mountRaw,
        (grpc::ClientContext * context),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaclaw::MountRequest, cinaclaw::MountReply>*),
        AsyncmountRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaclaw::MountRequest, cinaclaw::MountReply>*),
        PrepareAsyncmountRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq),
        (override));
    MOCK_METHOD(grpc::ClientAsyncResponseReaderInterface<cinaclaw::PingReply>*,
                AsyncpingRaw,
                (grpc::ClientContext * context,
                 const cinaclaw::PingRequest& request,
                 grpc::CompletionQueue* cq),
                (override));
    MOCK_METHOD(grpc::ClientAsyncResponseReaderInterface<cinaclaw::PingReply>*,
                PrepareAsyncpingRaw,
                (grpc::ClientContext * context,
                 const cinaclaw::PingRequest& request,
                 grpc::CompletionQueue* cq),
                (override));
    MOCK_METHOD(
        (grpc::ClientReaderWriterInterface<cinaclaw::RecoverRequest, cinaclaw::RecoverReply>*),
        recoverRaw,
        (grpc::ClientContext * context),
        (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaclaw::RecoverRequest,
                                                        cinaclaw::RecoverReply>*),
                AsyncrecoverRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
                (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaclaw::RecoverRequest,
                                                        cinaclaw::RecoverReply>*),
                PrepareAsyncrecoverRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq),
                (override));
    MOCK_METHOD(
        (grpc::ClientReaderWriterInterface<cinaclaw::SSHInfoRequest, cinaclaw::SSHInfoReply>*),
        ssh_infoRaw,
        (grpc::ClientContext * context),
        (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaclaw::SSHInfoRequest,
                                                        cinaclaw::SSHInfoReply>*),
                Asyncssh_infoRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
                (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaclaw::SSHInfoRequest,
                                                        cinaclaw::SSHInfoReply>*),
                PrepareAsyncssh_infoRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq),
                (override));
    MOCK_METHOD(
        (grpc::ClientReaderWriterInterface<cinaclaw::StartRequest, cinaclaw::StartReply>*),
        startRaw,
        (grpc::ClientContext * context),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaclaw::StartRequest, cinaclaw::StartReply>*),
        AsyncstartRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaclaw::StartRequest, cinaclaw::StartReply>*),
        PrepareAsyncstartRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq),
        (override));
    MOCK_METHOD((grpc::ClientReaderWriterInterface<cinaclaw::StopRequest, cinaclaw::StopReply>*),
                stopRaw,
                (grpc::ClientContext * context),
                (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaclaw::StopRequest, cinaclaw::StopReply>*),
        AsyncstopRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaclaw::StopRequest, cinaclaw::StopReply>*),
        PrepareAsyncstopRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq),
        (override));
    MOCK_METHOD(
        (grpc::ClientReaderWriterInterface<cinaclaw::SuspendRequest, cinaclaw::SuspendReply>*),
        suspendRaw,
        (grpc::ClientContext * context),
        (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaclaw::SuspendRequest,
                                                        cinaclaw::SuspendReply>*),
                AsyncsuspendRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
                (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaclaw::SuspendRequest,
                                                        cinaclaw::SuspendReply>*),
                PrepareAsyncsuspendRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq),
                (override));
    MOCK_METHOD(
        (grpc::ClientReaderWriterInterface<cinaclaw::RestartRequest, cinaclaw::RestartReply>*),
        restartRaw,
        (grpc::ClientContext * context),
        (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaclaw::RestartRequest,
                                                        cinaclaw::RestartReply>*),
                AsyncrestartRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
                (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaclaw::RestartRequest,
                                                        cinaclaw::RestartReply>*),
                PrepareAsyncrestartRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq),
                (override));
    MOCK_METHOD(
        (grpc::ClientReaderWriterInterface<cinaclaw::DeleteRequest, cinaclaw::DeleteReply>*),
        deletRaw,
        (grpc::ClientContext * context),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaclaw::DeleteRequest, cinaclaw::DeleteReply>*),
        AsyncdeletRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaclaw::DeleteRequest, cinaclaw::DeleteReply>*),
        PrepareAsyncdeletRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq),
        (override));
    MOCK_METHOD(
        (grpc::ClientReaderWriterInterface<cinaclaw::UmountRequest, cinaclaw::UmountReply>*),
        umountRaw,
        (grpc::ClientContext * context),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaclaw::UmountRequest, cinaclaw::UmountReply>*),
        AsyncumountRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaclaw::UmountRequest, cinaclaw::UmountReply>*),
        PrepareAsyncumountRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq),
        (override));
    MOCK_METHOD(
        (grpc::ClientReaderWriterInterface<cinaclaw::VersionRequest, cinaclaw::VersionReply>*),
        versionRaw,
        (grpc::ClientContext * context),
        (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaclaw::VersionRequest,
                                                        cinaclaw::VersionReply>*),
                AsyncversionRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
                (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaclaw::VersionRequest,
                                                        cinaclaw::VersionReply>*),
                PrepareAsyncversionRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq),
                (override));
    MOCK_METHOD((grpc::ClientReaderWriterInterface<cinaclaw::GetRequest, cinaclaw::GetReply>*),
                getRaw,
                (grpc::ClientContext * context),
                (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaclaw::GetRequest, cinaclaw::GetReply>*),
        AsyncgetRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaclaw::GetRequest, cinaclaw::GetReply>*),
        PrepareAsyncgetRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq),
        (override));
    MOCK_METHOD((grpc::ClientReaderWriterInterface<cinaclaw::SetRequest, cinaclaw::SetReply>*),
                setRaw,
                (grpc::ClientContext * context),
                (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaclaw::SetRequest, cinaclaw::SetReply>*),
        AsyncsetRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaclaw::SetRequest, cinaclaw::SetReply>*),
        PrepareAsyncsetRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq),
        (override));
    MOCK_METHOD((grpc::ClientReaderWriterInterface<cinaclaw::KeysRequest, cinaclaw::KeysReply>*),
                keysRaw,
                (grpc::ClientContext * context),
                (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaclaw::KeysRequest, cinaclaw::KeysReply>*),
        AsynckeysRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaclaw::KeysRequest, cinaclaw::KeysReply>*),
        PrepareAsynckeysRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq),
        (override));
    MOCK_METHOD((grpc::ClientReaderWriterInterface<cinaclaw::AuthenticateRequest,
                                                   cinaclaw::AuthenticateReply>*),
                authenticateRaw,
                (grpc::ClientContext * context),
                (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaclaw::AuthenticateRequest,
                                                        cinaclaw::AuthenticateReply>*),
                AsyncauthenticateRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
                (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaclaw::AuthenticateRequest,
                                                        cinaclaw::AuthenticateReply>*),
                PrepareAsyncauthenticateRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq),
                (override));
    MOCK_METHOD(
        (grpc::ClientReaderWriterInterface<cinaclaw::SnapshotRequest, cinaclaw::SnapshotReply>*),
        snapshotRaw,
        (grpc::ClientContext * context),
        (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaclaw::SnapshotRequest,
                                                        cinaclaw::SnapshotReply>*),
                AsyncsnapshotRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
                (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaclaw::SnapshotRequest,
                                                        cinaclaw::SnapshotReply>*),
                PrepareAsyncsnapshotRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq),
                (override));
    MOCK_METHOD(
        (grpc::ClientReaderWriterInterface<cinaclaw::RestoreRequest, cinaclaw::RestoreReply>*),
        restoreRaw,
        (grpc::ClientContext * context),
        (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaclaw::RestoreRequest,
                                                        cinaclaw::RestoreReply>*),
                AsyncrestoreRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
                (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaclaw::RestoreRequest,
                                                        cinaclaw::RestoreReply>*),
                PrepareAsyncrestoreRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq),
                (override));
    MOCK_METHOD(
        (grpc::ClientReaderWriterInterface<cinaclaw::CloneRequest, cinaclaw::CloneReply>*),
        cloneRaw,
        (grpc::ClientContext * context),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaclaw::CloneRequest, cinaclaw::CloneReply>*),
        AsynccloneRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaclaw::CloneRequest, cinaclaw::CloneReply>*),
        PrepareAsynccloneRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq),
        (override));
    MOCK_METHOD((grpc::ClientReaderWriterInterface<cinaclaw::DaemonInfoRequest,
                                                   cinaclaw::DaemonInfoReply>*),
                daemon_infoRaw,
                (grpc::ClientContext * context),
                (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaclaw::DaemonInfoRequest,
                                                        cinaclaw::DaemonInfoReply>*),
                Asyncdaemon_infoRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
                (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaclaw::DaemonInfoRequest,
                                                        cinaclaw::DaemonInfoReply>*),
                PrepareAsyncdaemon_infoRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq),
                (override));
    MOCK_METHOD((grpc::ClientReaderWriterInterface<cinaclaw::WaitReadyRequest,
                                                   cinaclaw::WaitReadyReply>*),
                wait_readyRaw,
                (grpc::ClientContext * context),
                (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaclaw::WaitReadyRequest,
                                                        cinaclaw::WaitReadyReply>*),
                Asyncwait_readyRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
                (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaclaw::WaitReadyRequest,
                                                        cinaclaw::WaitReadyReply>*),
                PrepareAsyncwait_readyRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq),
                (override));
};
} // namespace cinaclaw::test
