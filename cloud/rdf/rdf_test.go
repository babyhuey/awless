package rdf

import "testing"

func TestAllLabelsHaveProperties(t *testing.T) {
	for prop, label := range Labels {
		if _, ok := Properties[label]; !ok {
			t.Fatalf("prop '%s' with label '%s' has no corresponding property", prop, label)
		}
	}
}

func TestAllPropertiesHaveLabel(t *testing.T) {
	for label, rdf := range Properties {
		found := false
		for _, v := range Labels {
			if label == v {
				found = true
			}
		}
		if !found && rdf.RdfType != RdfsSubProperty {
			t.Fatalf("rdf prop with label '%s' has no corresponding entry in labels", label)
		}
	}
}

func TestRDFPropertiesGet(t *testing.T) {
	t.Run("existing property", func(t *testing.T) {
		prop, err := Properties.Get(Account)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if prop.ID != Account {
			t.Fatalf("expected ID %q, got %q", Account, prop.ID)
		}
		if prop.RdfsLabel != "Account" {
			t.Fatalf("expected label 'Account', got %q", prop.RdfsLabel)
		}
	})

	t.Run("non-existing property", func(t *testing.T) {
		_, err := Properties.Get("cloud:nonexistent")
		if err == nil {
			t.Fatal("expected error for non-existing property")
		}
	})
}

func TestRDFPropertiesIsRDFProperty(t *testing.T) {
	t.Run("is an RDF property", func(t *testing.T) {
		if !Properties.IsRDFProperty(Account) {
			t.Fatalf("expected %q to be an RDF property", Account)
		}
	})

	t.Run("non-existing key", func(t *testing.T) {
		if Properties.IsRDFProperty("cloud:nonexistent") {
			t.Fatal("expected false for non-existing key")
		}
	})

	t.Run("sub-property is not RDF property", func(t *testing.T) {
		// PortRange is rdfs:subPropertyOf, not rdf:Property
		if Properties.IsRDFProperty(PortRange) {
			t.Fatalf("expected %q to not be an RDF property (it's a subProperty)", PortRange)
		}
	})
}

func TestRDFPropertiesIsRDFSubProperty(t *testing.T) {
	t.Run("is a sub-property", func(t *testing.T) {
		if !Properties.IsRDFSubProperty(PortRange) {
			t.Fatalf("expected %q to be an RDF sub-property", PortRange)
		}
	})

	t.Run("regular property is not sub-property", func(t *testing.T) {
		if Properties.IsRDFSubProperty(Account) {
			t.Fatalf("expected %q to not be an RDF sub-property", Account)
		}
	})

	t.Run("non-existing key", func(t *testing.T) {
		if Properties.IsRDFSubProperty("cloud:nonexistent") {
			t.Fatal("expected false for non-existing key")
		}
	})
}

func TestRDFPropertiesIsRDFList(t *testing.T) {
	t.Run("is a list", func(t *testing.T) {
		// Actions is defined by rdfs:list
		if !Properties.IsRDFList(Actions) {
			t.Fatalf("expected %q to be an RDF list", Actions)
		}
	})

	t.Run("not a list", func(t *testing.T) {
		// Account is defined by rdfs:Literal
		if Properties.IsRDFList(Account) {
			t.Fatalf("expected %q to not be an RDF list", Account)
		}
	})

	t.Run("non-existing key", func(t *testing.T) {
		if Properties.IsRDFList("cloud:nonexistent") {
			t.Fatal("expected false for non-existing key")
		}
	})
}

func TestRDFPropertiesGetRDFId(t *testing.T) {
	t.Run("existing label", func(t *testing.T) {
		// Find a label to test with
		var testLabel string
		var expectedID string
		for label, id := range Labels {
			testLabel = label
			expectedID = id
			break
		}
		id, err := Properties.GetRDFId(testLabel)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if id != expectedID {
			t.Fatalf("expected id %q, got %q", expectedID, id)
		}
	})

	t.Run("non-existing label", func(t *testing.T) {
		_, err := Properties.GetRDFId("NonExistentLabel")
		if err == nil {
			t.Fatal("expected error for non-existing label")
		}
	})
}

func TestRDFPropertiesGetDataType(t *testing.T) {
	t.Run("string data type", func(t *testing.T) {
		dt, err := Properties.GetDataType(Account)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if dt != XsdString {
			t.Fatalf("expected %q, got %q", XsdString, dt)
		}
	})

	t.Run("boolean data type", func(t *testing.T) {
		dt, err := Properties.GetDataType(ActionsEnabled)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if dt != XsdBoolean {
			t.Fatalf("expected %q, got %q", XsdBoolean, dt)
		}
	})

	t.Run("int data type", func(t *testing.T) {
		dt, err := Properties.GetDataType(ActiveServicesCount)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if dt != XsdInt {
			t.Fatalf("expected %q, got %q", XsdInt, dt)
		}
	})

	t.Run("non-existing property", func(t *testing.T) {
		_, err := Properties.GetDataType("cloud:nonexistent")
		if err == nil {
			t.Fatal("expected error for non-existing property")
		}
	})
}

