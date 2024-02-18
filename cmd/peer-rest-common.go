// Copyright (c) 2015-2024 MinIO, Inc.
//
// This file is part of MinIO Object Storage stack
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package cmd

const (
	peerRESTVersion       = "v37" // Add 'metrics' option for ServerInfo
	peerRESTVersionPrefix = SlashSeparator + peerRESTVersion
	peerRESTPrefix        = minioReservedBucketPath + "/peer"
	peerRESTPath          = peerRESTPrefix + peerRESTVersionPrefix
)

const (
	peerRESTMethodHealth                 = "/health"
	peerRESTMethodServerInfo             = "/serverinfo"
	peerRESTMethodLocalStorageInfo       = "/localstorageinfo"
	peerRESTMethodCPUInfo                = "/cpuinfo"
	peerRESTMethodDiskHwInfo             = "/diskhwinfo"
	peerRESTMethodNetHwInfo              = "/nethwinfo"
	peerRESTMethodOsInfo                 = "/osinfo"
	peerRESTMethodMemInfo                = "/meminfo"
	peerRESTMethodProcInfo               = "/procinfo"
	peerRESTMethodSysErrors              = "/syserrors"
	peerRESTMethodSysServices            = "/sysservices"
	peerRESTMethodSysConfig              = "/sysconfig"
	peerRESTMethodGetBucketStats         = "/getbucketstats"
	peerRESTMethodGetAllBucketStats      = "/getallbucketstats"
	peerRESTMethodVerifyBinary           = "/verifybinary"
	peerRESTMethodCommitBinary           = "/commitbinary"
	peerRESTMethodSignalService          = "/signalservice"
	peerRESTMethodBackgroundHealStatus   = "/backgroundhealstatus"
	peerRESTMethodGetLocks               = "/getlocks"
	peerRESTMethodStartProfiling         = "/startprofiling"
	peerRESTMethodDownloadProfilingData  = "/downloadprofilingdata"
	peerRESTMethodLog                    = "/log"
	peerRESTMethodGetBandwidth           = "/bandwidth"
	peerRESTMethodGetMetacacheListing    = "/getmetacache"
	peerRESTMethodUpdateMetacacheListing = "/updatemetacache"
	peerRESTMethodGetPeerMetrics         = "/peermetrics"
	peerRESTMethodGetPeerBucketMetrics   = "/peerbucketmetrics"
	peerRESTMethodSpeedTest              = "/speedtest"
	peerRESTMethodDriveSpeedTest         = "/drivespeedtest"
	peerRESTMethodStopRebalance          = "/stoprebalance"
	peerRESTMethodGetLastDayTierStats    = "/getlastdaytierstats"
	peerRESTMethodDevNull                = "/devnull"
	peerRESTMethodNetperf                = "/netperf"
	peerRESTMethodMetrics                = "/metrics"
	peerRESTMethodResourceMetrics        = "/resourcemetrics"
	peerRESTMethodGetReplicationMRF      = "/getreplicationmrf"
	peerRESTMethodGetSRMetrics           = "/getsrmetrics"
)

const (
	peerRESTBucket         = "bucket"
	peerRESTBuckets        = "buckets"
	peerRESTUser           = "user"
	peerRESTGroup          = "group"
	peerRESTUserTemp       = "user-temp"
	peerRESTPolicy         = "policy"
	peerRESTUserOrGroup    = "user-or-group"
	peerRESTUserType       = "user-type"
	peerRESTIsGroup        = "is-group"
	peerRESTSignal         = "signal"
	peerRESTSubSys         = "sub-sys"
	peerRESTProfiler       = "profiler"
	peerRESTSize           = "size"
	peerRESTConcurrent     = "concurrent"
	peerRESTDuration       = "duration"
	peerRESTStorageClass   = "storage-class"
	peerRESTEnableSha256   = "enableSha256"
	peerRESTMetricsTypes   = "types"
	peerRESTDisk           = "disk"
	peerRESTHost           = "host"
	peerRESTJobID          = "job-id"
	peerRESTDepID          = "depID"
	peerRESTStartRebalance = "start-rebalance"
	peerRESTMetrics        = "metrics"
	peerRESTDryRun         = "dry-run"

	peerRESTURL         = "url"
	peerRESTSha256Sum   = "sha256sum"
	peerRESTReleaseInfo = "releaseinfo"

	peerRESTListenBucket = "bucket"
	peerRESTListenPrefix = "prefix"
	peerRESTListenSuffix = "suffix"
	peerRESTListenEvents = "events"
)
