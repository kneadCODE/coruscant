# =============================================================================
# Local Values - Provider Registration Defaults
# =============================================================================

locals {
  # Azure default resource providers (RegistrationFree - always registered, cannot be unregistered)
  # Include these in all subscriptions to prevent Terraform state conflicts
  # Source: Azure Portal > Subscription > Resource providers > RegistrationFree
  azure_default_providers = [
    "Microsoft.ADHybridHealthService", # AD Hybrid Health monitoring
    "Microsoft.Authorization",         # RBAC and policy
    "Microsoft.Billing",               # Billing information
    "Microsoft.ChangeSafety",          # Change safety controls
    "Microsoft.ClassicSubscription",   # Classic subscription management
    "Microsoft.Commerce",              # Commerce and marketplace
    "Microsoft.Consumption",           # Cost management and billing
    "Microsoft.CostManagement",        # Cost analysis and budgets
    "Microsoft.Features",              # Feature flags and previews
    "Microsoft.MarketplaceOrdering",   # Marketplace agreements
    "Microsoft.Portal",                # Azure Portal customizations
    "Microsoft.ResourceGraph",         # Resource queries
    "Microsoft.ResourceIntelligence",  # Resource recommendations
    "Microsoft.ResourceNotifications", # Resource notifications
    "Microsoft.Resources",             # Core resource management
    "Microsoft.SerialConsole",         # VM serial console
    "microsoft.support",               # Support tickets (lowercase in Azure)
  ]

  # As recommended by CAF
  base_providers = [
    "Microsoft.ManagedIdentity", # Managed identities for services
    "Microsoft.Insights",        # Monitoring and diagnostics
  ]

  # =============================================================================
  # Subscription-specific provider sets based on actual infrastructure
  # =============================================================================

  # OBSERVABILITY SUBSCRIPTION: Monitoring infrastructure
  # Resources: OTEL collectors, Managed Grafana, Spoke VNet
  observability_providers = [
    "Microsoft.ManagedIdentity", # Managed identities for services
    "Microsoft.Insights",        # Monitoring and diagnostics
    "Microsoft.Compute",         # OTEL collectors (if VM-based)
    "Microsoft.Dashboard",       # Managed Grafana
    "Microsoft.Network",         # Spoke VNet, NSGs, Route Tables
  ]

  # APPLICATION LANDING ZONE SUBSCRIPTIONS (IAM, Payment, etc.)
  # Resources: AKS, PostgreSQL, Redis, Blob Storage, Spoke VNet
  app_landingzone_providers = [
    "Microsoft.ManagedIdentity",  # Managed identities for AKS
    "Microsoft.Insights",         # Monitoring and diagnostics
    "Microsoft.ContainerService", # AKS for application workloads
    "Microsoft.DBforPostgreSQL",  # PostgreSQL databases
    "Microsoft.Cache",            # Redis cache
    "Microsoft.Storage",          # Blob storage accounts
    "Microsoft.Network",          # Spoke VNet, NSGs, Route Tables, Private Endpoints
  ]
}
