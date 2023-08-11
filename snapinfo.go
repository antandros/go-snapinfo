package snapinfo

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/antandros/go-pkgparser"
	"github.com/antandros/go-pkgparser/model"
	"github.com/snapcore/snapd/client"
)

func Parse(p *pkgparser.Parser) error {
	snapclient := client.New(nil)
	snaps, err := snapclient.List(nil, nil)
	if err != nil {
		return err
	}

	for _, snapn := range snaps {
		packageItem := p.CreateModel()

		mapItems := map[string]interface{}{
			"Package":             snapn.Title,
			"Cnf-Visible-Pkgname": snapn.Title,
			"Version":             snapn.Version,
			"Description":         snapn.Description,
			"Installed-Size":      snapn.InstalledSize,
			"License":             snapn.License,
			"Homepage":            snapn.Website,
			"Revision":            snapn.Revision,
			"Status":              snapn.Status,
			"Section":             snapn.Channel,
			"Vendor":              snapn.Publisher.DisplayName,
			"Maintainer":          fmt.Sprintf("%s <%s>", snapn.Publisher.DisplayName, strings.ReplaceAll(snapn.Contact, "mailto:", "")),
		}
		keys := make([]string, 0, len(mapItems))
		for k := range mapItems {
			keys = append(keys, k)
		}
		volOf := reflect.ValueOf(snapn).Elem()
		typOf := reflect.TypeOf(*snapn)
		for i := 0; i < typOf.NumField(); i++ {
			field := typOf.Field(i)
			hasItem := false
			for k := range keys {
				if strings.EqualFold(keys[k], field.Name) {
					hasItem = true
				}
			}
			if !hasItem {
				valField := volOf.FieldByName(field.Name)
				if valField.IsValid() && valField.Interface() != nil {
					mapItems[field.Name] = valField
				}

			}
		}

		for key, valn := range mapItems {
			packageItem, err = p.SetValue(key, fmt.Sprintf("%v", valn), packageItem)
			if err != nil {
				fmt.Println("Error", err.Error(), key)
			}
		}
		p.Packages = append(p.Packages, packageItem)
	}
	return nil
}
func GetPackages() ([]model.Package, error) {
	var packages []model.Package
	p := new(pkgparser.Parser)
	p.Model = model.Package{}
	err := p.StructParse()
	if err != nil {
		return nil, err
	}
	err = Parse(p)
	if err != nil {
		return nil, err
	}
	for _, i := range p.Packages {
		item, ok := i.(*model.Package)
		if !ok {
			return nil, errors.New("struct conversion failed")

		}
		packages = append(packages, *item)
	}
	return packages, nil
}
