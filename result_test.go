package sparql

import (
	"encoding/xml"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
)

func TestDecodeXMLQueryResult(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		reader := ioutil.NopCloser(strings.NewReader(``))
		if _, err := DecodeXMLQueryResult(reader); err != io.EOF {
			t.Errorf("DecodeXMLQueryResult() error = %v", err)
			return
		}
	})
	t.Run("success", func(t *testing.T) {
		reader := ioutil.NopCloser(strings.NewReader(`<head></head>`))
		got, err := DecodeXMLQueryResult(reader)
		if err != nil {
			t.Errorf("DecodeXMLQueryResult() error = %v", err)
			return
		}

		if got, want := got.Variables(), []string{}; !reflect.DeepEqual(got, want) {
			t.Errorf("DecodeXMLQueryResult() = %v, want %v", got, want)
		}
	})
}

func Test_decodeVariables(t *testing.T) {
	t.Run("bad XML", func(t *testing.T) {
		decoder := xml.NewDecoder(strings.NewReader(`<head><variable></head>`))
		_, err := decodeVariables(decoder)
		if _, ok := err.(*xml.SyntaxError); !ok {
			t.Errorf("decodeVariables() error = %v", err)
			return
		}
	})
	t.Run("success", func(t *testing.T) {
		decoder := xml.NewDecoder(strings.NewReader(`<head><variable name="foo" /></head>`))
		got, err := decodeVariables(decoder)
		if err != nil {
			t.Errorf("decodeVariables() error = %v", err)
			return
		}
		want := []string{"foo"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("decodeVariables() = %v, want %v", got, want)
		}
		if _, err := decodeVariables(decoder); err != io.EOF {
			t.Errorf("decodeVariables() error = %v", err)
			return
		}
	})
}

func TestXMLQueryResult_Variables(t *testing.T) {
	x := &XMLQueryResult{
		variables: []string{"foo"},
	}
	want := []string{"foo"}
	if got := x.Variables(); !reflect.DeepEqual(got, want) {
		t.Errorf("XMLQueryResult.Variables() = %v, want %v", got, want)
	}
}

func TestXMLQueryResult_Next(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		x := &XMLQueryResult{
			decoder: xml.NewDecoder(strings.NewReader("")),
		}
		if _, err := x.Next(); err != io.EOF {
			t.Errorf("XMLQueryResult.Next() error = %v", err)
			return
		}
	})
	t.Run("success", func(t *testing.T) {
		x := &XMLQueryResult{
			decoder: xml.NewDecoder(strings.NewReader(`<result>
<binding name="x"><bnode>r2</bnode></binding>
</result>`)),
		}
		got, err := x.Next()
		if err != nil {
			t.Errorf("XMLQueryResult.Next() error = %v", err)
			return
		}
		want := map[string]Value{"x": BNode("r2")}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("XMLQueryResult.Next() = %v, want %v", got, want)
		}
	})
}

func Test_decodeResult(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		decoder := xml.NewDecoder(strings.NewReader(`<result></result>`))
		if _, err := decoder.Token(); err != nil {
			t.Fatal(err)
		}
		got, err := decodeResult(decoder, 0)
		if err != nil {
			t.Errorf("decodeResult() error = %v", err)
		}
		want := make(map[string]Value, 0)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("decodeResult() = %v, want %v", got, want)
		}
	})
	t.Run("bad result", func(t *testing.T) {
		decoder := xml.NewDecoder(strings.NewReader(`<result>`))
		if _, err := decoder.Token(); err != nil {
			t.Fatal(err)
		}
		_, err := decodeResult(decoder, 0)
		if _, ok := err.(*xml.SyntaxError); !ok {
			t.Errorf("decodeResult() error = %v", err)
		}
	})
	t.Run("bad binding", func(t *testing.T) {
		decoder := xml.NewDecoder(strings.NewReader(`<result><binding></result>`))
		if _, err := decoder.Token(); err != nil {
			t.Fatal(err)
		}
		_, err := decodeResult(decoder, 0)
		if _, ok := err.(*xml.SyntaxError); !ok {
			t.Errorf("decodeResult() error = %v", err)
		}
	})
	t.Run("success", func(t *testing.T) {
		decoder := xml.NewDecoder(strings.NewReader(`<result>
<binding name="x"><bnode>r2</bnode></binding>
<binding name="hpage"><uri>http://work.example.org/bob/</uri></binding>
<binding name="name"><literal xml:lang="en">Bob</literal></binding>
<binding name="age"><literal datatype="http://www.w3.org/2001/XMLSchema#integer">30</literal></binding>
<binding name="mbox"><uri>mailto:bob@work.example.org</uri></binding>
</result>`))
		if _, err := decoder.Token(); err != nil {
			t.Fatal(err)
		}
		got, err := decodeResult(decoder, 5)
		if err != nil {
			t.Errorf("decodeResult() error = %v", err)
			return
		}
		want := map[string]Value{
			"x":     BNode("r2"),
			"hpage": URI("http://work.example.org/bob/"),
			"name": Literal{
				Value:       "Bob",
				LanguageTag: "en",
			},
			"age": Literal{
				Value:    "30",
				DataType: URI("http://www.w3.org/2001/XMLSchema#integer"),
			},
			"mbox": URI("mailto:bob@work.example.org"),
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("decodeResult() = %v, want %v", got, want)
		}
	})
}

