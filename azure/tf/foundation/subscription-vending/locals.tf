# =============================================================================
# Azure default RegistrationFree resource providers that are auto registered:
# =============================================================================
# "Microsoft.ADHybridHealthService", # AD Hybrid Health monitoring
# "Microsoft.Authorization",         # RBAC and policy
# "Microsoft.Billing",               # Billing information
# "Microsoft.ChangeSafety",          # Change safety controls
# "Microsoft.ClassicSubscription",   # Classic subscription management
# "Microsoft.Commerce",              # Commerce and marketplace
# "Microsoft.Consumption",           # Cost management and billing
# "Microsoft.CostManagement",        # Cost analysis and budgets
# "Microsoft.Features",              # Feature flags and previews
# "Microsoft.MarketplaceOrdering",   # Marketplace agreements
# "Microsoft.Portal",                # Azure Portal customizations
# "Microsoft.ResourceGraph",         # Resource queries
# "Microsoft.ResourceIntelligence",  # Resource recommendations
# "Microsoft.ResourceNotifications", # Resource notifications
# "Microsoft.Resources",             # Core resource management
# "Microsoft.SerialConsole",         # VM serial console
# "microsoft.support",               # Support tickets (lowercase in Azure)

# =============================================================================
# Azure default Registration Required resource providers that are auto registered:
# =============================================================================


# =============================================================================
# Local Values - Provider Registration Defaults
# =============================================================================

locals {

  base_resource_providers = [
    "Microsoft.Advisor",            # Azure Advisor recommendations
    "Microsoft.Automation",         # Automation accounts and runbooks
    "Microsoft.GuestConfiguration", # Guest configuration policies
    "microsoft.insights",           # Monitoring and diagnostics
    "Microsoft.ManagedIdentity",    # Managed identities for services
    "Microsoft.PolicyInsights",     # Policy compliance and insights
    "Microsoft.Security",           # Security Center and Defender for Cloud
  ]
}
