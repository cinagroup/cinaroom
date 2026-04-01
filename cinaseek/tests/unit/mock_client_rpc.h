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

#include <cinaseek/rpc/cinaseek.grpc.pb.h>

using namespace testing;

namespace cinaseek::test
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

class MockRpcStub : public cinaseek::Rpc::StubInterface
{
public:
    MockRpcStub() = default;

    MOCK_METHOD(grpc::Status,
                ping,
                (grpc::ClientContext * context,
                 const cinaseek::PingRequest& request,
                 cinaseek::PingReply* response),
                (override));

    // originally private
    // NB: we're sort of relying on gRPC implementation here... but it's only for tests and we can
    // update as needed
    MOCK_METHOD(
        (grpc::ClientReaderWriterInterface<cinaseek::LaunchRequest, cinaseek::LaunchReply>*),
        createRaw,
        (grpc::ClientContext * context),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaseek::LaunchRequest, cinaseek::LaunchReply>*),
        AsynccreateRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaseek::LaunchRequest, cinaseek::LaunchReply>*),
        PrepareAsynccreateRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq),
        (override));
    MOCK_METHOD(
        (grpc::ClientReaderWriterInterface<cinaseek::LaunchRequest, cinaseek::LaunchReply>*),
        launchRaw,
        (grpc::ClientContext * context),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaseek::LaunchRequest, cinaseek::LaunchReply>*),
        AsynclaunchRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaseek::LaunchRequest, cinaseek::LaunchReply>*),
        PrepareAsynclaunchRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq),
        (override));
    MOCK_METHOD(
        (grpc::ClientReaderWriterInterface<cinaseek::PurgeRequest, cinaseek::PurgeReply>*),
        purgeRaw,
        (grpc::ClientContext * context),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaseek::PurgeRequest, cinaseek::PurgeReply>*),
        AsyncpurgeRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaseek::PurgeRequest, cinaseek::PurgeReply>*),
        PrepareAsyncpurgeRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq),
        (override));
    MOCK_METHOD((grpc::ClientReaderWriterInterface<cinaseek::FindRequest, cinaseek::FindReply>*),
                findRaw,
                (grpc::ClientContext * context),
                (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaseek::FindRequest, cinaseek::FindReply>*),
        AsyncfindRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaseek::FindRequest, cinaseek::FindReply>*),
        PrepareAsyncfindRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq),
        (override));
    MOCK_METHOD((grpc::ClientReaderWriterInterface<cinaseek::InfoRequest, cinaseek::InfoReply>*),
                infoRaw,
                (grpc::ClientContext * context),
                (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaseek::InfoRequest, cinaseek::InfoReply>*),
        AsyncinfoRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaseek::InfoRequest, cinaseek::InfoReply>*),
        PrepareAsyncinfoRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq),
        (override));
    MOCK_METHOD((grpc::ClientReaderWriterInterface<cinaseek::ListRequest, cinaseek::ListReply>*),
                listRaw,
                (grpc::ClientContext * context),
                (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaseek::ListRequest, cinaseek::ListReply>*),
        AsynclistRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaseek::ListRequest, cinaseek::ListReply>*),
        PrepareAsynclistRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq),
        (override));
    MOCK_METHOD(
        (grpc::ClientReaderWriterInterface<cinaseek::NetworksRequest, cinaseek::NetworksReply>*),
        networksRaw,
        (grpc::ClientContext * context),
        (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaseek::NetworksRequest,
                                                        cinaseek::NetworksReply>*),
                AsyncnetworksRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
                (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaseek::NetworksRequest,
                                                        cinaseek::NetworksReply>*),
                PrepareAsyncnetworksRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq),
                (override));
    MOCK_METHOD(
        (grpc::ClientReaderWriterInterface<cinaseek::MountRequest, cinaseek::MountReply>*),
        mountRaw,
        (grpc::ClientContext * context),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaseek::MountRequest, cinaseek::MountReply>*),
        AsyncmountRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaseek::MountRequest, cinaseek::MountReply>*),
        PrepareAsyncmountRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq),
        (override));
    MOCK_METHOD(grpc::ClientAsyncResponseReaderInterface<cinaseek::PingReply>*,
                AsyncpingRaw,
                (grpc::ClientContext * context,
                 const cinaseek::PingRequest& request,
                 grpc::CompletionQueue* cq),
                (override));
    MOCK_METHOD(grpc::ClientAsyncResponseReaderInterface<cinaseek::PingReply>*,
                PrepareAsyncpingRaw,
                (grpc::ClientContext * context,
                 const cinaseek::PingRequest& request,
                 grpc::CompletionQueue* cq),
                (override));
    MOCK_METHOD(
        (grpc::ClientReaderWriterInterface<cinaseek::RecoverRequest, cinaseek::RecoverReply>*),
        recoverRaw,
        (grpc::ClientContext * context),
        (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaseek::RecoverRequest,
                                                        cinaseek::RecoverReply>*),
                AsyncrecoverRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
                (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaseek::RecoverRequest,
                                                        cinaseek::RecoverReply>*),
                PrepareAsyncrecoverRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq),
                (override));
    MOCK_METHOD(
        (grpc::ClientReaderWriterInterface<cinaseek::SSHInfoRequest, cinaseek::SSHInfoReply>*),
        ssh_infoRaw,
        (grpc::ClientContext * context),
        (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaseek::SSHInfoRequest,
                                                        cinaseek::SSHInfoReply>*),
                Asyncssh_infoRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
                (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaseek::SSHInfoRequest,
                                                        cinaseek::SSHInfoReply>*),
                PrepareAsyncssh_infoRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq),
                (override));
    MOCK_METHOD(
        (grpc::ClientReaderWriterInterface<cinaseek::StartRequest, cinaseek::StartReply>*),
        startRaw,
        (grpc::ClientContext * context),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaseek::StartRequest, cinaseek::StartReply>*),
        AsyncstartRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaseek::StartRequest, cinaseek::StartReply>*),
        PrepareAsyncstartRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq),
        (override));
    MOCK_METHOD((grpc::ClientReaderWriterInterface<cinaseek::StopRequest, cinaseek::StopReply>*),
                stopRaw,
                (grpc::ClientContext * context),
                (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaseek::StopRequest, cinaseek::StopReply>*),
        AsyncstopRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaseek::StopRequest, cinaseek::StopReply>*),
        PrepareAsyncstopRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq),
        (override));
    MOCK_METHOD(
        (grpc::ClientReaderWriterInterface<cinaseek::SuspendRequest, cinaseek::SuspendReply>*),
        suspendRaw,
        (grpc::ClientContext * context),
        (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaseek::SuspendRequest,
                                                        cinaseek::SuspendReply>*),
                AsyncsuspendRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
                (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaseek::SuspendRequest,
                                                        cinaseek::SuspendReply>*),
                PrepareAsyncsuspendRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq),
                (override));
    MOCK_METHOD(
        (grpc::ClientReaderWriterInterface<cinaseek::RestartRequest, cinaseek::RestartReply>*),
        restartRaw,
        (grpc::ClientContext * context),
        (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaseek::RestartRequest,
                                                        cinaseek::RestartReply>*),
                AsyncrestartRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
                (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaseek::RestartRequest,
                                                        cinaseek::RestartReply>*),
                PrepareAsyncrestartRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq),
                (override));
    MOCK_METHOD(
        (grpc::ClientReaderWriterInterface<cinaseek::DeleteRequest, cinaseek::DeleteReply>*),
        deletRaw,
        (grpc::ClientContext * context),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaseek::DeleteRequest, cinaseek::DeleteReply>*),
        AsyncdeletRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaseek::DeleteRequest, cinaseek::DeleteReply>*),
        PrepareAsyncdeletRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq),
        (override));
    MOCK_METHOD(
        (grpc::ClientReaderWriterInterface<cinaseek::UmountRequest, cinaseek::UmountReply>*),
        umountRaw,
        (grpc::ClientContext * context),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaseek::UmountRequest, cinaseek::UmountReply>*),
        AsyncumountRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaseek::UmountRequest, cinaseek::UmountReply>*),
        PrepareAsyncumountRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq),
        (override));
    MOCK_METHOD(
        (grpc::ClientReaderWriterInterface<cinaseek::VersionRequest, cinaseek::VersionReply>*),
        versionRaw,
        (grpc::ClientContext * context),
        (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaseek::VersionRequest,
                                                        cinaseek::VersionReply>*),
                AsyncversionRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
                (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaseek::VersionRequest,
                                                        cinaseek::VersionReply>*),
                PrepareAsyncversionRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq),
                (override));
    MOCK_METHOD((grpc::ClientReaderWriterInterface<cinaseek::GetRequest, cinaseek::GetReply>*),
                getRaw,
                (grpc::ClientContext * context),
                (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaseek::GetRequest, cinaseek::GetReply>*),
        AsyncgetRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaseek::GetRequest, cinaseek::GetReply>*),
        PrepareAsyncgetRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq),
        (override));
    MOCK_METHOD((grpc::ClientReaderWriterInterface<cinaseek::SetRequest, cinaseek::SetReply>*),
                setRaw,
                (grpc::ClientContext * context),
                (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaseek::SetRequest, cinaseek::SetReply>*),
        AsyncsetRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaseek::SetRequest, cinaseek::SetReply>*),
        PrepareAsyncsetRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq),
        (override));
    MOCK_METHOD((grpc::ClientReaderWriterInterface<cinaseek::KeysRequest, cinaseek::KeysReply>*),
                keysRaw,
                (grpc::ClientContext * context),
                (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaseek::KeysRequest, cinaseek::KeysReply>*),
        AsynckeysRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaseek::KeysRequest, cinaseek::KeysReply>*),
        PrepareAsynckeysRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq),
        (override));
    MOCK_METHOD((grpc::ClientReaderWriterInterface<cinaseek::AuthenticateRequest,
                                                   cinaseek::AuthenticateReply>*),
                authenticateRaw,
                (grpc::ClientContext * context),
                (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaseek::AuthenticateRequest,
                                                        cinaseek::AuthenticateReply>*),
                AsyncauthenticateRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
                (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaseek::AuthenticateRequest,
                                                        cinaseek::AuthenticateReply>*),
                PrepareAsyncauthenticateRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq),
                (override));
    MOCK_METHOD(
        (grpc::ClientReaderWriterInterface<cinaseek::SnapshotRequest, cinaseek::SnapshotReply>*),
        snapshotRaw,
        (grpc::ClientContext * context),
        (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaseek::SnapshotRequest,
                                                        cinaseek::SnapshotReply>*),
                AsyncsnapshotRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
                (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaseek::SnapshotRequest,
                                                        cinaseek::SnapshotReply>*),
                PrepareAsyncsnapshotRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq),
                (override));
    MOCK_METHOD(
        (grpc::ClientReaderWriterInterface<cinaseek::RestoreRequest, cinaseek::RestoreReply>*),
        restoreRaw,
        (grpc::ClientContext * context),
        (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaseek::RestoreRequest,
                                                        cinaseek::RestoreReply>*),
                AsyncrestoreRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
                (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaseek::RestoreRequest,
                                                        cinaseek::RestoreReply>*),
                PrepareAsyncrestoreRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq),
                (override));
    MOCK_METHOD(
        (grpc::ClientReaderWriterInterface<cinaseek::CloneRequest, cinaseek::CloneReply>*),
        cloneRaw,
        (grpc::ClientContext * context),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaseek::CloneRequest, cinaseek::CloneReply>*),
        AsynccloneRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
        (override));
    MOCK_METHOD(
        (grpc::ClientAsyncReaderWriterInterface<cinaseek::CloneRequest, cinaseek::CloneReply>*),
        PrepareAsynccloneRaw,
        (grpc::ClientContext * context, grpc::CompletionQueue* cq),
        (override));
    MOCK_METHOD((grpc::ClientReaderWriterInterface<cinaseek::DaemonInfoRequest,
                                                   cinaseek::DaemonInfoReply>*),
                daemon_infoRaw,
                (grpc::ClientContext * context),
                (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaseek::DaemonInfoRequest,
                                                        cinaseek::DaemonInfoReply>*),
                Asyncdaemon_infoRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
                (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaseek::DaemonInfoRequest,
                                                        cinaseek::DaemonInfoReply>*),
                PrepareAsyncdaemon_infoRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq),
                (override));
    MOCK_METHOD((grpc::ClientReaderWriterInterface<cinaseek::WaitReadyRequest,
                                                   cinaseek::WaitReadyReply>*),
                wait_readyRaw,
                (grpc::ClientContext * context),
                (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaseek::WaitReadyRequest,
                                                        cinaseek::WaitReadyReply>*),
                Asyncwait_readyRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq, void* tag),
                (override));
    MOCK_METHOD((grpc::ClientAsyncReaderWriterInterface<cinaseek::WaitReadyRequest,
                                                        cinaseek::WaitReadyReply>*),
                PrepareAsyncwait_readyRaw,
                (grpc::ClientContext * context, grpc::CompletionQueue* cq),
                (override));
};
} // namespace cinaseek::test