func Test_decodeBinding(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		decoder := xml.NewDecoder(strings.NewReader(`<binding name="x"></binding>`))
		token, err := decoder.Token()
		if err != nil {
			t.Fatal(err)
		}
		element := token.(xml.StartElement)
		if _, _, err = decodeBinding(decoder, &element); err != io.EOF {
			t.Errorf("decodeBinding() error = %v", err)
		}
	})
	t.Run("bnode", func(t *testing.T) {
		decoder := xml.NewDecoder(strings.NewReader(`<binding name="x">
<bnode>r2</bnode>
</binding>`))
		token, err := decoder.Token()
		if err != nil {
			t.Fatal(err)
		}
		element := token.(xml.StartElement)
		got, got1, err := decodeBinding(decoder, &element)
		if err != nil {
			t.Errorf("decodeBinding() error = %v", err)
			return
		}
		want, want1 := "x", BNode("r2")
		if got != want {
			t.Errorf("decodeBinding() got = %v, want %v", got, want)
		}
		if !reflect.DeepEqual(got1, want1) {
			t.Errorf("decodeBinding() got1 = %v, want %v", got1, want1)
		}
	})
	t.Run("bad bnode", func(t *testing.T) {
		decoder := xml.NewDecoder(strings.NewReader(`<binding name="x">
<bnode>r2
</binding>`))
		token, err := decoder.Token()
		if err != nil {
			t.Fatal(err)
		}
		element := token.(xml.StartElement)
		_, _, err = decodeBinding(decoder, &element)
		if _, ok := err.(*xml.SyntaxError); !ok {
			t.Errorf("decodeBinding() error = %v", err)
		}
	})
	t.Run("uri", func(t *testing.T) {
		decoder := xml.NewDecoder(strings.NewReader(`<binding name="x">
<uri>http://work.example.org/bob/</uri>
</binding>`))
		token, err := decoder.Token()
		if err != nil {
			t.Fatal(err)
		}
		element := token.(xml.StartElement)
		got, got1, err := decodeBinding(decoder, &element)
		if err != nil {
			t.Errorf("decodeBinding() error = %v", err)
			return
		}
		want, want1 := "x", URI("http://work.example.org/bob/")
		if got != want {
			t.Errorf("decodeBinding() got = %v, want %v", got, want)
		}
		if !reflect.DeepEqual(got1, want1) {
			t.Errorf("decodeBinding() got1 = %v, want %v", got1, want1)
		}
	})
	t.Run("bad uri", func(t *testing.T) {
		decoder := xml.NewDecoder(strings.NewReader(`<binding name="x">
<uri>http://work.example.org/bob/
</binding>`))
		token, err := decoder.Token()
		if err != nil {
			t.Fatal(err)
		}
		element := token.(xml.StartElement)
		_, _, err = decodeBinding(decoder, &element)
		if _, ok := err.(*xml.SyntaxError); !ok {
			t.Errorf("decodeBinding() error = %v", err)
		}
	})
	t.Run("literal with lang", func(t *testing.T) {
		decoder := xml.NewDecoder(strings.NewReader(`<binding name="x">
<literal xml:lang="en">Bob</literal>
</binding>`))
		token, err := decoder.Token()
		if err != nil {
			t.Fatal(err)
		}
		element := token.(xml.StartElement)
		got, got1, err := decodeBinding(decoder, &element)
		if err != nil {
			t.Errorf("decodeBinding() error = %v", err)
			return
		}
		want, want1 := "x", Literal{
			Value:       "Bob",
			LanguageTag: "en",
		}
		if got != want {
			t.Errorf("decodeBinding() got = %v, want %v", got, want)
		}
		if !reflect.DeepEqual(got1, want1) {
			t.Errorf("decodeBinding() got1 = %v, want %v", got1, want1)
		}
	})
	t.Run("literal with datatype", func(t *testing.T) {
		decoder := xml.NewDecoder(strings.NewReader(`<binding name="x">
<literal datatype="http://www.w3.org/2001/XMLSchema#integer">30</literal>
</binding>`))
		token, err := decoder.Token()
		if err != nil {
			t.Fatal(err)
		}
		element := token.(xml.StartElement)
		got, got1, err := decodeBinding(decoder, &element)
		if err != nil {
			t.Errorf("decodeBinding() error = %v", err)
			return
		}
		want, want1 := "x", Literal{
			Value:    "30",
			DataType: URI("http://www.w3.org/2001/XMLSchema#integer"),
		}
		if got != want {
			t.Errorf("decodeBinding() got = %v, want %v", got, want)
		}
		if !reflect.DeepEqual(got1, want1) {
			t.Errorf("decodeBinding() got1 = %v, want %v", got1, want1)
		}
	})
	t.Run("literal with datatype and lang", func(t *testing.T) {
		decoder := xml.NewDecoder(strings.NewReader(`<binding name="x">
<literal xml:lang="en" datatype="http://www.w3.org/2001/XMLSchema#string">foo</literal>
</binding>`))
		token, err := decoder.Token()
		if err != nil {
			t.Fatal(err)
		}
		element := token.(xml.StartElement)
		got, got1, err := decodeBinding(decoder, &element)
		if err != nil {
			t.Errorf("decodeBinding() error = %v", err)
			return
		}
		want, want1 := "x", Literal{
			Value:       "foo",
			DataType:    URI("http://www.w3.org/2001/XMLSchema#string"),
			LanguageTag: "en",
		}
		if got != want {
			t.Errorf("decodeBinding() got = %v, want %v", got, want)
		}
		if !reflect.DeepEqual(got1, want1) {
			t.Errorf("decodeBinding() got1 = %v, want %v", got1, want1)
		}
	})
	t.Run("bad literal", func(t *testing.T) {
		decoder := xml.NewDecoder(strings.NewReader(`<binding name="x">
<literal>foo
</binding>`))
		token, err := decoder.Token()
		if err != nil {
			t.Fatal(err)
		}
		element := token.(xml.StartElement)
		_, _, err = decodeBinding(decoder, &element)
		if _, ok := err.(*xml.SyntaxError); !ok {
			t.Errorf("decodeBinding() error = %v", err)
		}
	})
	t.Run("unknown binding", func(t *testing.T) {
		decoder := xml.NewDecoder(strings.NewReader(`<binding name="x">
<foo>bar</foo>
</binding>`))
		token, err := decoder.Token()
		if err != nil {
			t.Fatal(err)
		}
		element := token.(xml.StartElement)
		if _, _, err := decodeBinding(decoder, &element); err == nil {
			t.Errorf("decodeBinding() error = %v", err)
		}
	})
}

