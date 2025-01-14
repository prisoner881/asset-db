// Copyright © by Jeff Foley 2017-2024. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.
// SPDX-License-Identifier: Apache-2.0

package repository

import (
	"encoding/json"
	"fmt"
	"time"

	oam "github.com/owasp-amass/open-asset-model"
	oamtls "github.com/owasp-amass/open-asset-model/certificate"
	"github.com/owasp-amass/open-asset-model/contact"
	"github.com/owasp-amass/open-asset-model/domain"
	"github.com/owasp-amass/open-asset-model/fingerprint"
	"github.com/owasp-amass/open-asset-model/network"
	"github.com/owasp-amass/open-asset-model/org"
	"github.com/owasp-amass/open-asset-model/people"
	oamreg "github.com/owasp-amass/open-asset-model/registration"
	"github.com/owasp-amass/open-asset-model/service"
	"github.com/owasp-amass/open-asset-model/source"
	"github.com/owasp-amass/open-asset-model/url"
	"gorm.io/datatypes"
)

// Asset represents an asset stored in the database.
type Asset struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement:true"`                               // The unique identifier of the asset.
	CreatedAt time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP();column=created_at"` // The creation timestamp of the asset.
	LastSeen  time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP();column=last_seen"`  // The last seen timestamp of the asset.
	Type      string         // The type of the asset.
	Content   datatypes.JSON // The JSON-encoded content of the asset.
}

// Relation represents a relationship between two assets stored in the database.
type Relation struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement:true"`              // The unique identifier of the relation.
	CreatedAt   time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP();"` // The creation timestamp of the relation.
	LastSeen    time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP();"` // The last seen timestamp of the relation.
	Type        string    // The type of the relation.
	FromAssetID uint64    // The ID of the asset from which the relation originates.
	ToAssetID   uint64    // The ID of the asset to which the relation points.
	FromAsset   Asset     // The asset from which the relation originates.
	ToAsset     Asset     // The asset to which the relation points.
}

// Parse parses the content of the asset into the corresponding Open Asset Model (OAM) asset type.
// It returns the parsed asset and an error, if any.
func (a *Asset) Parse() (oam.Asset, error) {
	var err error
	var asset oam.Asset

	switch a.Type {
	case string(oam.FQDN):
		var fqdn domain.FQDN

		err = json.Unmarshal(a.Content, &fqdn)
		asset = &fqdn
	case string(oam.NetworkEndpoint):
		var ne domain.NetworkEndpoint

		err = json.Unmarshal(a.Content, &ne)
		asset = &ne
	case string(oam.IPAddress):
		var ip network.IPAddress

		err = json.Unmarshal(a.Content, &ip)
		asset = &ip
	case string(oam.AutonomousSystem):
		var as network.AutonomousSystem

		err = json.Unmarshal(a.Content, &as)
		asset = &as
	case string(oam.AutnumRecord):
		var ar oamreg.AutnumRecord

		err = json.Unmarshal(a.Content, &ar)
		asset = &ar
	case string(oam.Netblock):
		var netblock network.Netblock

		err = json.Unmarshal(a.Content, &netblock)
		asset = &netblock
	case string(oam.IPNetRecord):
		var ipnetrec oamreg.IPNetRecord

		err = json.Unmarshal(a.Content, &ipnetrec)
		asset = &ipnetrec
	case string(oam.SocketAddress):
		var sa network.SocketAddress

		err = json.Unmarshal(a.Content, &sa)
		asset = &sa
	case string(oam.DomainRecord):
		var dr oamreg.DomainRecord

		err = json.Unmarshal(a.Content, &dr)
		asset = &dr
	case string(oam.Fingerprint):
		var fingerprint fingerprint.Fingerprint

		err = json.Unmarshal(a.Content, &fingerprint)
		asset = &fingerprint
	case string(oam.Organization):
		var organization org.Organization

		err = json.Unmarshal(a.Content, &organization)
		asset = &organization
	case string(oam.Person):
		var person people.Person

		err = json.Unmarshal(a.Content, &person)
		asset = &person
	case string(oam.Phone):
		var phone contact.Phone

		err = json.Unmarshal(a.Content, &phone)
		asset = &phone
	case string(oam.EmailAddress):
		var emailAddress contact.EmailAddress

		err = json.Unmarshal(a.Content, &emailAddress)
		asset = &emailAddress
	case string(oam.Location):
		var location contact.Location

		err = json.Unmarshal(a.Content, &location)
		asset = &location
	case string(oam.ContactRecord):
		var cr contact.ContactRecord

		err = json.Unmarshal(a.Content, &cr)
		asset = &cr
	case string(oam.TLSCertificate):
		var tlsCertificate oamtls.TLSCertificate

		err = json.Unmarshal(a.Content, &tlsCertificate)
		asset = &tlsCertificate
	case string(oam.URL):
		var url url.URL

		err = json.Unmarshal(a.Content, &url)
		asset = &url
	case string(oam.Source):
		var src source.Source

		err = json.Unmarshal(a.Content, &src)
		asset = &src
	case string(oam.Service):
		var serv service.Service

		err = json.Unmarshal(a.Content, &serv)
		asset = &serv
	default:
		return nil, fmt.Errorf("unknown asset type: %s", a.Type)
	}

	return asset, err
}

