package main

import (
	"github.com/kneadCODE/go-diagrams/diagram"
	"github.com/kneadCODE/go-diagrams/nodes/azure"
)

func genGovernance() error {
	d, err := diagram.New(
		diagram.Filename("governance"),
		diagram.Label("Governance"),
		diagram.Direction("TD"),
		defaultDiagramAttributes(),
	)
	if err != nil {
		return err
	}

	// Level 0: Tenant Root
	mgTenantRoot := azure.General.Managementgroups(diagram.NodeLabel("Tenant Root Group"))

	// Level 1: Foundation MG & Platform Root MG
	mgFoundation := azure.General.Managementgroups(diagram.NodeLabel("mg-foundation"))
	mgRoot := azure.General.Managementgroups(diagram.NodeLabel("mg-coruscant-root"))
	d.Connect(mgTenantRoot, mgFoundation)
	d.Connect(mgTenantRoot, mgRoot)

	// Level 2: Root's children (declare left to right for balanced layout)
	mgPlt := azure.General.Managementgroups(diagram.NodeLabel("mg-coruscant-platform"))
	mgLZ := azure.General.Managementgroups(diagram.NodeLabel("mg-coruscant-lz"))
	mgSand := azure.General.Managementgroups(diagram.NodeLabel("mg-coruscant-sandbox"))
	d.Connect(mgRoot, mgPlt)
	d.Connect(mgRoot, mgLZ)
	d.Connect(mgRoot, mgSand)

	// Level 3: Platform's children (declare left to right)
	mgMgmt := azure.General.Managementgroups(diagram.NodeLabel("mg-coruscant-mgmt"))
	mgIdentity := azure.General.Managementgroups(diagram.NodeLabel("mg-coruscant-identity"))
	mgSec := azure.General.Managementgroups(diagram.NodeLabel("mg-coruscant-security"))
	mgConn := azure.General.Managementgroups(diagram.NodeLabel("mg-coruscant-connectivity"))
	d.Connect(mgPlt, mgMgmt)
	d.Connect(mgPlt, mgIdentity)
	d.Connect(mgPlt, mgSec)
	d.Connect(mgPlt, mgConn)

	// Level 3: LZ's children (declare left to right)
	mgCorp := azure.General.Managementgroups(diagram.NodeLabel("mg-coruscant-corp"))
	mgOnline := azure.General.Managementgroups(diagram.NodeLabel("mg-coruscant-online"))
	d.Connect(mgLZ, mgCorp)
	d.Connect(mgLZ, mgOnline)

	// Level 4: Subscriptions under Platform MGs
	subMgmt := azure.General.Subscriptions(diagram.NodeLabel("sub-coruscant-mgmt"))
	d.Connect(mgMgmt, subMgmt)

	subIdentity := azure.General.Subscriptions(diagram.NodeLabel("sub-coruscant-identity"))
	d.Connect(mgIdentity, subIdentity)

	subSecProd := azure.General.Subscriptions(diagram.NodeLabel("sub-coruscant-security-prod"))
	subSecStg := azure.General.Subscriptions(diagram.NodeLabel("sub-coruscant-security-stg"))
	d.Connect(mgSec, subSecProd)
	d.Connect(mgSec, subSecStg)

	subHubProd := azure.General.Subscriptions(diagram.NodeLabel("sub-coruscant-hub-prod"))
	subHubStg := azure.General.Subscriptions(diagram.NodeLabel("sub-coruscant-hub-stg"))
	subEdgeProd := azure.General.Subscriptions(diagram.NodeLabel("sub-coruscant-edge-prod"))
	subEdgeStg := azure.General.Subscriptions(diagram.NodeLabel("sub-coruscant-edge-stg"))
	d.Connect(mgConn, subHubProd)
	d.Connect(mgConn, subHubStg)
	d.Connect(mgConn, subEdgeProd)
	d.Connect(mgConn, subEdgeStg)

	genSubFoundation(d, mgFoundation)

	return d.Render()
}

func genSubFoundation(d *diagram.Diagram, mgFoundation *diagram.Node) {
	// Create subscription node
	subFound := azure.General.Subscriptions(diagram.NodeLabel("sub-coruscant-foundation"))
	d.Connect(mgFoundation, subFound)

	d.Connect(subFound, azure.General.Policy(diagram.NodeLabel("Policy")))
	d.Connect(subFound, azure.General.CostManagement(diagram.NodeLabel("Cost Management")))
	d.Connect(subFound, azure.General.CostAlerts(diagram.NodeLabel("Cost Alerts")))

	// Create resource group with connections to subscription
	rgNode := azure.General.Resourcegroups(diagram.NodeLabel("rg-coruscant-tfbackend-01"))
	d.Connect(subFound, rgNode)

	// Add storage accounts connected to resource group
	stTFFoundation := azure.Storage.StorageAccounts(diagram.NodeLabel("stknccoruscanttfbfoundation01"))
	stTFDeployment := azure.Storage.StorageAccounts(diagram.NodeLabel("stknccoruscanttfbdeployment01"))
	d.Connect(rgNode, stTFFoundation)
	d.Connect(rgNode, stTFDeployment)

	// Each storage account has its own Defender enabled
	defenderFoundation := azure.Security.DefenderForCloud(diagram.NodeLabel("Defender"))
	defenderDeployment := azure.Security.DefenderForCloud(diagram.NodeLabel("Defender"))
	d.Connect(stTFFoundation, defenderFoundation)
	d.Connect(stTFDeployment, defenderDeployment)
}
