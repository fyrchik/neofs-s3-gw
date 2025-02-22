# Changelog

This document outlines major changes between releases.

## 0.17.0 (30 Aug 2021)
With this release we introduce [ceph-based](https://github.com/ceph/s3-tests) S3 compatibility results.

### Added
* Versioning support (#122)
* Ceph S3 compatibility results (#150)
* Handling `X-Amz-Expected-Bucket-Owner` header (#216)
* `X-Container-Id` header for `HeadBucket` response (#220)
* Basic ACL support (#49, #213)
* Caching (#206, #179)
* Metadata directive when copying (#191)
* Bucket name checking (189)
* Continuation token support (#112, #154, #180)
* Mapping `LocationConstraint` to `PlacementPolicy` (#89)
* Tagging support (#196)

### Changed
* Default placement policy is now configurable (#218) 
* README is split into different files (#210)
* Unified error handling (#184, #149, #89)
* Authmate issue-secret response contains container id (#163)
* Removed "github.com/nspcc-dev/neofs-node" dependency (#234)

### Fixed
* ListObjects results are now the same for different users (#230)
* Error response for invalid authentication header is now correct (#199)
* Saving object metadata (#198)
* Range header handling (#194)
* Correct NotFound status (#118)
* HeadObject for "directories" (#160)
* Fetch-owner parameter support (#159)

## 0.16.0 (16 Jul 2021)

With this release we publish S3 gateway source code. It includes various S3
compatibility improvements, support of bucket management, unified secp256r1
cryptography with NEP-6 wallet support.

### Fixed
 * Allowed no-sign request (#65)
 * Bearer token attached to all requests (#84)
 * Time format in responses (#133)
 * Max-keys checked in ListObjects (#135)
 * Lost metadat in the objects (#131)
 * Unique bucket name check (#125)

### Added
 * Bucket management operations (#47, #72)
 * Node-specific owner IDs in bearer tokens (#83)
 * AWS CLI usage section in README (#77)
 * List object paging (#97)
 * Lifetime for the tokens in auth-mate (#108)
 * Support of range in GetObject request (#96)
 * Support of NEP-6 wallets instead of binary encoded keys (#92)
 * Support of JSON encoded rules in auth-mate (#71)
 * Support of delimiters in ListObjects (#98)
 * Support of object ETag (#93)
 * Support of time-based conditional CopyObject and GetObject (#94)

### Changed
 * Accesskey format: now `0` used as a delimiter between container ID and object 
   ID instead of `_` (#164)
 * Accessbox is encoded in protobuf format (#48)
 * Authentication uses secp256r1 instead of ed25519 (#75)
 * Improved integration with NeoFS SDK and NeoFS API Go (#78, #88)
 * Optimized object put execution (#155)

### Removed
 * GRPC keepalive options (#73)

## 0.15.0 (10 Jun 2021)

This release brings S3 gateway to the current state of NeoFS and fixes some
bugs, no new significant features introduced (other than moving here already
existing authmate component).

New features:
 * authmate was moved into this repository and is now built along with the
   gateway itself (#46)

Behavior changes:
 * neofs-s3-gate was renamed to neofs-s3-gw (#50)

Improvements:
 * better Makefile (#43, #45, #55)
 * stricter linters (#45)
 * removed non-standard errors package from dependencies (#54)
 * refactoring, reusing new sdk-go component (#60, #62, #63)
 * updated neofs-api-go for compatibility with current NeoFS node 0.21.0 (#60, #68)
 * extended README (#67, #76)

Bugs fixed:
 * wrong (as per AWS specification) access key ID generated (#64)

## Older versions

Please refer to [Github
releases](https://github.com/nspcc-dev/neofs-s3-gw/releases/) for older
releases.
