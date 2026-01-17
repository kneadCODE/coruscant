# Azure Global Architecture

## Executive Summary

This document outlines our Azure global architecture for a production-grade, GDPR-compliant SaaS platform. Our architecture uses **4 Azure regions across 2 paired clusters** with **separate hub infrastructure per environment** to balance regulatory compliance, operational resilience, security isolation, and cost efficiency.

**Selected Configuration:**

- **Global Pair (Non-GDPR):** East US 2 (Active) + Central US (Standby)
- **GDPR Pair (EU/EEA):** West Europe (Active) + North Europe (Standby)

**Hub Strategy:**

- **6 Total Hubs:** 4 Production + 2 Non-Production
- Separate Prod/NonProd hubs for security isolation and testing independence
- Shared connectivity components (ER Circuit, LNG) for cost optimization

---

## Table of Contents

1. [Why Two Geographic Pairs](#why-two-geographic-pairs)
2. [Why Separate Production and Non-Production Hubs](#why-separate-production-and-non-production-hubs)
3. [Region Selection Rationale](#region-selection-rationale)
4. [Hub Architecture Strategy](#hub-architecture-strategy)
5. [Network Connectivity](#network-connectivity)
6. [Data Services Architecture](#data-services-architecture)
7. [Observability & Security](#observability--security)
8. [Backup & Disaster Recovery](#backup--disaster-recovery)
9. [DevOps & Infrastructure Management](#devops--infrastructure-management)
10. [Cost Optimization Strategies](#cost-optimization-strategies)
11. [Decision Audit Trail](#decision-audit-trail)

---

## Why Two Geographic Pairs

### 1. GDPR Compliance Requirements

**The Challenge:**
Our application serves a global customer base spanning:

- **Americas:** US, Canada, Argentina, Brazil, Mexico
- **Europe:** UK, EU-27, Switzerland
- **Asia-Pacific:** Japan, Singapore, Indonesia, Australia, New Zealand
- **Middle East:** UAE, Saudi Arabia
- **Other:** India, South Africa

Under **GDPR (General Data Protection Regulation)**, personal data of EU/UK residents must:

- Remain within EU/EEA boundaries unless specific adequacy agreements exist
- Not be transferred to jurisdictions without adequate data protection
- Be subject to EU data protection laws and enforcement

**The Solution:**
By maintaining a **dedicated GDPR cluster** in European regions:

- ✅ EU customer data never leaves EU boundaries
- ✅ Full compliance with GDPR Article 44-50 (data transfers)
- ✅ Simplified regulatory audits and certifications
- ✅ Clear technical enforcement of data residency

### 2. Operational Resilience & Blast Radius Containment

**Geographic Risk Isolation:**

- Catastrophic events in North America don't impact EU operations
- Regulatory actions are geographically isolated
- Network routing issues don't cascade globally

**Regulatory Boundary Enforcement:**

- Clear technical separation prevents accidental data spillover
- Simplified compliance attestation
- Each geography can operate independently

### 3. Performance Optimization

**Latency Optimization by Market:**

| Customer Location | Routing | Typical RTT |
|-------------------|---------|-------------|
| US, Canada, LatAm | → East US 2 | 20-140ms |
| EU, UK, Middle East | → West Europe | 10-110ms |
| Asia-Pacific | → East US 2 | 150-200ms |

**Benefits:**

- 30-50% lower latency for EU customers vs. single US deployment
- Compliance with data residency while maintaining performance
- Foundation for future regional expansion (Singapore, Japan, etc.)

---

## Why Separate Production and Non-Production Hubs

### The Architecture Decision

We implement **separate hub infrastructure per environment** rather than shared hubs because our requirements prioritize:

1. **Security isolation** - SOC 2 Type 2, ISO 27001, PCI-DSS compliance
2. **Testing independence** - Ability to test infrastructure changes safely
3. **GitOps maturity** - Infrastructure-as-Code makes managing multiple hubs trivial

### Key Benefits of Separate Hubs

**1. Security Isolation & Blast Radius**

**Separate Hubs:**

- ✅ NonProd firewall misconfiguration **cannot** affect Production (different firewall)
- ✅ NonProd VPN tunnel issues **cannot** affect Production connectivity
- ✅ Network-level air gap between environments
- ✅ Compromising NonProd doesn't provide network path to Production

**vs Shared Hub:**

- ⚠️ Single firewall = single point of failure for both environments
- ⚠️ Misconfigured rule can impact both Prod and NonProd
- ⚠️ Requires complex rule management to maintain isolation

**2. True Infrastructure Testing**

**Test Scenarios Enabled by Separate Hubs:**

| Scenario | Separate Hubs | Shared Hub |
|----------|---------------|------------|
| **Test new firewall rules** | Add to NonProd FW → validate → apply to Prod | Cannot test without affecting Prod |
| **Test VPN failover** | Break NonProd VPN → verify procedures | Cannot test (shared gateway) |
| **Test network architecture migration** | Migrate NonProd → validate for months → migrate Prod | Cannot test (shared infrastructure) |
| **Test disaster recovery** | Simulate NonProd hub failure → Prod unaffected | Simulation affects both environments |

With separate hubs, **NonProd architecturally mirrors Production**, enabling true production-parity testing.

**3. Compliance & Audit Simplicity**

**SOC 2 Type 2 Audit:**

- **Separate Hubs:** "Production and Non-Production are physically separate networks" → Simple 5-minute conversation
- **Shared Hub:** Must demonstrate RBAC + firewall rules + NSGs + GitOps controls → Complex 2-hour documentation review

**ISO 27001 Control A.13.1.3 (Network Segregation):**

- **Separate Hubs:** Native compliance through physical separation
- **Shared Hub:** Requires compensating controls and extensive documentation

**PCI-DSS Requirement 1.2 (Network Segmentation):**

- **Separate Hubs:** Clear segmentation that auditors can visualize
- **Shared Hub:** Must prove effective logical segmentation (quarterly pen-tests)

**4. Operational Clarity with GitOps**

With Infrastructure-as-Code (OpenTofu + GitHub Actions):

```
/infrastructure
  ├─ modules/
  │  └─ hub/              # Single hub module
  ├─ prod-hub/
  │  └─ main.tf           # Deploys hub module with prod config
  └─ nonprod-hub/
     └─ main.tf           # Deploys hub module with nonprod config
```

**Managing 2 hubs vs 1 shared hub:**

- Same Terraform module, different variable files
- Separate CODEOWNERS: Prod hub = Security team, NonProd hub = Platform team
- NonProd changes don't require Production-level change control
- Clear separation of concerns in IaC repository

**The complexity argument doesn't apply when using GitOps.**

### What Gets Shared vs Separated

**SHARED Between Prod and NonProd:**

- ✅ **ExpressRoute Circuit** (physical circuit)
- ✅ **Local Network Gateway** (logical representation of on-prem)
- ✅ **Security LAW + Sentinel** (security needs unified view across environments)

**SEPARATED Per Environment:**

- ✅ **Hub VNet** (complete network isolation)
- ✅ **Azure Firewall** (independent security policies)
- ✅ **VPN Gateway** (separate tunnels to on-prem)
- ✅ **ER Gateway** (both connect to shared circuit)
- ✅ **EPA Connectors** (separate identity access)
- ✅ **Azure Bastion** (independent breakglass)
- ✅ **Management LAW** (separate operational logs)
- ✅ **Recovery Services Vault** (different retention policies)
- ✅ **Backup Vault** (different backup strategies)

---

## Region Selection Rationale

### Global Pair: East US 2 + Central US

#### Primary: East US 2 (Virginia)

**Why Selected:**

- ✅ Azure's flagship US region - first to receive new features
- ✅ Mature 3 Availability Zones (highest availability SLA)
- ✅ Excellent submarine cable infrastructure (MAREA to Europe, BRUSA to South America)
- ✅ Optimal latency balance: 85ms to Europe, 150ms to APAC, 140ms to LatAm
- ✅ Highest capacity and proven enterprise-scale reliability
- ✅ Baseline Azure pricing (1.00x reference point)

#### Standby: Central US (Iowa)

**Why Selected:**

- ✅ Official Azure paired region (native GRS/GZRS support)
- ✅ Geographic separation: 1,100km from East US 2 (regional disaster mitigation)
- ✅ Mature 3 Availability Zones
- ✅ Slightly lower cost: 2% cheaper than East US 2
- ✅ Geographic center of North America (balanced coast-to-coast latency)
- ✅ Gold feature tier (identical service catalog to East US 2)

**Alternatives Considered:**

| Region Pair | Why Rejected |
|-------------|--------------|
| South Central US + North Central US | North Central US AZs still rolling out (completion late 2026); 10ms worse global latency; East US 2 gets new features first |
| West US 2 + West US 3 | Higher cost; worse latency to Europe and East Coast US |

### GDPR Pair: West Europe + North Europe

#### Primary: West Europe (Netherlands)

**Why Selected:**

- ✅ Azure's flagship EU region - largest European datacenter hub
- ✅ Central European location: optimal latency to UK (35ms), Germany (40ms), France (30ms)
- ✅ Mature 3 Availability Zones
- ✅ EU baseline pricing (1.00x EU reference)
- ✅ Excellent global connectivity to Middle East, India, Africa
- ✅ GDPR compliant - data remains within EU boundaries
- ✅ Native GRS pairing with North Europe (automatic replication)

#### Standby: North Europe (Ireland)

**Why Selected:**

- ✅ Official Azure paired region (native GRS/GZRS support)
- ✅ Geographic separation: 750km from Netherlands (different disaster zone)
- ✅ Mature 3 Availability Zones
- ✅ Identical pricing to West Europe
- ✅ Gold feature tier (complete service catalog parity)
- ✅ GDPR compliant - data remains within EU/EEA

**Alternatives Considered:**

| Region | Why Rejected |
|--------|--------------|
| Sweden Central + Sweden South | 8% cheaper BUT Sweden South is restricted "Alternate" region; native GRS only goes to Sweden South (not West Europe); requires manual Object Replication; operational complexity outweighs savings |
| Germany West Central | 8% cost premium over West Europe with no latency benefit; only justified for strict German data residency |
| UK South | 8% premium; worse connectivity to continental Europe; post-Brexit data adequacy concerns |
| France Central | Similar price to West Europe but smaller capacity; West Europe is more established |

**Operational Simplicity Matters:**

West Europe + North Europe provides **zero-configuration disaster recovery** via native Azure GRS/GZRS, eliminating:

- Custom replication scripts and monitoring
- Manual failover procedures
- Ongoing maintenance of custom DR architecture
- Training overhead for operations team

This is the configuration used by most Fortune 500 companies for European Azure deployments.

---

## Hub Architecture Strategy

### Complete Hub Topology

```
┌─────────────────────────────────────────────────────────┐
│            GLOBAL PAIR (Non-GDPR)                       │
├─────────────────────────────────────────────────────────┤
│ East US 2 (Active):                                     │
│ ├─ Production Hub                                       │
│ │  ├─ Azure Firewall Premium                           │
│ │  ├─ VPN Gateway                                       │
│ │  ├─ ExpressRoute Gateway                              │
│ │  ├─ ExpressRoute Circuit (shared)                    │
│ │  ├─ EPA Connectors                                    │
│ │  ├─ Azure Bastion (breakglass)                       │
│ │  └─ P2S VPN (breakglass)                             │
│ │                                                        │
│ └─ Non-Production Hub                                   │
│    ├─ Azure Firewall Standard                           │
│    ├─ VPN Gateway (separate tunnel)                    │
│    ├─ ExpressRoute Gateway (shares circuit)            │
│    ├─ EPA Connectors                                    │
│    ├─ Azure Bastion                                     │
│    └─ Routes to: Dev, QA, Staging spokes               │
│                                                          │
│ Central US (Standby):                                   │
│ ├─ Production Hub                                       │
│ │  ├─ Azure Firewall Standard                          │
│ │  ├─ EPA Connectors                                    │
│ │  ├─ No VPN/ER (uses East US 2 via peering)          │
│ │  └─ Serves 10% read-only traffic                     │
│ │                                                        │
│ └─ Non-Production Hub (On-Demand)                       │
│    ├─ VNet exists (FREE when empty)                    │
│    ├─ Deploy for DR exercises: 2× per year, 15 days   │
│    └─ Cost-optimized: Only pay when running            │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│            GDPR PAIR (EU/EEA)                           │
├─────────────────────────────────────────────────────────┤
│ West Europe (Active):                                   │
│ ├─ Production Hub                                       │
│ │  ├─ Azure Firewall Premium                           │
│ │  ├─ VPN Gateway                                       │
│ │  ├─ ExpressRoute Gateway                              │
│ │  ├─ ExpressRoute Circuit (shared)                    │
│ │  ├─ EPA Connectors                                    │
│ │  ├─ Azure Bastion (breakglass)                       │
│ │  └─ P2S VPN (breakglass)                             │
│ │                                                        │
│ └─ No Non-Production (cost optimization)               │
│                                                          │
│ North Europe (Standby):                                 │
│ ├─ Production Hub                                       │
│ │  ├─ Azure Firewall Standard                          │
│ │  ├─ EPA Connectors                                    │
│ │  ├─ No VPN/ER (uses West Europe via peering)        │
│ │  └─ Serves 10% read-only traffic                     │
│ │                                                        │
│ └─ No Non-Production (cost optimization)               │
└─────────────────────────────────────────────────────────┘

TOTAL HUBS: 6
├─ 4 Production Hubs (all regions)
└─ 2 Non-Production Hubs (Global pair only)
```

### Hub Components by Region

| Component | East US 2 (Active) | Central US (Standby) | West Europe (Active) | North Europe (Standby) |
|-----------|-------------------|----------------------|---------------------|------------------------|
| **Prod Hub** | Full (Premium FW) | Minimal (Std FW) | Full (Premium FW) | Minimal (Std FW) |
| **NonProd Hub** | Full (Std FW) | On-demand | None | None |
| **VPN Gateway** | In both hubs | Uses East US 2 | In Prod hub | Uses West EU |
| **ER Gateway** | In both hubs | Uses East US 2 | In Prod hub | Uses West EU |
| **ER Circuit** | 1 circuit (shared) | Shares East US 2 | 1 circuit (shared) | Shares West EU |
| **EPA Connectors** | Both hubs | Both hubs | Prod hub only | Prod hub only |
| **Bastion** | Both hubs | Prod hub + on-demand | Prod hub only | Prod hub only |
| **P2S VPN** | Prod hub | Deploy on-demand | Prod hub | Deploy on-demand |

### Hub Sizing Strategy

**Active Region Production Hubs:**

- Azure Firewall Premium (TLS inspection, IDPS)
- VPN Gateway: VpnGw2 or higher
- ER Gateway: Standard or HighPerformance
- Azure Bastion: Standard SKU

**Active Region Non-Production Hubs:**

- Azure Firewall Standard (cost optimization)
- VPN Gateway: VpnGw1 (sufficient for dev/test)
- ER Gateway: Standard
- Azure Bastion: Basic SKU

**Standby Region Production Hubs:**

- Azure Firewall Standard (handle 10-30% traffic)
- No VPN/ER (use active via peering)
- Azure Bastion: Basic SKU or on-demand

**Standby Region Non-Production Hubs:**

- On-demand deployment only
- Deploy 2× per year for DR exercises
- Empty VNet costs $0 when not in use

---

## Network Connectivity

### Hub Peering Strategy

**Within Each Pair (Always Connected):**

```
Global Pair:
East US 2 Prod Hub ↔ Central US Prod Hub (global VNet peering)
East US 2 NonProd Hub ↔ Central US NonProd Hub (global VNet peering)

GDPR Pair:
West Europe Prod Hub ↔ North Europe Prod Hub (global VNet peering)
```

**Between Pairs (NOT Connected):**

```
Global Pair ↔ GDPR Pair: NO CONNECTION

Why:
- GDPR isolation (US data should not mix with EU data)
- Clearer audit story
- Each geography self-sufficient
- DevOps infrastructure duplicated to both pairs
```

### ExpressRoute Strategy

**Physical Circuits: 2 Total**

- 1 circuit in East US 2 region (serves Global pair)
- 1 circuit in West Europe region (serves GDPR pair)

**Logical Gateways:**

- East US 2 Prod Hub: ER Gateway → connects to circuit
- East US 2 NonProd Hub: ER Gateway → connects to same circuit
- Central US: Uses East US 2 gateways via hub peering
- West Europe Prod Hub: ER Gateway → connects to circuit
- North Europe: Uses West Europe gateway via hub peering

**Cost Optimization:**

- Physical circuit is expensive (dedicated fiber)
- Share circuit between Prod and NonProd environments
- Standby regions don't need separate circuits
- On-prem typically has 1-2 circuits total, not 4+

### VPN Gateway Strategy

**Site-to-Site VPN:**

- Separate VPN tunnels per environment (Prod vs NonProd)
- Active regions: Always-on VPN gateways
- Standby regions: Use active region gateways via peering
- Can be deployed on-demand to standby if primary region fails

**Point-to-Site VPN (Breakglass):**

- Deployed in active Production hubs only
- Emergency access when EPA or Entra Private Access fails
- Deploy to standby on-demand if needed during DR

### Entra Private Access (EPA) Strategy

**Primary Access Method:**

- EPA connectors are the primary method for human access to private networks
- Deployed in all hubs (Prod and NonProd, Active and Standby)
- Separate connector groups per environment
- Identity-based access control (no VPN needed for normal operations)

**Breakglass Access (EPA Failure Scenario):**

- Azure Bastion: Jump box access to VMs
- P2S VPN: Direct VPN when EPA completely fails
- Only in active Production hubs normally
- Deploy to standby if primary region + EPA both fail

---

## Data Services Architecture

### Overview of Data Services

Our platform uses a hybrid approach combining Azure-managed PaaS services with self-hosted open-source databases to balance manageability, performance, and feature requirements.

**Azure-Managed Services:**

- Azure Database for PostgreSQL Flexible Server
- Azure Storage Account
- Azure Managed Instance for Apache Cassandra
- Azure Cache for Redis

**Self-Hosted Services:**

- Elasticsearch (search and analytics)
- TimescaleDB (time-series data)
- MongoDB (document store)
- Apache Kafka (event streaming, short-term persistence)

### Data Service Replication Strategy

#### Azure Database for PostgreSQL Flexible Server

**Configuration:**

```
Global Pair:
├─ East US 2 (Primary)
│  ├─ PostgreSQL Flexible Server (Zone-redundant HA)
│  ├─ 3 Availability Zones
│  └─ Read replicas in Central US (async replication)
│
└─ Central US (Standby)
   └─ Read replica (promoted to primary during failover)

GDPR Pair:
├─ West Europe (Primary)
│  ├─ PostgreSQL Flexible Server (Zone-redundant HA)
│  └─ Read replicas in North Europe
│
└─ North Europe (Standby)
   └─ Read replica
```

**Replication Details:**

- Built-in Azure geo-replication (private endpoints only)
- RPO: <1 minute (async replication lag)
- RTO: <5 minutes (promote replica to primary)
- No cross-region VNet peering required (uses Azure backbone)
- Private Link for all connections

**Backup Strategy:**

- Automated backups with PITR (Point-In-Time Recovery) up to 35 days
- Backups stored with GRS (automatically replicated to paired region)
- Production: 35-day retention
- Non-Production: 7-day retention

#### Azure Storage Account

**Configuration:**

```
Production Storage:
├─ East US 2: GZRS (Geo-Zone-Redundant Storage)
│  ├─ Synchronous replication across 3 AZs
│  └─ Async replication to Central US
│
└─ West Europe: GZRS
   ├─ Synchronous replication across 3 AZs
   └─ Async replication to North Europe

Non-Production Storage:
└─ LRS (Locally Redundant Storage) - Cost optimization
```

**Features:**

- 16 nines durability (99.99999999999999%)
- Private endpoints for all access
- Immutable storage for compliance data
- Lifecycle management policies

#### Azure Managed Instance for Apache Cassandra

**Configuration:**

```
Global Pair:
├─ East US 2 (Primary datacenter)
│  ├─ 3 nodes across 3 Availability Zones
│  └─ Cassandra replication factor: 3
│
└─ Central US (Standby datacenter)
   └─ 3 nodes (async replication via Cassandra protocol)

GDPR Pair:
├─ West Europe (Primary datacenter)
└─ North Europe (Standby datacenter)
```

**Replication Details:**

- Multi-datacenter replication (Cassandra native)
- NetworkTopologyStrategy for geo-replication
- Eventually consistent reads from standby regions
- Private VNet integration (no public endpoints)

**Consistency Configuration:**

- Local_Quorum for writes (within same datacenter)
- Local_Quorum for reads (low latency)
- Quorum for critical reads (ensure consistency)

#### Azure Cache for Redis

**Configuration:**

```
Production:
├─ Enterprise tier with active geo-replication
├─ Zone-redundant (3 replicas across AZs)
└─ Active-Active replication between regions

Per Region:
├─ Primary cache in active region
├─ Read replica in standby region
└─ Automatic failover configured
```

**Use Cases:**

- Session management
- API response caching
- Rate limiting counters
- Real-time leaderboards

### Self-Hosted Data Services

#### Elasticsearch Cluster

**Deployment Model:**

```
Per Region (Production):
├─ Master nodes: 3 (across 3 AZs)
├─ Data nodes: 6+ (scaled based on load)
├─ Coordinating nodes: 2 (query routing)
└─ Deployed in dedicated spoke VNet

Replication:
├─ East US 2 ↔ Central US: Cross-Cluster Replication (CCR)
└─ West Europe ↔ North Europe: Cross-Cluster Replication (CCR)
```

**Network Configuration:**

- Spoke-to-spoke VNet peering for CCR traffic
- Private IPs only (no public endpoints)
- Firewall rules control which spokes can connect
- Elasticsearch security features enabled (TLS, authentication)

**Backup Strategy:**

- Snapshot to Azure Storage Account (GZRS)
- Automated daily snapshots
- Retention: 30 days production, 7 days non-production

**Why Self-Hosted:**

- Advanced search features not available in Azure Cognitive Search
- Custom analyzers and plugins
- Fine-grained performance tuning
- Cost optimization for large datasets

#### TimescaleDB (PostgreSQL Extension)

**Deployment Model:**

```
Per Region (Production):
├─ Primary VM: Premium SSD, zone-redundant
├─ Standby VM: Async streaming replication
└─ Deployed in database spoke VNet

Replication:
├─ East US 2 ↔ Central US: PostgreSQL streaming replication
└─ West Europe ↔ North Europe: PostgreSQL streaming replication
```

**Network Configuration:**

- Spoke-to-spoke VNet peering for replication traffic
- Private IPs with NSG rules
- TLS required for all connections
- Connection pooling via PgBouncer

**Backup Strategy:**

- pg_basebackup + WAL archiving to Azure Storage (GZRS)
- PITR capability up to 30 days
- Continuous WAL archiving

**Why Self-Hosted:**

- TimescaleDB-specific features (continuous aggregates, compression)
- Not available as Azure managed service
- Performance requirements for time-series queries

#### MongoDB Cluster

**Deployment Model:**

```
Per Region (Production):
├─ Replica Set: 3 members (across 3 AZs)
├─ Config Servers: 3 (for sharded cluster)
├─ Mongos Routers: 2+ (application connection points)
└─ Deployed in database spoke VNet

Replication:
├─ East US 2 ↔ Central US: MongoDB global cluster
└─ West Europe ↔ North Europe: MongoDB global cluster
```

**Network Configuration:**

- Spoke-to-spoke VNet peering between regions
- Private IPs with MongoDB authentication
- TLS/SSL for all inter-node communication
- Firewall rules restrict access to application spokes only

**Backup Strategy:**

- Continuous backup via mongodump to Azure Storage (GZRS)
- Oplog tailing for PITR
- Retention: 30 days production, 7 days non-production

**Why Self-Hosted:**

- Specific MongoDB features required
- Azure Cosmos DB MongoDB API has limitations
- Cost optimization for workload patterns
- Full control over sharding strategy

#### Apache Kafka Cluster

**Deployment Model:**

```
Per Region (Production):
├─ Kafka Brokers: 6+ (across 3 AZs, 2 per AZ)
├─ Zookeeper: 3 nodes (or KRaft mode)
├─ Schema Registry: 2 nodes
├─ Kafka Connect: 2+ nodes
└─ Deployed in messaging spoke VNet

Replication:
├─ East US 2 ↔ Central US: MirrorMaker 2 (disaster recovery)
└─ West Europe ↔ North Europe: MirrorMaker 2
```

**Network Configuration:**

- Spoke-to-spoke peering for MirrorMaker replication
- Private IPs with SASL/SCRAM authentication
- TLS encryption for all traffic
- Dedicated spoke for messaging infrastructure

**Retention Strategy:**

- Short-term persistence: 7 days (operational events)
- Long-term storage: Archive to Azure Storage via Kafka Connect
- Compacted topics for changelog streams

**Why Self-Hosted vs Azure Event Hubs:**

- Kafka-specific features (transactions, exactly-once semantics)
- Custom Kafka Connect connectors
- Kafka Streams applications
- Fine-grained performance control
- Cost optimization for high-throughput scenarios

### Data Replication Network Architecture

**Key Principle: Spoke-to-Spoke Replication**

```
Database replication does NOT go through hub firewalls.
Direct spoke-to-spoke VNet peering is used.

Example - PostgreSQL Replication:
East US 2 DB Spoke ↔ VNet Peering ↔ Central US DB Spoke
    (Direct connection, private IPs)

Hub Firewall Role:
- Controls which spokes can peer with each other
- Monitors traffic (flow logs)
- Does NOT route database replication traffic
```

**Network Flow for Self-Hosted Services:**

```
East US 2:
├─ App Spoke → DB Spoke (same region, intra-region peering)
└─ DB Spoke → Central US DB Spoke (cross-region peering)
   ├─ Private IP replication
   ├─ TLS encrypted
   └─ Firewall monitors but doesn't route
```

**Benefits of Spoke-to-Spoke Replication:**

- ✅ Lower latency (direct connection, no hub transit)
- ✅ Hub firewall doesn't become bottleneck
- ✅ Simplified routing (no UDRs needed)
- ✅ Clear network segmentation (DB spoke isolated)

### Data Service Backup Summary

| Service | Backup Method | Retention (Prod) | Retention (NonProd) | Storage Location |
|---------|---------------|------------------|---------------------|------------------|
| **PostgreSQL Flexible** | Automated + PITR | 35 days | 7 days | GRS (automatic) |
| **Storage Account** | GZRS replication | N/A (continuous) | N/A | Paired region |
| **Cassandra** | Cassandra backup | 30 days | 7 days | Azure Storage (GZRS) |
| **Redis** | RDB snapshots | 7 days | 1 day | Azure Storage (GZRS) |
| **Elasticsearch** | Snapshots | 30 days | 7 days | Azure Storage (GZRS) |
| **TimescaleDB** | pg_basebackup + WAL | 30 days (PITR) | 7 days | Azure Storage (GZRS) |
| **MongoDB** | mongodump + oplog | 30 days (PITR) | 7 days | Azure Storage (GZRS) |
| **Kafka** | Topic replication | 7 days (in Kafka) | 7 days | Archive to Storage |

---

## Observability & Security

### Log Analytics Workspace Strategy

#### Architecture Overview

We implement **geographic separation for all observability infrastructure** to ensure GDPR compliance and data residency. This means EU infrastructure logs never leave EU boundaries.

**Total: 4 Log Analytics Workspaces**

```
┌─────────────────────────────────────────────────────┐
│              GLOBAL PAIR (East US 2)                │
├─────────────────────────────────────────────────────┤
│                                                      │
│ Management LAW - Global                             │
│ ├─ Ingests: Operational logs from Global pair      │
│ │  ├─ East US 2: All resources (Prod + NonProd)   │
│ │  └─ Central US: All resources (Prod + NonProd)   │
│ ├─ Log Types:                                       │
│ │  ├─ VM metrics (CPU, memory, disk, network)      │
│ │  ├─ Database performance metrics                 │
│ │  ├─ Storage account metrics                      │
│ │  ├─ Kafka/Elasticsearch/MongoDB metrics          │
│ │  ├─ Autoscaling events                           │
│ │  ├─ Backup job logs                              │
│ │  ├─ Database slow query logs                     │
│ │  └─ Application performance data                 │
│ ├─ Retention: 30 days (hot querying)               │
│ ├─ Access: Platform team, SRE, Database team       │
│ └─ Cost: Pay-per-GB ingestion + queries            │
│                                                      │
│ Security LAW - Global + Sentinel                    │
│ ├─ Ingests: Security logs from Global pair         │
│ │  ├─ East US 2: All security events              │
│ │  └─ Central US: All security events              │
│ ├─ Log Types:                                       │
│ │  ├─ Azure Firewall logs (all traffic, threats)  │
│ │  ├─ NSG Flow Logs (network patterns)             │
│ │  ├─ WAF logs (Application Gateway, Front Door)   │
│ │  ├─ VPN Gateway diagnostics                      │
│ │  ├─ Azure AD Sign-in logs (Global users)        │
│ │  ├─ Azure AD Audit logs                          │
│ │  ├─ Entra Private Access logs                    │
│ │  ├─ Key Vault access logs                        │
│ │  ├─ Database audit logs (auth failures)          │
│ │  ├─ Defender for Cloud alerts                    │
│ │  ├─ Activity Logs (resource changes)             │
│ │  └─ VM authentication logs                       │
│ ├─ Microsoft Sentinel: Enabled ✅                  │
│ ├─ Retention: 30 days (hot querying)               │
│ ├─ Access: Security team (Global scope)            │
│ └─ Monitors: Global pair only (US/APAC customers)  │
│                                                      │
└─────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────┐
│             GDPR PAIR (West Europe)                 │
├─────────────────────────────────────────────────────┤
│                                                      │
│ Management LAW - GDPR                               │
│ ├─ Ingests: Operational logs from GDPR pair        │
│ │  ├─ West Europe: All resources (Prod only)       │
│ │  └─ North Europe: All resources (Prod only)       │
│ ├─ Same log types as Global Management LAW         │
│ ├─ Retention: 30 days                              │
│ ├─ Access: Platform team (GDPR scope)              │
│ └─ Data residency: EU only ✅                      │
│                                                      │
│ Security LAW - GDPR + Sentinel                      │
│ ├─ Ingests: Security logs from GDPR pair           │
│ │  ├─ West Europe: All security events             │
│ │  └─ North Europe: All security events             │
│ ├─ Same log types as Global Security LAW           │
│ ├─ Microsoft Sentinel: Enabled ✅                  │
│ ├─ Retention: 30 days                              │
│ ├─ Access: Security team (GDPR scope)              │
│ ├─ Monitors: GDPR pair only (EU customers)         │
│ └─ Data residency: EU only ✅                      │
│                                                      │
└─────────────────────────────────────────────────────┘
```

---

### Why Geographic Separation?

**GDPR Compliance:**

- EU infrastructure logs (even operational metrics) remain in EU
- No EU data transferred to US infrastructure
- Simplified compliance attestation
- Clear data residency boundaries

**Operational Benefits:**

- Each geography can operate independently
- Different teams can manage different regions
- Regional incident response doesn't require cross-geography access
- Aligns with application data residency model

**Consistency:**

- Mirrors storage account geographic split
- Mirrors hub infrastructure geographic split
- Consistent architecture across all layers

---

### Log Categorization: Management vs Security

#### Management LAW (Operational Telemetry)

**Purpose:** Performance monitoring, capacity planning, troubleshooting

**Log Types:**

- **Compute Metrics:** VM CPU, memory, disk I/O, network throughput
- **Database Performance:** Query execution time, connection counts, replication lag
- **Storage Metrics:** IOPS, throughput, capacity, transaction counts
- **Application Performance:** Response times, error rates, dependency calls
- **Infrastructure Events:** Autoscaling actions, backup job status, resource provisioning
- **Self-Hosted Services:** Elasticsearch cluster health, Kafka broker metrics, MongoDB replica set status

**Who Uses It:**

- Platform Engineering team
- Site Reliability Engineers (SRE)
- Database Administrators
- Application Developers (for performance tuning)

**Query Patterns:**

- "Why is this VM slow?"
- "Is database replication lagging?"
- "Do we need to scale up?"
- "What caused this application timeout?"

---

#### Security LAW (Threat Detection)

**Purpose:** Security monitoring, threat detection, compliance, incident response

**Log Types:**

**Network Security:**

- Azure Firewall logs (allowed/denied traffic, threat intelligence matches)
- NSG Flow Logs (network traffic patterns, anomalies)
- WAF logs (Application Gateway, Front Door - injection attempts, suspicious patterns)
- VPN Gateway diagnostics (connection attempts, authentication failures)
- DDoS Protection logs (attack patterns, mitigation actions)
- Private DNS query logs (DNS tunneling detection)

**Identity & Access:**

- Azure AD Sign-in logs (failed authentications, risky sign-ins, impossible travel)
- Azure AD Audit logs (role changes, permission modifications)
- Entra Private Access logs (connector health, user access patterns)
- Conditional Access logs (policy evaluations, blocks)
- Privileged Identity Management (PIM) logs (elevation requests, approvals)

**Security Services:**

- Microsoft Defender for Cloud alerts (vulnerabilities, misconfigurations)
- Defender for Servers (malware detection, suspicious processes)
- Defender for Databases (SQL injection attempts, anomalous access)
- Defender for Storage (malware uploads, suspicious downloads)

**Resource Access:**

- Key Vault access logs (secret retrieval, who accessed what)
- Storage Account access logs (blob access patterns, unauthorized attempts)
- Database audit logs (authentication failures, privilege escalation attempts, schema changes)

**Infrastructure Security:**

- Activity Logs (resource deletions, configuration changes, who did what)
- Resource Manager operations (deployment failures, policy violations)
- VM authentication logs (SSH/RDP access attempts, privilege escalations)

**Who Uses It:**

- Security Operations Center (SOC)
- Security Engineers
- Compliance/Audit teams
- Incident Response teams

**Query Patterns:**

- "Is there suspicious login activity?"
- "Did someone access production secrets?"
- "Are there SQL injection attempts?"
- "Who modified this firewall rule?"
- "Is there lateral movement between environments?"

---

### Microsoft Sentinel Configuration

#### Dual Sentinel Architecture

We run **two separate Sentinel instances** to maintain complete geographic data isolation:

```
┌─────────────────────────────────────────────────────┐
│         Sentinel Instance - Global (East US 2)      │
├─────────────────────────────────────────────────────┤
│                                                      │
│ Data Sources:                                       │
│ ├─ Global Security LAW (primary)                    │
│ ├─ Global Management LAW (context)                  │
│ └─ Monitors: Global pair infrastructure             │
│                                                      │
│ Analytics Rules:                                     │
│ ├─ Brute force detection (multiple failed logins)  │
│ ├─ Privilege escalation detection                   │
│ ├─ Lateral movement (Dev → Prod access attempts)   │
│ ├─ Data exfiltration (unusual outbound traffic)    │
│ ├─ Anomalous database access patterns               │
│ ├─ Suspicious Key Vault access                      │
│ ├─ Firewall rule changes                            │
│ └─ After-hours admin activity                       │
│                                                      │
│ Threat Intelligence:                                │
│ ├─ Microsoft Threat Intelligence feed               │
│ ├─ Known malicious IP addresses                     │
│ └─ Indicators of Compromise (IoCs)                  │
│                                                      │
│ Automation & Response:                              │
│ ├─ Playbooks (Azure Logic Apps):                    │
│ │  ├─ Auto-block malicious IPs in Firewall         │
│ │  ├─ Disable compromised user accounts             │
│ │  ├─ Isolate affected VMs                          │
│ │  └─ Create incident tickets (ServiceNow/Jira)    │
│ └─ Notification channels (Slack, PagerDuty, Email) │
│                                                      │
│ Workbooks (Dashboards):                             │
│ ├─ Security Overview dashboard                      │
│ ├─ Firewall traffic analysis                        │
│ ├─ Identity & Access dashboard                      │
│ └─ Compliance reporting                             │
│                                                      │
│ Access: Global Security team                         │
└─────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────┐
│        Sentinel Instance - GDPR (West Europe)       │
├─────────────────────────────────────────────────────┤
│                                                      │
│ Data Sources:                                       │
│ ├─ GDPR Security LAW (primary)                      │
│ ├─ GDPR Management LAW (context)                    │
│ └─ Monitors: GDPR pair infrastructure               │
│                                                      │
│ Same analytics rules, playbooks, and workbooks      │
│ Access: GDPR Security team (can be same people)     │
│ Data residency: EU only ✅                          │
│                                                      │
└─────────────────────────────────────────────────────┘
```

---

### Why Separate Sentinel Instances?

**GDPR Compliance:**

- EU security data (even metadata about infrastructure) stays in EU
- No EU security logs analyzed by US-based infrastructure
- Clear data processing boundaries for DPO (Data Protection Officer)

**Operational Independence:**

- Each geography can respond to incidents independently
- EU team can investigate EU incidents without US infrastructure dependency
- Maintains regional operational resilience

**Consistency:**

- Mirrors the geographic split in LAW, Storage, and Application data
- Simpler architecture to explain to auditors
- No exceptions to "EU data stays in EU" rule

**Cross-Environment Monitoring (Within Geography):**

- Each Sentinel instance monitors BOTH Prod and NonProd in its geography
- Detects lateral movement Dev → Prod within same geography
- Example: Credential theft in Global NonProd → use in Global Prod ✅
- Does NOT detect: Global → GDPR attacks (different customer bases anyway)

---

## Unified Security View (Optional)

For organizations with a centralized Global SOC team that needs visibility across both geographies:

**Option: Microsoft Defender XDR**

- Provides unified view across multiple Sentinel instances
- Security team can see both Global and GDPR incidents in one portal
- Data still remains in respective regions (Defender XDR just provides UI layer)
- No data mixing - just visualization

**Alternative: Dual Dashboard Access**

- Security analysts have access to both Sentinel instances
- Toggle between Global and GDPR dashboards
- Most incidents are region-specific anyway (different customer bases)

---

### Log Ingestion: No Cross-Region Bandwidth Charges

**Critical Cost Factor:**

```
All Diagnostic Settings → LAW/Storage are FREE of bandwidth charges:

✅ Central US → East US 2 LAW: FREE
✅ North Europe → West Europe LAW: FREE
✅ East US 2 → East US 2 Storage: FREE
✅ Central US → East US 2 Storage: FREE
✅ West Europe → West Europe Storage: FREE
✅ North Europe → West Europe Storage: FREE

Even inter-continental would be free (but we don't do this for GDPR):
✅ North Europe → East US 2 LAW: FREE (but we avoid for compliance)
```

**From Microsoft Documentation:**
> "Log data ingestion is free. Bandwidth charges don't apply for ingestion of data to Log Analytics or Storage Accounts from Azure resources via Diagnostic Settings, regardless of source and destination region."

---

### Long-Term Log Archival Strategy

#### Storage Account Architecture

**2 Storage Accounts (Geographic Separation):**

```
┌─────────────────────────────────────────────────────┐
│        Log Archive Storage - Global (East US 2)     │
├─────────────────────────────────────────────────────┤
│                                                      │
│ Type: StorageV2 (General Purpose v2)               │
│ Replication: GRS (automatic to Central US)         │
│ Default Access Tier: Cool                           │
│ Network: Private endpoint only                      │
│                                                      │
│ Receives logs from (via Diagnostic Settings):      │
│ ├─ East US 2: All resources (Prod + NonProd)       │
│ └─ Central US: All resources (Prod + NonProd)       │
│    └─ No bandwidth charge ✅                        │
│                                                      │
│ Lifecycle Policy:                                   │
│ ├─ Days 0-30: Cool tier (~$0.01/GB/month)         │
│ ├─ Days 31-90: Cold tier (~$0.004/GB/month)       │
│ ├─ Days 91-365: Archive tier (~$0.002/GB/month)   │
│ └─ Days 366+: Auto-delete                          │
│                                                      │
│ Containers:                                         │
│ ├─ insights-logs-azurefirewallnetworkrule/         │
│ ├─ insights-logs-azurefirewallapplicationrule/     │
│ ├─ insights-logs-nsgflowlogs/                      │
│ ├─ insights-logs-azureactivity/                    │
│ └─ ... (one per log category)                      │
│                                                      │
└─────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────┐
│       Log Archive Storage - GDPR (West Europe)      │
├─────────────────────────────────────────────────────┤
│                                                      │
│ Same configuration as Global archive                │
│                                                      │
│ Receives logs from (via Diagnostic Settings):      │
│ ├─ West Europe: All resources (Prod only)          │
│ └─ North Europe: All resources (Prod only)          │
│    └─ No bandwidth charge ✅                        │
│                                                      │
│ Data residency: EU only ✅                          │
│                                                      │
└─────────────────────────────────────────────────────┘
```

---

#### Complete Log Retention Timeline

```
┌─────────────────────────────────────────────────────┐
│              COMPLETE LOG LIFECYCLE                 │
├─────────────────────────────────────────────────────┤
│                                                      │
│ DAY 0-30 (Hot - Active Querying):                  │
│ ├─ Primary: Log Analytics Workspace                │
│ │  ├─ Purpose: Dashboards, alerts, investigations  │
│ │  ├─ Query Speed: Very fast (indexed)             │
│ │  └─ Use Cases: Daily ops, incident response      │
│ └─ Backup: Storage Account (Cool tier)             │
│    ├─ Purpose: Compliance backup                    │
│    └─ Query Speed: Medium (if needed)              │
│                                                      │
│ DAY 31-90 (Warm - Occasional Access):              │
│ ├─ Location: Storage Account (Cold tier)           │
│ ├─ Purpose: Compliance, historical analysis        │
│ ├─ Query Speed: Slow (need to rehydrate)           │
│ └─ Cost: 60% cheaper than Cool tier                │
│                                                      │
│ DAY 91-365 (Cold - Compliance Archive):            │
│ ├─ Location: Storage Account (Archive tier)        │
│ ├─ Purpose: Regulatory compliance (SOC 2, ISO)     │
│ ├─ Query Speed: Very slow (hours to rehydrate)     │
│ ├─ Cost: 80% cheaper than Cool tier                │
│ └─ Access: Rarely (forensics, audits)              │
│                                                      │
│ DAY 366+ (Retention Complete):                      │
│ └─ Automatically deleted (unless extended for PCI)  │
│                                                      │
└─────────────────────────────────────────────────────┘
```

**Accessing Archived Logs:**

When you need logs older than 30 days:

1. **Recent (31-90 days - Cold tier):**
   - Rehydrate to Cool/Hot tier (takes minutes)
   - Query via Azure Storage Explorer or Azure Portal
   - Cost: Standard Cold tier retrieval charges

2. **Old (91-365 days - Archive tier):**
   - Rehydrate to Cool/Hot tier (takes hours)
   - Set rehydration priority: Standard (12-24 hours) or High (1 hour, 10x cost)
   - Query after rehydration completes
   - Cost: Archive retrieval charges + temporary Hot storage

3. **Compliance Exports:**
   - For audits: Export specific date ranges to Hot tier
   - Provide to auditors as CSV/JSON
   - Delete temporary copies after audit completes

---

## Database Monitoring & Security Logs

### Azure-Managed Database Services

**Azure Database for PostgreSQL Flexible Server:**

```
Diagnostic Settings:
├─ To Management LAW:
│  ├─ PostgreSQL Server Logs (slow queries, deadlocks)
│  ├─ Performance metrics (CPU, memory, connections)
│  └─ Replication lag metrics
│
└─ To Security LAW:
   ├─ Audit logs (all DDL/DCL statements)
   ├─ Authentication failures
   └─ Privilege changes
```

**Azure Managed Instance for Apache Cassandra:**

```
Diagnostic Settings:
├─ To Management LAW:
│  ├─ Cluster health metrics
│  ├─ Node performance (CPU, disk, network)
│  └─ Query performance statistics
│
└─ To Security LAW:
   ├─ Authentication events
   └─ Configuration changes
```

**Azure Cache for Redis:**

```
Diagnostic Settings:
├─ To Management LAW:
│  ├─ Performance metrics (cache hits, evictions)
│  ├─ Memory usage
│  └─ Connection counts
│
└─ To Security LAW:
   └─ Access logs (who connected, when)
```

---

### Self-Hosted Database Services

**Elasticsearch, TimescaleDB, MongoDB, Kafka:**

```
Log Collection via Azure Monitor Agent:
├─ Application logs → rsyslog → Azure Monitor Agent → LAW
│
├─ To Management LAW:
│  ├─ Service health (cluster status, node health)
│  ├─ Performance metrics (query latency, throughput)
│  ├─ Resource utilization (CPU, memory, disk)
│  └─ Replication status
│
└─ To Security LAW:
   ├─ Authentication logs (SSH, database auth)
   ├─ Authorization failures
   ├─ Configuration file changes
   └─ Suspicious query patterns
```

**Custom Metrics Exporter:**

```
Each self-hosted service:
├─ Metrics endpoint (Prometheus format)
├─ Azure Monitor Agent scrapes metrics
├─ Sends to appropriate LAW (Management or Security)
└─ Custom dashboards in Azure Monitor
```

---

## Access Control Strategy

### Role-Based Access (RBAC) on Log Analytics

**Global Management LAW:**

```
├─ Log Analytics Reader:
│  ├─ Platform Engineering team (all data)
│  ├─ SRE team (all data)
│  └─ Database team (all data)
│
├─ Log Analytics Contributor:
│  └─ Platform Lead (for workspace config)
│
└─ Resource-Level RBAC:
   ├─ Application Team A: Reader (filtered to their resources only)
   └─ Application Team B: Reader (filtered to their resources only)
```

**Global Security LAW:**

```
├─ Log Analytics Reader:
│  ├─ Security Operations team (all data)
│  ├─ Compliance team (read-only)
│  └─ Incident Response team (all data)
│
├─ Microsoft Sentinel Responder:
│  └─ SOC analysts (can manage incidents, run playbooks)
│
└─ Microsoft Sentinel Contributor:
   └─ Security Engineering (can modify analytics rules)
```

**GDPR LAWs:**

```
Same RBAC structure as Global
Access limited to EU-focused teams (or same global teams)
Data residency maintained through RBAC
```

---

## Monitoring the Monitoring

**Health Monitoring for Observability Infrastructure:**

```
Monitor the monitoring systems themselves:

LAW Health:
├─ Daily data ingestion volume (alerts if drops >20%)
├─ Query performance (alerts if >5 seconds)
├─ Workspace health status
└─ Data export failures

Sentinel Health:
├─ Analytics rule execution status
├─ Playbook run failures
├─ Data connector health
└─ Incident response SLA compliance

Storage Account Health:
├─ Lifecycle policy execution
├─ Replication status (GRS lag)
├─ Access failures (auth issues)
└─ Capacity warnings (approaching limits)

Diagnostic Settings:
├─ Failed log delivery
├─ Throttling events
└─ Configuration drift detection
```

---

## Cost Management

**Observability Infrastructure Costs:**

```
Log Analytics:
├─ Ingestion: Pay-per-GB (with commitment tiers for savings)
├─ Retention: First 30 days free, then $0.10/GB/month
├─ Queries: Pay for data scanned
└─ Optimization: Filter logs at source, use Basic Logs tier for high-volume

Microsoft Sentinel:
├─ Ingestion: Same as LAW (share cost)
├─ Analytics: Additional cost per GB ingested
├─ Playbook runs: Logic App execution costs
└─ Optimization: Use built-in rules, efficient queries

Storage Accounts:
├─ Cool tier: ~$0.01/GB/month
├─ Cold tier: ~$0.004/GB/month
├─ Archive tier: ~$0.002/GB/month
├─ Retrieval: Pay when accessing Cold/Archive
└─ Optimization: Aggressive lifecycle policies, compress before storing

Estimated Monthly Costs (example for 50GB/day ingestion):
├─ LAW ingestion: ~$3,000/month
├─ Sentinel: ~$2,000/month
├─ Storage (1 year): ~$150/month
└─ Total: ~$5,150/month for complete observability
```

---

## Summary: Complete Observability Architecture

**Total Components:**

- **4 Log Analytics Workspaces:** 2 Management (Global + GDPR), 2 Security (Global + GDPR)
- **2 Sentinel Instances:** Global + GDPR
- **2 Storage Accounts:** Global Archive + GDPR Archive
- **0 Bandwidth Charges:** All diagnostic settings ingestion is free

**Key Principles:**

- ✅ Geographic data isolation for GDPR compliance
- ✅ Separate operational and security logs for clarity
- ✅ 30-day hot querying in LAW, long-term archive in storage
- ✅ Cost-optimized lifecycle policies (Cool → Cold → Archive → Delete)
- ✅ Dual Sentinel for independent regional security monitoring
- ✅ Consistent architecture across all observability layers

**Data Flow:**

```
Azure Resource (Firewall, VM, Database, etc.)
    ↓
Diagnostic Settings (configured once)
    ├─ Destination 1: Log Analytics Workspace (30 days)
    │   └─ Microsoft Sentinel (real-time threat detection)
    └─ Destination 2: Storage Account (1 year, lifecycle managed)
        ├─ Days 0-30: Cool tier
        ├─ Days 31-90: Cold tier
        ├─ Days 91-365: Archive tier
        └─ Day 366+: Auto-delete
```

---

**Does this updated section address your concerns? The key changes:**

1. ✅ Split Management LAW by geography (Global + GDPR)
2. ✅ Split Sentinel by geography (Global + GDPR)
3. ✅ Explained why Sentinel should also be split (consistency with storage account split)
4. ✅ Clarified what logs go to Management LAW vs Security LAW
5. ✅ Confirmed no bandwidth charges for diagnostic settings → LAW/Storage (even cross-region)
6. ✅ Total: 4 LAW + 2 Sentinel + 2 Storage Accounts

### Log Analytics Workspace Strategy

**4 Total Workspaces (Environment + Geography Scoped):**

```
Management Logs (Operational Telemetry):
├─ Prod Management LAW (East US 2)
│  └─ Ingests: All Production hub logs (all 4 regions)
│     ├─ Firewall logs
│     ├─ VM metrics
│     ├─ NSG flow logs
│     ├─ Database metrics (PostgreSQL, Cassandra, Redis)
│     └─ Self-hosted service logs (Elasticsearch, MongoDB, Kafka)
│
└─ NonProd Management LAW (East US 2)
   └─ Ingests: All Non-Production hub logs
      └─ Same log types as Prod

Security Logs (Threat Detection):
├─ Centralized Security LAW (East US 2) ✅
│  ├─ Ingests: Security events from ALL environments
│  │  ├─ Azure AD logs
│  │  ├─ Defender for Cloud alerts
│  │  ├─ Threat intelligence feeds
│  │  ├─ Database audit logs
│  │  └─ Security events (Prod + NonProd)
│  │
│  └─ Microsoft Sentinel enabled here ✅
│     └─ Monitors: All environments for cross-environment attacks
│
└─ GDPR Security LAW (West Europe) - Optional
   └─ If GDPR requires EU security logs remain in EU
```

**Why This Structure:**

**Separate Management Logs (Prod vs NonProd):**

- ✅ Different access controls (Prod team sees only Prod logs)
- ✅ Cost tracking per environment
- ✅ Different retention policies (Prod: 90 days, NonProd: 30 days)

**Centralized Security Logs:**

- ✅ Security team needs unified view across ALL environments
- ✅ Attack patterns often span Prod ↔ NonProd (creds stolen in Dev, used in Prod)
- ✅ Sentinel correlation rules work across environments
- ✅ Industry best practice: Centralized SIEM

**Cross-Region Ingestion Cost:**

- **FREE!** Microsoft does NOT charge for cross-region log ingestion to Log Analytics
- Central US → East US 2 LAW: No bandwidth charges
- North Europe → West Europe LAW: No bandwidth charges

### Microsoft Sentinel Configuration

**1 Primary Sentinel Instance (Shared Across Environments):**

```
Microsoft Sentinel (on Centralized Security LAW):
├─ Monitors: ALL environments (Prod + NonProd, all 4 regions)
├─ Analytics Rules:
│  ├─ Detect lateral movement Dev → Prod
│  ├─ Detect anomalous database access patterns
│  ├─ Detect credential theft across environments
│  ├─ Detect data exfiltration attempts
│  └─ Tested in NonProd before enabling in Prod ✅
│
├─ Data Connectors:
│  ├─ Azure AD
│  ├─ Azure Firewall
│  ├─ PostgreSQL audit logs
│  ├─ Cassandra audit logs
│  ├─ VM syslog (Elasticsearch, MongoDB, Kafka hosts)
│  └─ Custom log ingestion (application logs)
│
├─ Automation & Response:
│  ├─ Playbooks for incident response
│  ├─ Automatic threat remediation
│  └─ Integration with ticketing systems
│
└─ Unified Incident View:
   └─ Security analyst sees complete attack chain
```

**Why Centralized Sentinel (Not Separate Per Environment):**

- ✅ Real attack pattern: Compromise NonProd → Pivot to Prod
- ✅ Sentinel correlation requires seeing all data
- ✅ Security Operations Center (SOC) needs single pane of glass
- ✅ SOC 2 / ISO 27001 REQUIRE centralized security monitoring
- ✅ Industry standard: Even Fortune 500 run 1-2 Sentinel instances globally

**Optional: 2nd Sentinel for GDPR (if required):**

- If GDPR regulations require EU security logs stay in EU
- Deploy 2nd Sentinel in West Europe
- Primarily for compliance, not operational need

### Database Monitoring & Performance

**Azure-Managed Services:**

- Azure Monitor integration (automatic)
- Query Performance Insights
- Intelligent Performance recommendations
- Automated alerts on resource utilization

**Self-Hosted Services Monitoring:**

```
Elasticsearch:
├─ Elasticsearch monitoring API
├─ Metrics exported to Azure Monitor
├─ Cluster health dashboards
└─ Index performance tracking

TimescaleDB:
├─ pg_stat_statements for query analysis
├─ TimescaleDB telemetry
├─ Metrics to Azure Monitor via custom agent
└─ Replication lag monitoring

MongoDB:
├─ MongoDB Ops Manager or Cloud Manager
├─ Metrics exported to Azure Monitor
├─ Replica set health monitoring
└─ Sharding performance tracking

Kafka:
├─ JMX metrics collection
├─ Kafka Manager UI
├─ Metrics exported to Azure Monitor
└─ Consumer lag monitoring
```

---

## Backup & Disaster Recovery

### Recovery Services Vault Strategy

**2 Vaults (Environment Scoped):**

```
Production RSV (East US 2):
├─ Backs up: Production VMs in all 4 regions
│  ├─ Self-hosted database VMs (Elasticsearch, MongoDB, etc.)
│  ├─ Kafka broker VMs
│  ├─ Application servers
│  └─ Jump boxes
├─ Replication: Azure Site Recovery to Central US
├─ Storage: GRS (automatic geo-replication to Central US)
├─ Retention: 90 days
├─ Features:
│  ├─ Immutable backups (ransomware protection)
│  ├─ Multi-User Authorization (MUA)
│  └─ Resource Guard (requires Security Officer approval)

Non-Production RSV (East US 2):
├─ Backs up: NonProd VMs in East US 2 + Central US
├─ Storage: LRS (cost optimization)
├─ Retention: 7-14 days
├─ Features: Standard (no MUA/immutability needed)
```

**Why Separate Vaults:**

- ✅ Different retention policies (Prod: long, NonProd: short)
- ✅ Different recovery priorities (Prod first, NonProd can wait)
- ✅ Ransomware air-gap (compromised NonProd admin can't delete Prod backups)
- ✅ Different cost profiles (Prod: GRS, NonProd: LRS)
- ✅ Multi-User Authorization only on Prod (Security Officer approval required)

### Backup Vault Strategy (PaaS Services)

**2 Vaults (Environment Scoped):**

```
Production Backup Vault (East US 2):
├─ Backs up:
│  ├─ Azure Storage Accounts (blob backup)
│  ├─ Azure PostgreSQL Flexible Server (incremental)
│  ├─ Azure Managed Cassandra (automated)
│  └─ Azure Cache for Redis (RDB snapshots)
├─ Storage: GRS
├─ Immutable: Enabled
└─ Retention: 30-90 days depending on service

Non-Production Backup Vault (East US 2):
├─ Backs up: NonProd PaaS resources (if any)
├─ Storage: LRS
└─ Retention: 7 days
```

### Application-Level Backup Strategy

**Self-Hosted Databases:**

```
Elasticsearch:
├─ Automated snapshots to Azure Storage Account (GZRS)
├─ Incremental snapshots (efficient storage)
├─ Retention: 30 days
└─ Restore tested quarterly

TimescaleDB:
├─ Continuous WAL archiving to Azure Storage (GZRS)
├─ Base backups via pg_basebackup
├─ PITR capability up to 30 days
└─ Automated via cron + Azure CLI

MongoDB:
├─ mongodump to Azure Storage Account (GZRS)
├─ Oplog tailing for PITR
├─ Retention: 30 days
└─ Automated via cron + azcopy

Kafka:
├─ Topic data is transient (7-day retention in Kafka)
├─ Long-term data archived to Azure Storage via Kafka Connect
├─ Schema Registry backed up to Azure Storage
└─ Consumer offsets preserved during failover
```

### Disaster Recovery Testing

**Annual DR Exercise (Comprehensive):**

- Test failover of entire Production environment
- Validate RTO/RPO meets requirements
- Test database failover procedures
- Verify self-hosted service replication
- Document lessons learned

**Quarterly DR Exercise (Non-Production):**

- Deploy Central US NonProd hub (on-demand)
- Test failover procedures in NonProd environment
- Test self-hosted database failover (Elasticsearch, MongoDB, Kafka)
- Validate runbooks and automation
- Train teams on DR procedures
- Tear down after testing (cost optimization)

**Monthly Tabletop Exercise:**

- Walk through DR procedures
- Update runbooks based on infrastructure changes
- Review database backup/restore procedures
- Ensure team familiarity

### RTO/RPO Targets by Service

| Service | RPO Target | RTO Target | Failover Method |
|---------|-----------|-----------|-----------------|
| **PostgreSQL Flexible** | <1 minute | <5 minutes | Promote read replica |
| **Storage Account (GZRS)** | ~15 minutes | Variable | Microsoft-managed failover |
| **Cassandra** | <5 minutes | <10 minutes | Client failover to standby DC |
| **Redis** | <1 minute | <5 minutes | Active geo-replication |
| **Elasticsearch** | <5 minutes | ~15 minutes | Switch DNS to standby cluster |
| **TimescaleDB** | <5 minutes | ~15 minutes | Promote streaming replica |
| **MongoDB** | <10 seconds | <5 minutes | Replica set election |
| **Kafka** | <30 seconds | <10 minutes | MirrorMaker active, client reconnect |

---

## DevOps & Infrastructure Management

### DevOps Infrastructure Topology

**Global Pair DevOps (East US 2 + Central US):**

```
East US 2:
├─ GitHub Actions self-hosted runners
├─ Azure Container Registry (ACR)
├─ HashiCorp Vault (secrets management)
├─ Terraform/OpenTofu state storage (Azure Storage)
├─ Artifact repositories
└─ CI/CD pipeline infrastructure

Central US:
└─ DR runners (on-demand or 10% capacity)
```

**GDPR Pair DevOps (West Europe + North Europe):**

```
West Europe:
├─ GitHub Actions self-hosted runners (GDPR deployments)
├─ Azure Container Registry (geo-replicated)
├─ HashiCorp Vault (separate cluster or replicated)
├─ Terraform/OpenTofu state storage
└─ GDPR-specific tooling

North Europe:
└─ DR infrastructure (minimal or on-demand)
```

**Why Duplicate DevOps (Not Connect Pairs):**

- ✅ GDPR isolation (no network path between US and EU)
- ✅ Each geography self-sufficient
- ✅ Better latency for EU deployments
- ✅ Cleaner audit story (no US → EU connectivity)
- ✅ Cost is reasonable (DevOps infrastructure relative to total spend)

### Infrastructure-as-Code Strategy

**Repository Structure:**

```
/infrastructure
├─ modules/
│  ├─ hub/                    # Reusable hub module
│  ├─ spoke/                  # Reusable spoke module
│  ├─ observability/          # LAW, Sentinel, monitoring
│  ├─ postgres-flexible/      # PostgreSQL module
│  ├─ cassandra/              # Cassandra module
│  ├─ self-hosted-db/         # VM-based databases
│  └─ kafka-cluster/          # Kafka cluster module
│
├─ global-pair/
│  ├─ eastus2-prod-hub/
│  ├─ eastus2-nonprod-hub/
│  ├─ eastus2-data-services/  # DB deployments
│  ├─ centralus-prod-hub/
│  └─ centralus-data-services/
│
├─ gdpr-pair/
│  ├─ westeurope-prod-hub/
│  ├─ westeurope-data-services/
│  └─ ...
│
└─ global-resources/
   ├─ azure-front-door/
   ├─ dns/
   └─ key-vault/
```

**GitHub Actions Workflow:**

```
1. PR opened → Terraform plan → Review by CODEOWNERS
2. Infrastructure validation:
   ├─ Terraform validate
   ├─ tflint (linting)
   ├─ checkov (security scanning)
   └─ Cost estimation (Infracost)
3. PR approved → Merge to main
4. Terraform apply → Deploy to affected region/environment
5. Post-deployment validation:
   ├─ Connectivity tests
   ├─ Database replication checks
   └─ Health checks
6. Slack notification → Deployment complete
```

**CODEOWNERS Configuration:**

```
# Production infrastructure requires Security team approval
/global-pair/*prod-hub/          @security-team @platform-team
/gdpr-pair/*prod-hub/            @security-team @platform-team
/global-pair/*data-services/     @data-team @security-team

# NonProd infrastructure can be approved by Platform team
/global-pair/*nonprod-hub/       @platform-team
/global-pair/*nonprod-data/      @data-team @platform-team
```

### Subscription Strategy

**Azure Landing Zone Subscriptions:**

```
Management Group Hierarchy:
├─ Platform
│  ├─ Management (1 subscription - global)
│  │  └─ Log Analytics, monitoring, automation accounts
│  │
│  ├─ Identity (1 subscription - global)
│  │  └─ Azure AD, Domain Controllers (if needed)
│  │
│  └─ Connectivity
│     ├─ Connectivity-Prod (1 subscription)
│     │  └─ All 4 Production hubs
│     └─ Connectivity-NonProd (1 subscription)
│        └─ 2 Non-Production hubs
│
└─ Landing Zones
   ├─ Prod Workloads
   │  ├─ Prod-App (1+ subscriptions)
   │  │  └─ Application spokes in all 4 regions
   │  ├─ Prod-Data (1+ subscriptions)
   │  │  ├─ Azure-managed data services
   │  │  └─ Self-hosted database VMs
   │  └─ RBAC: Production team only
   │
   └─ NonProd Workloads
      ├─ NonProd-App (1+ subscriptions)
      │  └─ Dev, QA, Staging spokes
      ├─ NonProd-Data (1+ subscriptions)
      │  └─ NonProd databases and test data
      └─ RBAC: Dev + Platform teams
```

**Key Principles:**

- Subscriptions are environment-scoped (Prod vs NonProd), not region-scoped
- Data services may warrant separate subscriptions for cost tracking
- Subscription limits are high (unlikely to hit for most deployments)
- Use resource groups + tags for detailed cost allocation

---

## Cost Optimization Strategies

### Hub Infrastructure Optimization

**Sizing by Role:**

- Active Prod hubs: Premium/Standard SKUs (full features)
- Active NonProd hubs: Standard/Basic SKUs (cost-optimized)
- Standby Prod hubs: Standard SKUs (handle 10-30% traffic)
- Standby NonProd hubs: On-demand only (deploy 2× per year)

**Shared Components:**

- ExpressRoute circuits shared between Prod and NonProd
- LNG (Local Network Gateway) shared
- Security LAW + Sentinel shared across environments

**Estimated Savings:**

- On-demand Central US NonProd hub: Save ~85% vs always-on
- Shared ER circuits: Save ~50% vs dedicated per environment
- Standby hubs without VPN/ER: Save ~25% per standby hub

### Data Services Optimization

**Azure-Managed Services:**

- Use Reserved Instances for production databases (30-50% discount)
- Right-size database SKUs based on actual usage
- Use Zone-Redundant HA in active regions only (standby uses replicas)
- Enable auto-pause for non-production databases (if applicable)

**Self-Hosted Services:**

- Use B-series or D-series VMs for cost efficiency
- Implement auto-scaling for non-production environments
- Use Premium SSD only where IOPS requirements justify it
- Shut down non-production VMs outside business hours (potential 60% savings)

**Storage Optimization:**

- Production: GZRS for critical data, LRS for logs/temp data
- Non-Production: LRS only (no geo-redundancy needed)
- Implement lifecycle policies (move old data to Cool/Archive tiers)
- Use managed disks for VMs, but not oversized

### Backup Optimization

**Retention Policies:**

- Production: 30-90 days (regulatory requirements)
- Non-Production: 7-14 days (sufficient for dev/test)
- Self-hosted DBs: Incremental backups to minimize storage

**Storage Tier Usage:**

- Recent backups (0-7 days): Hot tier
- Older backups (8-30 days): Cool tier
- Compliance backups (31+ days): Archive tier

### Monitoring & Observability Optimization

**Log Retention:**

- Security logs: 90-180 days (compliance requirement)
- Operational logs Prod: 90 days
- Operational logs NonProd: 30 days
- Archive old logs to Azure Storage (Cool/Archive tier)

**Log Ingestion:**

- Filter noisy logs at source (reduce ingestion volume)
- Use Basic Log Analytics for high-volume, low-value logs
- Sample logs where appropriate (e.g., sample 10% of debug logs)

### Network Optimization

**VNet Peering:**

- Global peering has cost ($0.035/GB)
- Minimize cross-region data transfer where possible
- Use spoke-to-spoke peering for database replication (most efficient)

**Egress Optimization:**

- All 4 regions in Zone 1 ($0.087/GB egress)
- Use Azure Front Door for edge caching (reduce origin egress)
- Implement CDN for static content
- Compress data before egress where applicable

---

## Decision Audit Trail

### Hub Architecture Decision

**Decision: Separate Hubs per Environment**

| Factor | Weight | Separate Hubs | Shared Hub | Winner |
|--------|--------|---------------|------------|--------|
| Security Isolation | 25% | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | Separate |
| Testing Independence | 20% | ⭐⭐⭐⭐⭐ | ⭐⭐ | Separate |
| Compliance/Audit | 20% | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | Separate |
| Operational Complexity | 15% | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | Shared |
| Cost | 20% | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | Shared |

**Rationale:**

- GitOps maturity eliminates operational complexity concerns
- SOC 2 + ISO 27001 + PCI-DSS compliance benefits justify cost
- True infrastructure testing capability is critical
- Network-level isolation is strongest security posture

### Region Pair Selection

| Factor | Weight | Impact on Decision |
|--------|--------|-------------------|
| GDPR Compliance | 30% | Mandated separate EU cluster |
| Availability Zone Maturity | 25% | All 4 regions must have 3 AZs now |
| Operational Simplicity | 20% | Native GRS > manual replication |
| Global Latency | 15% | East US 2 + West Europe optimal |
| Cost | 10% | Considered but not primary driver |

### Data Services: Managed vs Self-Hosted

| Service | Decision | Rationale |
|---------|----------|-----------|
| **PostgreSQL** | Azure Flexible Server | Built-in HA, PITR, automatic backups, geo-replication |
| **Cassandra** | Azure Managed | Multi-region built-in, reduced ops overhead |
| **Redis** | Azure Cache | Enterprise features, active geo-replication, managed |
| **Elasticsearch** | Self-Hosted | Advanced features, custom plugins, cost at scale |
| **TimescaleDB** | Self-Hosted | Not available as managed service, specific requirements |
| **MongoDB** | Self-Hosted | Full control over sharding, specific feature requirements |
| **Kafka** | Self-Hosted | Kafka-specific features, Connect ecosystem, cost |

### Cross-Region Connectivity Decision

**Decision: Do NOT Connect Global ↔ GDPR Pairs**

**Rationale:**

- ✅ GDPR isolation (clear regulatory boundary)
- ✅ Simpler audit story (no US → EU data paths)
- ✅ Each geography self-sufficient
- ✅ Duplicate DevOps is reasonable cost
- ❌ Would require complex firewall rules to prevent data mixing
- ❌ Would require proving to auditors that customer data can't cross

---

## Conclusion

Our **4-region, 6-hub architecture** provides:

✅ **Regulatory Compliance** - GDPR-compliant EU data residency
✅ **High Availability** - 3 Availability Zones in all 4 regions
✅ **Security Isolation** - Separate Prod/NonProd hubs with network-level separation
✅ **Disaster Recovery** - Native GRS/GZRS for managed services, spoke-to-spoke replication for self-hosted
✅ **Global Performance** - <150ms latency for 90% of users
✅ **Testing Independence** - True production-parity infrastructure testing
✅ **Operational Simplicity** - GitOps-managed, repeatable deployments
✅ **Cost Efficiency** - Optimized within compliance constraints
✅ **Data Flexibility** - Hybrid approach balancing managed services with self-hosted requirements

**This configuration represents industry best practices for global, compliant, highly-available enterprise SaaS deployments.**

---

## References

- [Azure Region Pairs Documentation](https://learn.microsoft.com/en-us/azure/reliability/regions-paired)
- [GDPR Article 44-50: Data Transfers](https://gdpr-info.eu/chapter-5/)
- [Azure Storage Redundancy Options](https://learn.microsoft.com/en-us/azure/storage/common/storage-redundancy)
- [Azure Front Door Documentation](https://learn.microsoft.com/en-us/azure/frontdoor/)
- [Azure Database for PostgreSQL Flexible Server](https://learn.microsoft.com/en-us/azure/postgresql/flexible-server/)
- [Azure Managed Instance for Apache Cassandra](https://learn.microsoft.com/en-us/azure/managed-instance-apache-cassandra/)
- [Azure Cache for Redis](https://learn.microsoft.com/en-us/azure/azure-cache-for-redis/)
