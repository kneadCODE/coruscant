package main

import (
	"github.com/kneadCODE/go-diagrams/diagram"
	"github.com/kneadCODE/go-diagrams/nodes/azure"
)

func edgeDiagramAttributes() diagram.Option {
	return diagram.WithAttributes(map[string]string{
		"nodesep":   "1.5",      // Wider horizontal spacing between nodes to prevent label overlap
		"ranksep":   "1.5",      // Vertical spacing between ranks
		"splines":   "polyline", // Polyline edges for better balanced tree layout
		"overlap":   "scale",    // Scale layout to remove overlaps
		"ordering":  "out",      // Preserve edge order for balanced tree layout
		"newrank":   "true",     // Use new ranking algorithm for better layout
	})
}

func genEdge() error {
	d, err := diagram.New(
		diagram.Filename("edge"),
		diagram.Label("ConnEdge"),
		diagram.Direction("TD"),
		edgeDiagramAttributes(),
	)
	if err != nil {
		return err
	}

	sub := azure.General.Subscriptions(diagram.NodeLabel("sub-coruscant-edge-prod"))

	d.Connect(sub, azure.General.Policy(diagram.NodeLabel("Policy")))
	d.Connect(sub, azure.General.CostManagementBilling(diagram.NodeLabel("Management & Billing")))
	d.Connect(sub, azure.General.CostBudgets(diagram.NodeLabel("Budgets")))
	d.Connect(sub, azure.General.CostAlerts(diagram.NodeLabel("Alerts")))

	rgVNETSEA := azure.General.Resourcegroups(diagram.NodeLabel("rg-coruscant-vnet-edge-prod-sea"))
	d.Connect(sub, rgVNETSEA)

	nwSEA := azure.Network.NetworkWatcher(diagram.NodeLabel("nw-coruscant-edge-prod-sea-01"))
	d.Connect(rgVNETSEA, nwSEA)

	vnetSEA := azure.Network.VirtualNetworks(diagram.NodeLabel("vnet-coruscant-edge-prod-sea-01"))
	d.Connect(rgVNETSEA, vnetSEA)

	rgAKS := azure.General.Resourcegroups(diagram.NodeLabel("rg-coruscant-edge-aks-prod-01"))
	d.Connect(sub, rgAKS)
	d.Connect(rgAKS, azure.Identity.ManagedIdentities(diagram.NodeLabel("id-coruscant-edge-aks-prod-01")))
	d.Connect(rgAKS, azure.Compute.KubernetesServices(diagram.NodeLabel("aks-coruscant-edge-prod-01")))
	d.Connect(rgAKS, azure.Compute.Vmss(diagram.NodeLabel("vmss-coruscant-edge-aks-nodepool-system-prod-01")))
	d.Connect(rgAKS, azure.Compute.Vmss(diagram.NodeLabel("vmss-coruscant-edge-aks-nodepool-apps-prod-01")))
	d.Connect(rgAKS, azure.Compute.Disks(diagram.NodeLabel("osdisk")))
	d.Connect(rgAKS, azure.Compute.Disks(diagram.NodeLabel("disk")))
	d.Connect(rgAKS, azure.Network.PrivateEndpoint(diagram.NodeLabel("pep-coruscant-edge-aks-kube-apiserver-prod-01")))
	d.Connect(rgAKS, azure.Network.NetworkInterfaces(diagram.NodeLabel("nic-coruscant-edge-aks-kube-apiserver-prod-01")))
	d.Connect(rgAKS, azure.Network.LoadBalancers(diagram.NodeLabel("lb-coruscant-edge-aks-internal-prod-01")))
	d.Connect(rgAKS, azure.Network.NetworkSecurityGroups(diagram.NodeLabel("nsg-coruscant-edge-aks-nodepool-prod-01")))

	rgAFD := azure.General.Resourcegroups(diagram.NodeLabel("rg-coruscant-edge-afd-prod"))
	d.Connect(sub, rgAFD)

	d.Connect(rgAFD, azure.Network.FrontDoors(diagram.NodeLabel("afd-coruscant-edge-prod-01")))
	d.Connect(rgAFD, azure.Network.PrivateLinkService(diagram.NodeLabel("pls-coruscant-edge-agw-prod-01")))

	return d.Render()
}
