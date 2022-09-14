// sitemap.go

// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    sitemap, err := Deserialize(bytes)
//    bytes, err = sitemap.Serialize()

package models

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
)

func Deserialize(data []byte) (Sitemap, error) {
	var r Sitemap
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Sitemap) SerializeToJSON() ([]byte, error) {
	// return json.MarshalIndent(r, "", "    ")
	return json.Marshal(r)
}

func (r *Sitemap) SerializeToXML() ([]byte, error) {
	return xml.MarshalIndent(r, "", "    ")
}

func (r *Sitemap) SerializeToFile(filePath string) error {
	marshalled, err := json.MarshalIndent(r, "", "    ")
	if err != nil {
		return fmt.Errorf("unable to marshall to json file %v due to error: \n\t%v", filePath, err)
	}
	jsonFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("unable to create json file %v due to error: \n\t%v", filePath, err)
	}

	_, err = jsonFile.Write(marshalled)
	if err != nil {
		return fmt.Errorf("unable to Write to json file %v due to error: \n\t%v", filePath, err)
	}
	return nil
}

type Sitemap struct {
	XMLName xml.Name `xml:"sitemap" json:"-"`
	Xmlns   string   `json:"xmlns" xml:"xmlns,attr"`
	Urlset  URLSet   `json:"urlset" xml:"urlset"`
}

func NewSitemap(urlset URLSet) *Sitemap {
	return &Sitemap{Urlset: urlset, Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9"}
}
