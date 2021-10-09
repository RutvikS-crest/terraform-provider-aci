package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateTACACSDestination(port string, host string, tacacs_monitoring_destination_group string, description string, nameAlias string, tacacsTacacsDestAttr models.TACACSDestinationAttributes) (*models.TACACSDestination, error) {
	rn := fmt.Sprintf(models.RntacacsTacacsDest, host, port)
	parentDn := fmt.Sprintf(models.ParentDntacacsTacacsDest, tacacs_monitoring_destination_group)
	tacacsTacacsDest := models.NewTACACSDestination(rn, parentDn, description, nameAlias, tacacsTacacsDestAttr)
	err := sm.Save(tacacsTacacsDest)
	return tacacsTacacsDest, err
}

func (sm *ServiceManager) ReadTACACSDestination(port string, host string, tacacs_monitoring_destination_group string) (*models.TACACSDestination, error) {
	dn := fmt.Sprintf(models.DntacacsTacacsDest, tacacs_monitoring_destination_group, host, port)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	tacacsTacacsDest := models.TACACSDestinationFromContainer(cont)
	return tacacsTacacsDest, nil
}

func (sm *ServiceManager) DeleteTACACSDestination(port string, host string, tacacs_monitoring_destination_group string) error {
	dn := fmt.Sprintf(models.DntacacsTacacsDest, tacacs_monitoring_destination_group, host, port)
	return sm.DeleteByDn(dn, models.TacacstacacsdestClassName)
}

func (sm *ServiceManager) UpdateTACACSDestination(port string, host string, tacacs_monitoring_destination_group string, description string, nameAlias string, tacacsTacacsDestAttr models.TACACSDestinationAttributes) (*models.TACACSDestination, error) {
	rn := fmt.Sprintf(models.RntacacsTacacsDest, host, port)
	parentDn := fmt.Sprintf(models.ParentDntacacsTacacsDest, tacacs_monitoring_destination_group)
	tacacsTacacsDest := models.NewTACACSDestination(rn, parentDn, description, nameAlias, tacacsTacacsDestAttr)
	tacacsTacacsDest.Status = "modified"
	err := sm.Save(tacacsTacacsDest)
	return tacacsTacacsDest, err
}

func (sm *ServiceManager) ListTACACSDestination(tacacs_monitoring_destination_group string) ([]*models.TACACSDestination, error) {
	dnUrl := fmt.Sprintf("%s/uni/fabric/tacacsgroup-%s/tacacsTacacsDest.json", models.BaseurlStr, tacacs_monitoring_destination_group)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.TACACSDestinationListFromContainer(cont)
	return list, err
}

// func (sm *ServiceManager) CreateRelationfileRsARemoteHostToEpg(parentDn, annotation, tDn string) error {
// 	dn := fmt.Sprintf("%s/rsARemoteHostToEpg", parentDn)
// 	containerJSON := []byte(fmt.Sprintf(`{
// 		"%s": {
// 			"attributes": {
// 				"dn": "%s",
// 				"annotation": "%s",
// 				"tDn": "%s"
// 			}
// 		}
// 	}`, "fileRsARemoteHostToEpg", dn, annotation, tDn))

// 	jsonPayload, err := container.ParseJSON(containerJSON)
// 	if err != nil {
// 		return err
// 	}
// 	req, err := sm.client.MakeRestRequest("POST", fmt.Sprintf("%s.json", sm.MOURL), jsonPayload, true)
// 	if err != nil {
// 		return err
// 	}
// 	cont, _, err := sm.client.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Printf("%+v", cont)
// 	return nil
// }

// func (sm *ServiceManager) DeleteRelationfileRsARemoteHostToEpg(parentDn string) error {
// 	dn := fmt.Sprintf("%s/rsARemoteHostToEpg", parentDn)
// 	return sm.DeleteByDn(dn, "fileRsARemoteHostToEpg")
// }

// func (sm *ServiceManager) ReadRelationfileRsARemoteHostToEpg(parentDn string) (interface{}, error) {
// 	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "fileRsARemoteHostToEpg")
// 	cont, err := sm.GetViaURL(dnUrl)
// 	contList := models.ListFromContainer(cont, "fileRsARemoteHostToEpg")

// 	if len(contList) > 0 {
// 		dat := models.G(contList[0], "tDn")
// 		return dat, err
// 	} else {
// 		return nil, err
// 	}
// }

// func (sm *ServiceManager) CreateRelationfileRsARemoteHostToEpp(parentDn, annotation, tDn string) error {
// 	dn := fmt.Sprintf("%s/rsARemoteHostToEpp", parentDn)
// 	containerJSON := []byte(fmt.Sprintf(`{
// 		"%s": {
// 			"attributes": {
// 				"dn": "%s",
// 				"annotation": "%s",
// 				"tDn": "%s"
// 			}
// 		}
// 	}`, "fileRsARemoteHostToEpp", dn, annotation, tDn))

// 	jsonPayload, err := container.ParseJSON(containerJSON)
// 	if err != nil {
// 		return err
// 	}
// 	req, err := sm.client.MakeRestRequest("POST", fmt.Sprintf("%s.json", sm.MOURL), jsonPayload, true)
// 	if err != nil {
// 		return err
// 	}
// 	cont, _, err := sm.client.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Printf("%+v", cont)
// 	return nil
// }

// func (sm *ServiceManager) DeleteRelationfileRsARemoteHostToEpp(parentDn string) error {
// 	dn := fmt.Sprintf("%s/rsARemoteHostToEpp", parentDn)
// 	return sm.DeleteByDn(dn, "fileRsARemoteHostToEpp")
// }

// func (sm *ServiceManager) ReadRelationfileRsARemoteHostToEpp(parentDn string) (interface{}, error) {
// 	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "fileRsARemoteHostToEpp")
// 	cont, err := sm.GetViaURL(dnUrl)
// 	contList := models.ListFromContainer(cont, "fileRsARemoteHostToEpp")

// 	if len(contList) > 0 {
// 		dat := models.G(contList[0], "tDn")
// 		return dat, err
// 	} else {
// 		return nil, err
// 	}
// }
