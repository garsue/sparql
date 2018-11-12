package client

import (
	"encoding/xml"
	"fmt"
	"io"
)

// QueryResult is a SPARQL query result.
type QueryResult interface {
	Variables() []string
	Next() (map[string]Value, error)
	Boolean() (bool, error)

	io.Closer
}

// Value is an interface holding one of the binding (or boolean) types:
// URI, Literal, BNode or bool.
type Value interface{}

type (
	// BNode represents blank node.
	BNode string

	// Literal http://www.w3.org/TR/2004/REC-rdf-concepts-20040210/#dfn-literal
	Literal struct {
		Value       string
		DataType    IRIRef
		LanguageTag string
	}
)

// UnmarshalXML unmarshals the literal element.
// `datatype` is decoded as URI. `lang` must have xml namespace.
func (l *Literal) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "datatype":
			l.DataType = URI(attr.Value)
		case "lang":
			if attr.Name.Space == "http://www.w3.org/XML/1998/namespace" {
				l.LanguageTag = attr.Value
			}
		}
	}
	return d.DecodeElement(&l.Value, &start)
}

// XMLQueryResult is the implementation to decode SPARQL Query Results XML Format.
// https://www.w3.org/TR/rdf-sparql-XMLres/
type XMLQueryResult struct {
	r         io.ReadCloser
	variables []string
	decoder   *xml.Decoder
}

// DecodeXMLQueryResult decodes responded XML Query Result.
func DecodeXMLQueryResult(r io.ReadCloser) (QueryResult, error) {
	decoder := xml.NewDecoder(r)
	variables, err := decodeVariables(decoder)
	if err != nil {
		return nil, err
	}
	return &XMLQueryResult{
		r:         r,
		variables: variables,
		decoder:   decoder,
	}, nil
}

type head struct {
	Variables []struct {
		Name string `xml:"name,attr"`
	} `xml:"variable"`
}

func decodeVariables(decoder *xml.Decoder) ([]string, error) {
	h, err := decodeHead(decoder)
	if err != nil {
		return nil, err
	}
	vs := make([]string, 0, len(h.Variables))
	for _, v := range h.Variables {
		vs = append(vs, v.Name)
	}
	return vs, nil
}

func decodeHead(decoder *xml.Decoder) (head, error) {
	for {
		t, err := startElement(decoder)
		if err != nil {
			return head{}, err
		}
		if t.Name.Local == "head" {
			var h head
			if err := decoder.DecodeElement(&h, &t); err != nil {
				return head{}, err
			}
			return h, nil
		}
	}
}

// Variables returns query variables.
func (x *XMLQueryResult) Variables() []string {
	return x.variables
}

func (x *XMLQueryResult) Next() (map[string]Value, error) {
	for {
		t, err := startElement(x.decoder)
		if err != nil {
			return nil, err
		}
		if t.Name.Local == "result" {
			return decodeResult(x.decoder, len(x.variables))
		}
	}
}

func decodeResult(decoder *xml.Decoder, size int) (map[string]Value, error) {
	bindings := make(map[string]Value, size)
	for {
		token, err := decoder.Token()
		if err != nil {
			return nil, err
		}
		switch token := token.(type) {
		case xml.StartElement:
			if token.Name.Local == "binding" {
				name, value, err := decodeBinding(decoder, &token)
				if err != nil {
					return nil, err
				}
				bindings[name] = value
			}
		case xml.EndElement:
			if token.Name.Local == "result" {
				return bindings, nil
			}
		}
	}
}

func decodeBinding(decoder *xml.Decoder, token *xml.StartElement) (string, Value, error) {
	name := nameAttr(token)
	se, err := startElement(decoder)
	if err != nil {
		return "", nil, err
	}
	switch se.Name.Local {
	case "uri":
		var uri URI
		if err := decoder.DecodeElement(&uri, &se); err != nil {
			return "", nil, err
		}
		return name, uri, nil
	case "literal":
		var literal Literal
		if err := decoder.DecodeElement(&literal, &se); err != nil {
			return "", nil, err
		}
		return name, literal, nil
	case "bnode":
		var bnode BNode
		if err := decoder.DecodeElement(&bnode, &se); err != nil {
			return "", nil, err
		}
		return name, bnode, nil
	default:
		return "", nil, fmt.Errorf("unknown binding %v", se.Name.Local)
	}
}

func nameAttr(token *xml.StartElement) string {
	for _, attr := range token.Attr {
		if attr.Name.Local == "name" {
			return attr.Value
		}
	}
	return ""
}

func startElement(r xml.TokenReader) (xml.StartElement, error) {
	for {
		token, err := r.Token()
		if err != nil {
			return xml.StartElement{}, err
		}
		se, ok := token.(xml.StartElement)
		if ok {
			return se, nil
		}
	}
}

func (x *XMLQueryResult) Boolean() (bool, error) {
	for {
		t, err := startElement(x.decoder)
		if err != nil {
			return false, err
		}
		if t.Name.Local == "boolean" {
			var boolean bool
			if err := x.decoder.DecodeElement(&boolean, &t); err != nil {
				return false, err
			}
			return boolean, nil
		}
	}
}

func (x *XMLQueryResult) Close() error {
	return x.r.Close()
}
