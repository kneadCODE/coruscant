package main

import (
	"github.com/kneadCODE/go-diagrams/diagram"
	"github.com/kneadCODE/go-diagrams/nodes/azure"
)

func genHub() error {
	d, err := diagram.New(
		diagram.Filename("hub"),
		diagram.Label("Hub"),
		diagram.Direction("TD"),
		defaultDiagramAttributes(),
	)
	if err != nil {
		return err
	}

	sub := azure.General.Subscriptions(diagram.NodeLabel("sub-coruscant-hub-prod"))

	d.Connect(sub, azure.General.Policy(diagram.NodeLabel("Policy")))
	d.Connect(sub, azure.General.CostManagementBilling(diagram.NodeLabel("Management & Billing")))
	d.Connect(sub, azure.General.CostBudgets(diagram.NodeLabel("Budgets")))
	d.Connect(sub, azure.General.CostAlerts(diagram.NodeLabel("Alerts")))

	rgEPA := azure.General.Resourcegroups(diagram.NodeLabel("rg-coruscant-hub-epa-prod-01")) // Entra Private Access
	d.Connect(sub, rgEPA)
	d.Connect(rgEPA, azure.Identity.EntraPrivateAccess(diagram.NodeLabel("epa-coruscant-prod-01")))

	rgBas := azure.General.Resourcegroups(diagram.NodeLabel("rg-coruscant-hub-bas-prod-01"))
	d.Connect(sub, rgBas)
	d.Connect(rgBas, azure.Network.Bastion(diagram.NodeLabel("bas-coruscant-prod-01"))) // For break-glass scenarios. Else use EPA

	rgVPN := azure.General.Resourcegroups(diagram.NodeLabel("rg-coruscant-hub-vpn-prod-01"))
	d.Connect(sub, rgVPN)
	d.Connect(rgVPN, azure.Network.VirtualNetworkGateways(diagram.NodeLabel("vgw-coruscant-p2s-prod-01"))) // For break-glass scenarios. Else use EPA

	rgVNETSEA := azure.General.Resourcegroups(diagram.NodeLabel("rg-coruscant-vnet-hub-prod-sea"))
	d.Connect(sub, rgVNETSEA)

	nwSEA := azure.Network.NetworkWatcher(diagram.NodeLabel("nw-coruscant-hub-prod-sea-01"))
	d.Connect(rgVNETSEA, nwSEA)

	vnetSEA := azure.Network.VirtualNetworks(diagram.NodeLabel("vnet-coruscant-hub-prod-sea-01"))
	d.Connect(rgVNETSEA, vnetSEA)

	rgFWSEA := azure.General.Resourcegroups(diagram.NodeLabel("rg-coruscant-hub-afw-prod-sea"))
	d.Connect(sub, rgFWSEA)
	d.Connect(rgFWSEA, azure.Network.Firewall(diagram.NodeLabel("afw-coruscant-hub-prod-sea-01")))
	d.Connect(rgFWSEA, azure.Network.PublicIpAddresses(diagram.NodeLabel("pip-coruscant-hub-afw-prod-sea-01")))

	return d.Render()
}