func TestRDFPropertiesGetLabel(t *testing.T) {
	t.Run("existing property", func(t *testing.T) {
		label, err := Properties.GetLabel(Account)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if label != "Account" {
			t.Fatalf("expected 'Account', got %q", label)
		}
	})

	t.Run("non-existing property", func(t *testing.T) {
		_, err := Properties.GetLabel("cloud:nonexistent")
		if err == nil {
			t.Fatal("expected error for non-existing property")
		}
	})
}

func TestRDFPropertiesGetDefinedBy(t *testing.T) {
	t.Run("literal property", func(t *testing.T) {
		def, err := Properties.GetDefinedBy(Account)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if def != RdfsLiteral {
			t.Fatalf("expected %q, got %q", RdfsLiteral, def)
		}
	})

	t.Run("list property", func(t *testing.T) {
		def, err := Properties.GetDefinedBy(Actions)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if def != RdfsList {
			t.Fatalf("expected %q, got %q", RdfsList, def)
		}
	})

	t.Run("non-existing property", func(t *testing.T) {
		_, err := Properties.GetDefinedBy("cloud:nonexistent")
		if err == nil {
			t.Fatal("expected error for non-existing property")
		}
	})
}

func TestNamespaceConstants(t *testing.T) {
	expected := map[string]string{
		"RdfsNS":     "rdfs",
		"RdfNS":      "rdf",
		"CloudNS":    "cloud",
		"CloudRelNS": "cloud-rel",
		"CloudOwlNS": "cloud-owl",
		"XsdNS":      "xsd",
		"NetNS":      "net",
		"NetowlNS":   "net-owl",
	}
	actuals := map[string]string{
		"RdfsNS":     RdfsNS,
		"RdfNS":      RdfNS,
		"CloudNS":    CloudNS,
		"CloudRelNS": CloudRelNS,
		"CloudOwlNS": CloudOwlNS,
		"XsdNS":      XsdNS,
		"NetNS":      NetNS,
		"NetowlNS":   NetowlNS,
	}
	for name, expVal := range expected {
		if actuals[name] != expVal {
			t.Fatalf("%s: got %q, want %q", name, actuals[name], expVal)
		}
	}
}

func TestRDFTermComposition(t *testing.T) {
	if RdfsLabel != "rdfs:label" {
		t.Fatalf("RdfsLabel: got %q, want %q", RdfsLabel, "rdfs:label")
	}
	if RdfType != "rdf:type" {
		t.Fatalf("RdfType: got %q, want %q", RdfType, "rdf:type")
	}
	if RdfProperty != "rdf:Property" {
		t.Fatalf("RdfProperty: got %q, want %q", RdfProperty, "rdf:Property")
	}
	if RdfsSubProperty != "rdfs:subPropertyOf" {
		t.Fatalf("RdfsSubProperty: got %q, want %q", RdfsSubProperty, "rdfs:subPropertyOf")
	}
	if RdfsList != "rdfs:list" {
		t.Fatalf("RdfsList: got %q, want %q", RdfsList, "rdfs:list")
	}
	if RdfsDefinedBy != "rdfs:isDefinedBy" {
		t.Fatalf("RdfsDefinedBy: got %q, want %q", RdfsDefinedBy, "rdfs:isDefinedBy")
	}
	if RdfsDataType != "rdfs:Datatype" {
		t.Fatalf("RdfsDataType: got %q, want %q", RdfsDataType, "rdfs:Datatype")
	}
	if RdfsLiteral != "rdfs:Literal" {
		t.Fatalf("RdfsLiteral: got %q, want %q", RdfsLiteral, "rdfs:Literal")
	}
	if RdfsClass != "rdfs:Class" {
		t.Fatalf("RdfsClass: got %q, want %q", RdfsClass, "rdfs:Class")
	}
	if XsdString != "xsd:string" {
		t.Fatalf("XsdString: got %q, want %q", XsdString, "xsd:string")
	}
	if XsdBoolean != "xsd:boolean" {
		t.Fatalf("XsdBoolean: got %q, want %q", XsdBoolean, "xsd:boolean")
	}
	if XsdInt != "xsd:int" {
		t.Fatalf("XsdInt: got %q, want %q", XsdInt, "xsd:int")
	}
	if XsdDateTime != "xsd:dateTime" {
		t.Fatalf("XsdDateTime: got %q, want %q", XsdDateTime, "xsd:dateTime")
	}
}