func TestXMLQueryResult_Boolean(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		x := &XMLQueryResult{
			decoder: xml.NewDecoder(strings.NewReader(``)),
		}
		if _, err := x.Boolean(); err != io.EOF {
			t.Errorf("XMLQueryResult.Boolean() error = %v", err)
			return
		}
	})
	t.Run("bad XML", func(t *testing.T) {
		x := &XMLQueryResult{
			decoder: xml.NewDecoder(strings.NewReader(`<boolean>`)),
		}
		_, err := x.Boolean()
		if _, ok := err.(*xml.SyntaxError); !ok {
			t.Errorf("XMLQueryResult.Boolean() error = %v", err)
			return
		}
	})
	t.Run("success", func(t *testing.T) {
		x := &XMLQueryResult{
			decoder: xml.NewDecoder(strings.NewReader(`<boolean>true</boolean>`)),
		}
		got, err := x.Boolean()
		if err != nil {
			t.Errorf("XMLQueryResult.Boolean() error = %v", err)
			return
		}
		if want := true; got != want {
			t.Errorf("XMLQueryResult.Boolean() = %v, want %v", got, want)
		}
	})
}

func TestXMLQueryResult_Close(t *testing.T) {
	x := &XMLQueryResult{
		r: ioutil.NopCloser(strings.NewReader("")),
	}
	if err := x.Close(); err != nil {
		t.Errorf("XMLQueryResult.Close() error = %v", err)
	}
}
