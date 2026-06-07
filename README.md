# NullCloud - backend

![CI](https://github.com/we-work-in-the-cloud/backend-nullcloud/actions/workflows/release.yml/badge.svg) ![GitHub release](https://img.shields.io/github/v/release/we-work-in-the-cloud/backend-nullcloud) ![Go version](https://img.shields.io/github/go-mod/go-version/we-work-in-the-cloud/backend-nullcloud) ![License](https://img.shields.io/github/license/we-work-in-the-cloud/backend-nullcloud) ![Downloads](https://img.shields.io/github/downloads/we-work-in-the-cloud/backend-nullcloud/total) ![Homebrew Cask](https://img.shields.io/badge/homebrew-cask-orange?logo=homebrew) ![Made with AI](https://img.shields.io/badge/made%20with%20AI%20-yes-green)

A fake cloud provider API — provision VPCs, subnets, and virtual server instances without any real infrastructure. Useful for demos, tests, and Terraform provider development.

> **Using Terraform?** The [terraform-provider-nullcloud](https://github.com/we-work-in-the-cloud/terraform-provider-nullcloud) wraps this API and lets you manage NullCloud resources with `.tf` files.

## Install

```sh
brew tap we-work-in-the-cloud/backend-nullcloud https://github.com/we-work-in-the-cloud/backend-nullcloud
brew install --cask nullcloud-backend
```

## Resources

- VPC
- Subnet
- Virtual Server Instance (VSI)
- Load Balancer
- Object Storage Bucket
- Managed Database
- Kubernetes Cluster

### Model Classes

```mermaid
classDiagram

    namespace Network {
        class VPC {
            string id
            string name
            string status
            string crn
            string region
            time.Time created_at
        }

        class Subnet {
            string id
            string name
            string status
            string crn
            string vpc_id
            string zone
            string cidr_block
            time.Time created_at
        }

        class LoadBalancer {
            string id
            string name
            string status
            string crn
            string protocol
            int port
            LoadBalancerTarget[] targets
            time.Time created_at
        }

        class LoadBalancerTarget {
            string type
            string id
        }
    }

    namespace Compute {
        class VSI {
            string id
            string name
            string status
            string crn
            string subnet_id
            string profile
            string image
            string primary_ip
            time.Time created_at
        }
        class KubernetesCluster {
            string id
            string name
            string status
            string crn
            string version
            int node_count
            string[] subnet_ids
            time.Time created_at
        }
    }

    namespace Data {
        class Bucket {
            string id
            string name
            string status
            string crn
            string region
            time.Time created_at
        }

        class Database {
            string id
            string name
            string status
            string crn
            string engine
            string version
            string plan
            string[] subnet_ids
            time.Time created_at
        }
    }

    VPC "1" -- "*" Subnet : contains
    LoadBalancerTarget --> VSI : points to vsi
    Subnet "1" -- "*" VSI : hosts
    LoadBalancerTarget --> KubernetesCluster : points to cluster
    Subnet "1" -- "*" Database : uses
    LoadBalancer "1" -- "*" LoadBalancerTarget : targets
    Subnet "1" -- "*" KubernetesCluster : uses
```

## Auth

All requests require an `Authorization` header. Any non-empty string works. Resources are scoped to the token — different tokens are isolated from each other.

## Persistence

| Mode | Flag | Notes |
|------|------|-------|
| In-memory | _(default)_ | State lost on restart |
| JSON file | `--store-file <path>` | State persisted to disk |

## Run

```sh
# In-memory (default)
./nullcloud-backend

# With file persistence
./nullcloud-backend --store-file store.json

# Custom port (default: 8080)
./nullcloud-backend --port 9090
```

## Build

```sh
make build        # all platforms → dist/
make build-linux_amd64  # single platform
```

## Test

```sh
make test
```

## Release

Push a tag to trigger a GitHub Actions release with pre-built binaries for Linux, macOS, and Windows (amd64/arm64).

```sh
git tag v1.0.0 && git push origin v1.0.0
```