func TestRelationConstants(t *testing.T) {
	if ParentOf != "cloud-rel:parentOf" {
		t.Fatalf("ParentOf: got %q, want %q", ParentOf, "cloud-rel:parentOf")
	}
	if ApplyOn != "cloud-rel:applyOn" {
		t.Fatalf("ApplyOn: got %q, want %q", ApplyOn, "cloud-rel:applyOn")
	}
	if ChildrenOfRel != "childrenOf" {
		t.Fatalf("ChildrenOfRel: got %q, want %q", ChildrenOfRel, "childrenOf")
	}
	if DependingOnRel != "dependingOn" {
		t.Fatalf("DependingOnRel: got %q, want %q", DependingOnRel, "dependingOn")
	}
}

func TestClassConstants(t *testing.T) {
	if Grant != "cloud-owl:Grant" {
		t.Fatalf("Grant: got %q, want %q", Grant, "cloud-owl:Grant")
	}
	if NetFirewallRule != "net-owl:FirewallRule" {
		t.Fatalf("NetFirewallRule: got %q, want %q", NetFirewallRule, "net-owl:FirewallRule")
	}
	if Permission != "cloud:permission" {
		t.Fatalf("Permission: got %q, want %q", Permission, "cloud:permission")
	}
}

func TestEmptyRDFProperties(t *testing.T) {
	empty := RDFProperties{}

	_, err := empty.Get("anything")
	if err == nil {
		t.Fatal("expected error from empty RDFProperties.Get")
	}

	if empty.IsRDFProperty("anything") {
		t.Fatal("expected false from empty RDFProperties.IsRDFProperty")
	}

	if empty.IsRDFSubProperty("anything") {
		t.Fatal("expected false from empty RDFProperties.IsRDFSubProperty")
	}

	if empty.IsRDFList("anything") {
		t.Fatal("expected false from empty RDFProperties.IsRDFList")
	}

	_, err = empty.GetDataType("anything")
	if err == nil {
		t.Fatal("expected error from empty RDFProperties.GetDataType")
	}

	_, err = empty.GetLabel("anything")
	if err == nil {
		t.Fatal("expected error from empty RDFProperties.GetLabel")
	}

	_, err = empty.GetDefinedBy("anything")
	if err == nil {
		t.Fatal("expected error from empty RDFProperties.GetDefinedBy")
	}
}

func TestCustomRDFProperties(t *testing.T) {
	props := RDFProperties{
		"test:prop1": {
			ID:            "test:prop1",
			RdfType:       RdfProperty,
			RdfsLabel:     "Prop1",
			RdfsDefinedBy: RdfsLiteral,
			RdfsDataType:  XsdString,
		},
		"test:prop2": {
			ID:            "test:prop2",
			RdfType:       RdfsSubProperty,
			RdfsLabel:     "Prop2",
			RdfsDefinedBy: RdfsList,
			RdfsDataType:  XsdInt,
		},
	}

	// Get
	p, err := props.Get("test:prop1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.RdfsLabel != "Prop1" {
		t.Fatalf("expected 'Prop1', got %q", p.RdfsLabel)
	}

	// IsRDFProperty
	if !props.IsRDFProperty("test:prop1") {
		t.Fatal("expected test:prop1 to be RDF property")
	}
	if props.IsRDFProperty("test:prop2") {
		t.Fatal("expected test:prop2 to not be RDF property")
	}

	// IsRDFSubProperty
	if !props.IsRDFSubProperty("test:prop2") {
		t.Fatal("expected test:prop2 to be RDF sub-property")
	}
	if props.IsRDFSubProperty("test:prop1") {
		t.Fatal("expected test:prop1 to not be RDF sub-property")
	}

	// IsRDFList
	if !props.IsRDFList("test:prop2") {
		t.Fatal("expected test:prop2 to be RDF list")
	}
	if props.IsRDFList("test:prop1") {
		t.Fatal("expected test:prop1 to not be RDF list")
	}

	// GetDataType
	dt, err := props.GetDataType("test:prop1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dt != XsdString {
		t.Fatalf("expected %q, got %q", XsdString, dt)
	}

	dt, err = props.GetDataType("test:prop2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dt != XsdInt {
		t.Fatalf("expected %q, got %q", XsdInt, dt)
	}

	// GetLabel
	label, err := props.GetLabel("test:prop2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if label != "Prop2" {
		t.Fatalf("expected 'Prop2', got %q", label)
	}

	// GetDefinedBy
	def, err := props.GetDefinedBy("test:prop2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if def != RdfsList {
		t.Fatalf("expected %q, got %q", RdfsList, def)
	}
}