// JSONQuery generates a JSON query expression based on the asset's content.
// It returns the generated JSON query expression and an error, if any.
func (a *Asset) JSONQuery() (*datatypes.JSONQueryExpression, error) {
	asset, err := a.Parse()
	if err != nil {
		return nil, err
	}

	jsonQuery := datatypes.JSONQuery("content")
	switch v := asset.(type) {
	case *domain.FQDN:
		return jsonQuery.Equals(v.Name, "name"), nil
	case *domain.NetworkEndpoint:
		return jsonQuery.Equals(v.Address, "address"), nil
	case *network.SocketAddress:
		return jsonQuery.Equals(v.Address.String(), "address"), nil
	case *network.IPAddress:
		return jsonQuery.Equals(v.Address.String(), "address"), nil
	case *network.AutonomousSystem:
		return jsonQuery.Equals(v.Number, "number"), nil
	case *network.Netblock:
		return jsonQuery.Equals(v.CIDR.String(), "cidr"), nil
	case *oamreg.IPNetRecord:
		return jsonQuery.Equals(v.Handle, "handle"), nil
	case *oamreg.AutnumRecord:
		return jsonQuery.Equals(v.Handle, "handle"), nil
	case *oamreg.DomainRecord:
		return jsonQuery.Equals(v.Domain, "domain"), nil
	case *fingerprint.Fingerprint:
		return jsonQuery.Equals(v.Value, "value"), nil
	case *org.Organization:
		return jsonQuery.Equals(v.Name, "name"), nil
	case *people.Person:
		return jsonQuery.Equals(v.FullName, "full_name"), nil
	case *contact.Phone:
		return jsonQuery.Equals(v.Raw, "raw"), nil
	case *contact.EmailAddress:
		return jsonQuery.Equals(v.Address, "address"), nil
	case *contact.Location:
		return jsonQuery.Equals(v.Address, "address"), nil
	case *contact.ContactRecord:
		return jsonQuery.Equals(v.DiscoveredAt, "discovered_at"), nil
	case *oamtls.TLSCertificate:
		return jsonQuery.Equals(v.SerialNumber, "serial_number"), nil
	case *url.URL:
		return jsonQuery.Equals(v.Raw, "url"), nil
	case *source.Source:
		return jsonQuery.Equals(v.Name, "name"), nil
	case *service.Service:
		return jsonQuery.Equals(v.Identifier, "identifier"), nil
	}

	return nil, fmt.Errorf("unknown asset type: %s", a.Type)
}
