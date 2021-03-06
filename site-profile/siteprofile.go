package siteprofile

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/ApogeeNetworking/dnacenter/models"
	"github.com/ApogeeNetworking/dnacenter/requests"
	"github.com/ApogeeNetworking/dnacenter/templates"
)

// Service ...package
type Service struct {
	baseURL string
	http    *requests.Req
}

// New creates an instance of a DNA-C NETWORK-PROFILE Service
func New(uri string, r *requests.Req) *Service {
	parsedURL, _ := url.Parse(uri)
	uriHost := parsedURL.Host
	ep := fmt.Sprintf("https://%s/api/v1/siteprofile", uriHost)
	return &Service{baseURL: ep, http: r}
}

// Profile ...
type Profile struct {
	ID        string `json:"siteProfileUuid"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	NameSpace string `json:"namespace"`
	SiteCount int    `json:"siteCount,omitempty"`
	Sites     []Site `json:"sites,omitempty"`
}

// Site ...
type Site struct {
	IsInherited bool   `json:"isInherited"`
	Name        string `json:"name"`
	ID          string `json:"uuid"`
}

// Get ...
func (s *Service) Get() ([]Profile, error) {
	type getResp struct {
		Response []Profile `json:"response"`
	}
	var resp getResp
	res, err := s.http.MakeReq(s.baseURL, "GET", nil)
	if err != nil {
		return resp.Response, fmt.Errorf("%v", err)
	}
	defer res.Body.Close()
	if err = json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return resp.Response, fmt.Errorf("%v", err)
	}
	return resp.Response, nil
}

// GetByID ...
func (s *Service) GetByID(id string) (Profile, error) {
	qs := "includeSites=true&excludeSettings=true&populated=false"
	uri := fmt.Sprintf("%s/%s?%s", s.baseURL, id, qs)
	type getResp struct {
		Response Profile `json:"response"`
	}
	var resp getResp
	res, err := s.http.MakeReq(uri, "GET", nil)
	if err != nil {
		return resp.Response, fmt.Errorf("%v", err)
	}
	defer res.Body.Close()
	if err = json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return resp.Response, fmt.Errorf("%v", err)
	}
	return resp.Response, nil
}

// GetSiteTemplates ...
func (s *Service) GetSiteTemplates(siteID string) []templates.Template {
	ep := fmt.Sprintf("/site/%v", siteID)
	res, err := s.http.MakeReq(s.baseURL+ep, "GET", nil)
	if err != nil {
	}
	defer res.Body.Close()

	type Attr struct {
		Key   string `json:"key"`
		Value string `json:"value"`
		Attr  []Attr `json:"attribs"`
	}
	type ProfAttr struct {
		ProfAttr []Attr `json:"profileAttributes"`
	}
	type SiteProfSiteTemps struct {
		Resp []ProfAttr `json:"response"`
	}
	var profAttrResp SiteProfSiteTemps
	json.NewDecoder(res.Body).Decode(&profAttrResp)
	profAttrs := profAttrResp.Resp
	var templs []templates.Template
	for _, attr := range profAttrs {
		for _, pAttr := range attr.ProfAttr {
			var templ templates.Template
			if pAttr.Key == "day0.templates" {
				for _, dfAttr := range pAttr.Attr {
					for _, dsAttr := range dfAttr.Attr {
						for _, dtAttr := range dsAttr.Attr {
							for _, dtagAttr := range dtAttr.Attr {
								if dtagAttr.Key == "template.id" {
									templ.ID = dtagAttr.Value
									for _, tmplAttr := range dtagAttr.Attr {
										if tmplAttr.Key == "template.name" {
											templ.Name = tmplAttr.Value
											templs = append(templs, templ)
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return templs
}

// Create ...
// func (s *Service) Create(name string) {
/*
	{
		"siteProfileUuid": "",
		"version" : 0,
		"name" name,
		"namespace" : "switching",
		"status": "",
		"lastUpdatedBy": "",
		"lastUpdatedDatetime": 0,
		"siteCount": 0,
		"profileAttributes": [],
		"id": "",
		"attributesList": [],
		"siteProfileType": "",
		"namingPrefix": "",
		"primaryDeviceType": "",
		"secondaryDeviceType": "",
		"interfaceList": [],
		"groupTypeList": [],
		"siteAssociationId": "",
		"sites": []
	}
*/
// }

// AssignSite ...
func (s *Service) AssignSite(profileID, siteID string) (models.Task, error) {
	// POST /siteprofile/{profileID}/site/{siteID}
	type spResp struct {
		Response models.Task `json:"response"`
	}
	var resp spResp
	uri := fmt.Sprintf("%s/%s/site/%s", s.baseURL, profileID, siteID)
	res, err := s.http.MakeReq(uri, "POST", nil)
	if err != nil {
		return resp.Response, fmt.Errorf("%v", err)
	}
	defer res.Body.Close()
	if err = json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return resp.Response, fmt.Errorf("%v", err)
	}
	return resp.Response, nil
}

// RemoveSite ...
func (s *Service) RemoveSite(profileID, siteID string) (models.Task, error) {
	/* DELETE /siteprofile/{profileID}/site/{siteID} */
	type spResp struct {
		Response models.Task `json:"response"`
	}
	var resp spResp
	uri := fmt.Sprintf("%s/%s/site/%s", s.baseURL, profileID, siteID)
	res, err := s.http.MakeReq(uri, "DELETE", nil)
	if err != nil {
		return resp.Response, fmt.Errorf("%v", err)
	}
	defer res.Body.Close()
	if err = json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return resp.Response, fmt.Errorf("%v", err)
	}
	return resp.Response, nil
}

/*
// POST to /api/v1/siteprofile
{
    "siteProfileUuid": "",
    "version": 0,
    "name": "2960Profile",
    "namespace": "switching",
    "status": "",
    "lastUpdatedBy": "",
    "lastUpdatedDatetime": 0,
    "siteCount": 0,
    "profileAttributes": [],
    "id": "",
    "attributesList": [],
    "siteProfileType": "",
    "namingPrefix": "",
    "primaryDeviceType": "",
    "secondaryDeviceType": "",
    "interfaceList": [],
    "groupTypeList": [],
    "siteAssociationId": "",
    "sites": []
}

{
    "response": {
        "siteProfileUuid": "dfa7f2df-15ec-4d42-ae23-61d438cdd71b",
        "version": 1,
        "name": "2960Profile",
        "namespace": "switching",
        "status": "draft",
        "lastUpdatedBy": "admin",
        "lastUpdatedDatetime": 1587758418061
    },
    "version": "1.0"
}
*/
