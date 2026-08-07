/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package match

import (
	"testing"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/graph/resourcetest"
)

func TestMatchers(t *testing.T) {
	tcases := []struct {
		match    cloud.Matcher
		resource cloud.Resource
		expect   bool
	}{
		{match: Property("Inexisting", "empty"), resource: resourcetest.Instance("i1").Build(), expect: false},
		{match: Property("Prop", "value"), resource: resourcetest.Instance("i1").Prop("Prop", "value").Build(), expect: true},
		{match: And(Property("Inexisting", "empty"), Property("Prop", "value")), resource: resourcetest.Instance("i1").Prop("Prop", "value").Build(), expect: false},
		{match: Or(Property("Inexisting", "empty"), Property("Prop", "value")), resource: resourcetest.Instance("i1").Prop("Prop", "value").Build(), expect: true},
		{match: Or(Property("Inexisting1", ""), Property("Inexisting2", "")), resource: resourcetest.Instance("i1").Build(), expect: false},
		{match: And(Property("Prop1", "value1"), Property("Prop2", "value2")), resource: resourcetest.Instance("i1").Prop("Prop1", "value1").Prop("Prop2", "value2").Build(), expect: true},
		{match: And(Property("Prop1", "value1"), Property("Prop2", "value2")), resource: resourcetest.Instance("i1").Prop("Prop1", "value1").Prop("Prop2", "value2").Build(), expect: true},
		{match: Property("Prop", 42).MatchString(), resource: resourcetest.Instance("i1").Prop("Prop", "42").Build(), expect: true},
		{match: Property("Prop", "WithCase").IgnoreCase(), resource: resourcetest.Instance("i1").Prop("Prop", "WITHCASE").Build(), expect: true},
		{match: Property("Prop", "42").IgnoreCase().MatchString(), resource: resourcetest.Instance("i1").Prop("Prop", 42).Build(), expect: true},
		{match: Property("Prop", "inside").Contains(), resource: resourcetest.Instance("i1").Prop("Prop", "Match inside the content").Build(), expect: true},
		{match: Tag("Key", "Val"), resource: resourcetest.Instance("i1").Prop("Tags", []string{"Key=Val"}).Build(), expect: true},
		{match: Tag("Key", "Notthis"), resource: resourcetest.Instance("i1").Prop("Tags", []string{"Key=Val"}).Build(), expect: false},
		{match: TagKey("Key"), resource: resourcetest.Instance("i1").Prop("Tags", []string{"Key=Val"}).Build(), expect: true},
		{match: TagKey("NotThis"), resource: resourcetest.Instance("i1").Prop("Tags", []string{"Key=Val"}).Build(), expect: false},
		{match: TagValue("Val"), resource: resourcetest.Instance("i1").Prop("Tags", []string{"Key=Val"}).Build(), expect: true},
		{match: TagValue("NotThis"), resource: resourcetest.Instance("i1").Prop("Tags", []string{"Key=Val"}).Build(), expect: false},
	}
	for i, tcase := range tcases {
		if got, want := tcase.match.Match(tcase.resource), tcase.expect; got != want {
			t.Fatalf("%d: got %t, want %t", i+1, got, want)
		}
	}
}

func TestAndCombinator(t *testing.T) {
	res := resourcetest.Instance("i1").Prop("Name", "web-server").Prop("State", "running").Build()

	t.Run("all matchers true returns true", func(t *testing.T) {
		m := And(Property("Name", "web-server"), Property("State", "running"))
		if !m.Match(res) {
			t.Fatal("expected And with all true matchers to return true")
		}
	})

	t.Run("one false matcher returns false", func(t *testing.T) {
		m := And(Property("Name", "web-server"), Property("State", "stopped"))
		if m.Match(res) {
			t.Fatal("expected And with one false matcher to return false")
		}
	})

	t.Run("first matcher false short-circuits", func(t *testing.T) {
		m := And(Property("State", "stopped"), Property("Name", "web-server"))
		if m.Match(res) {
			t.Fatal("expected And to return false when first matcher is false")
		}
	})

	t.Run("empty matchers returns false", func(t *testing.T) {
		m := And()
		if m.Match(res) {
			t.Fatal("expected And with no matchers to return false")
		}
	})

	t.Run("single matcher delegates correctly", func(t *testing.T) {
		m := And(Property("Name", "web-server"))
		if !m.Match(res) {
			t.Fatal("expected And with single true matcher to return true")
		}
		m = And(Property("Name", "wrong"))
		if m.Match(res) {
			t.Fatal("expected And with single false matcher to return false")
		}
	})

	t.Run("three matchers all true", func(t *testing.T) {
		res3 := resourcetest.Instance("i2").Prop("A", "1").Prop("B", "2").Prop("C", "3").Build()
		m := And(Property("A", "1"), Property("B", "2"), Property("C", "3"))
		if !m.Match(res3) {
			t.Fatal("expected And with three true matchers to return true")
		}
	})

	t.Run("three matchers middle false", func(t *testing.T) {
		res3 := resourcetest.Instance("i2").Prop("A", "1").Prop("B", "2").Prop("C", "3").Build()
		m := And(Property("A", "1"), Property("B", "wrong"), Property("C", "3"))
		if m.Match(res3) {
			t.Fatal("expected And to return false when middle matcher fails")
		}
	})
}

func TestOrCombinator(t *testing.T) {
	res := resourcetest.Instance("i1").Prop("Name", "web-server").Prop("State", "running").Build()

	t.Run("all matchers false returns false", func(t *testing.T) {
		m := Or(Property("Name", "wrong"), Property("State", "stopped"))
		if m.Match(res) {
			t.Fatal("expected Or with all false matchers to return false")
		}
	})

	t.Run("one true matcher returns true", func(t *testing.T) {
		m := Or(Property("Name", "wrong"), Property("State", "running"))
		if !m.Match(res) {
			t.Fatal("expected Or with one true matcher to return true")
		}
	})

	t.Run("all matchers true returns true", func(t *testing.T) {
		m := Or(Property("Name", "web-server"), Property("State", "running"))
		if !m.Match(res) {
			t.Fatal("expected Or with all true matchers to return true")
		}
	})

	t.Run("first matcher true short-circuits", func(t *testing.T) {
		m := Or(Property("Name", "web-server"), Property("Missing", "value"))
		if !m.Match(res) {
			t.Fatal("expected Or to return true when first matcher is true")
		}
	})

	t.Run("empty matchers returns false", func(t *testing.T) {
		m := Or()
		if m.Match(res) {
			t.Fatal("expected Or with no matchers to return false")
		}
	})

	t.Run("single matcher delegates correctly", func(t *testing.T) {
		m := Or(Property("Name", "web-server"))
		if !m.Match(res) {
			t.Fatal("expected Or with single true matcher to return true")
		}
		m = Or(Property("Name", "wrong"))
		if m.Match(res) {
			t.Fatal("expected Or with single false matcher to return false")
		}
	})

	t.Run("three matchers last true", func(t *testing.T) {
		m := Or(Property("Name", "wrong"), Property("State", "stopped"), Property("State", "running"))
		if !m.Match(res) {
			t.Fatal("expected Or to return true when last matcher matches")
		}
	})
}

func TestPropertyMatcher(t *testing.T) {
	t.Run("exact string match", func(t *testing.T) {
		res := resourcetest.Instance("i1").Prop("Name", "my-instance").Build()
		if !Property("Name", "my-instance").Match(res) {
			t.Fatal("expected exact string match to succeed")
		}
	})

	t.Run("string no match", func(t *testing.T) {
		res := resourcetest.Instance("i1").Prop("Name", "my-instance").Build()
		if Property("Name", "other").Match(res) {
			t.Fatal("expected non-matching string to fail")
		}
	})

	t.Run("missing property returns false", func(t *testing.T) {
		res := resourcetest.Instance("i1").Build()
		if Property("Missing", "value").Match(res) {
			t.Fatal("expected missing property to return false")
		}
	})

	t.Run("integer property match", func(t *testing.T) {
		res := resourcetest.Instance("i1").Prop("Port", 8080).Build()
		if !Property("Port", 8080).Match(res) {
			t.Fatal("expected integer property match to succeed")
		}
	})

	t.Run("integer property no match", func(t *testing.T) {
		res := resourcetest.Instance("i1").Prop("Port", 8080).Build()
		if Property("Port", 9090).Match(res) {
			t.Fatal("expected integer mismatch to fail")
		}
	})

	t.Run("bool property match", func(t *testing.T) {
		res := resourcetest.Instance("i1").Prop("Public", true).Build()
		if !Property("Public", true).Match(res) {
			t.Fatal("expected bool property match to succeed")
		}
	})

	t.Run("bool property no match", func(t *testing.T) {
		res := resourcetest.Instance("i1").Prop("Public", true).Build()
		if Property("Public", false).Match(res) {
			t.Fatal("expected bool mismatch to fail")
		}
	})

	t.Run("type mismatch without MatchString", func(t *testing.T) {
		res := resourcetest.Instance("i1").Prop("Port", 8080).Build()
		if Property("Port", "8080").Match(res) {
			t.Fatal("expected type mismatch (int vs string) to fail without MatchString")
		}
	})

	t.Run("MatchString converts types for comparison", func(t *testing.T) {
		res := resourcetest.Instance("i1").Prop("Port", 8080).Build()
		if !Property("Port", "8080").MatchString().Match(res) {
			t.Fatal("expected MatchString to convert and match")
		}
	})

	t.Run("IgnoreCase match", func(t *testing.T) {
		res := resourcetest.Instance("i1").Prop("Name", "MyInstance").Build()
		if !Property("Name", "myinstance").IgnoreCase().Match(res) {
			t.Fatal("expected case-insensitive match to succeed")
		}
	})

	t.Run("IgnoreCase no match", func(t *testing.T) {
		res := resourcetest.Instance("i1").Prop("Name", "MyInstance").Build()
		if Property("Name", "other").IgnoreCase().Match(res) {
			t.Fatal("expected case-insensitive non-match to fail")
		}
	})

	t.Run("Contains match", func(t *testing.T) {
		res := resourcetest.Instance("i1").Prop("Name", "my-web-server-prod").Build()
		if !Property("Name", "web-server").Contains().Match(res) {
			t.Fatal("expected Contains match to succeed")
		}
	})

	t.Run("Contains no match", func(t *testing.T) {
		res := resourcetest.Instance("i1").Prop("Name", "my-web-server-prod").Build()
		if Property("Name", "database").Contains().Match(res) {
			t.Fatal("expected Contains non-match to fail")
		}
	})

	t.Run("Contains with IgnoreCase", func(t *testing.T) {
		res := resourcetest.Instance("i1").Prop("Name", "My-Web-Server").Build()
		if !Property("Name", "web-server").Contains().IgnoreCase().Match(res) {
			t.Fatal("expected Contains with IgnoreCase to succeed")
		}
	})

	t.Run("empty string property", func(t *testing.T) {
		res := resourcetest.Instance("i1").Prop("Name", "").Build()
		if !Property("Name", "").Match(res) {
			t.Fatal("expected empty string match to succeed")
		}
	})

	t.Run("slice property exact match", func(t *testing.T) {
		res := resourcetest.Instance("i1").Prop("IPs", []string{"10.0.0.1", "10.0.0.2"}).Build()
		if !Property("IPs", []string{"10.0.0.1", "10.0.0.2"}).Match(res) {
			t.Fatal("expected slice property match to succeed")
		}
	})
}

func TestTagMatcher(t *testing.T) {
	t.Run("exact key value match", func(t *testing.T) {
		res := resourcetest.Instance("i1").Prop("Tags", []string{"env=prod"}).Build()
		if !Tag("env", "prod").Match(res) {
			t.Fatal("expected exact tag match to succeed")
		}
	})

	t.Run("wrong value", func(t *testing.T) {
		res := resourcetest.Instance("i1").Prop("Tags", []string{"env=prod"}).Build()
		if Tag("env", "staging").Match(res) {
			t.Fatal("expected wrong tag value to fail")
		}
	})

	t.Run("missing tag key", func(t *testing.T) {
		res := resourcetest.Instance("i1").Prop("Tags", []string{"env=prod"}).Build()
		if Tag("team", "backend").Match(res) {
			t.Fatal("expected missing tag key to fail")
		}
	})

	t.Run("empty tags", func(t *testing.T) {
		res := resourcetest.Instance("i1").Prop("Tags", []string{}).Build()
		if Tag("env", "prod").Match(res) {
			t.Fatal("expected empty tags to fail")
		}
	})

	t.Run("no tags property", func(t *testing.T) {
		res := resourcetest.Instance("i1").Build()
		if Tag("env", "prod").Match(res) {
			t.Fatal("expected missing Tags property to fail")
		}
	})

	t.Run("multiple tags match second", func(t *testing.T) {
		res := resourcetest.Instance("i1").Prop("Tags", []string{"env=prod", "team=backend"}).Build()
		if !Tag("team", "backend").Match(res) {
			t.Fatal("expected match on second tag to succeed")
		}
	})

	t.Run("wildcard in value", func(t *testing.T) {
		res := resourcetest.Instance("i1").Prop("Tags", []string{"Name=web-server-01"}).Build()
		if !Tag("Name", "web-server-*").Match(res) {
			t.Fatal("expected wildcard tag match to succeed")
		}
	})

	t.Run("wildcard no match", func(t *testing.T) {
		res := resourcetest.Instance("i1").Prop("Tags", []string{"Name=db-server-01"}).Build()
		if Tag("Name", "web-*").Match(res) {
			t.Fatal("expected wildcard non-match to fail")
		}
	})

	t.Run("tags property wrong type", func(t *testing.T) {
		res := resourcetest.Instance("i1").Prop("Tags", "not-a-slice").Build()
		if Tag("env", "prod").Match(res) {
			t.Fatal("expected non-slice Tags to fail")
		}
	})
}

func TestTagKeyMatcher(t *testing.T) {
	t.Run("key exists", func(t *testing.T) {
		res := resourcetest.Instance("i1").Prop("Tags", []string{"env=prod"}).Build()
		if !TagKey("env").Match(res) {
			t.Fatal("expected existing key to match")
		}
	})

	t.Run("key does not exist", func(t *testing.T) {
		res := resourcetest.Instance("i1").Prop("Tags", []string{"env=prod"}).Build()
		if TagKey("team").Match(res) {
			t.Fatal("expected missing key to fail")
		}
	})

	t.Run("no tags property", func(t *testing.T) {
		res := resourcetest.Instance("i1").Build()
		if TagKey("env").Match(res) {
			t.Fatal("expected missing Tags property to fail")
		}
	})

	t.Run("empty tags", func(t *testing.T) {
		res := resourcetest.Instance("i1").Prop("Tags", []string{}).Build()
		if TagKey("env").Match(res) {
			t.Fatal("expected empty tags to fail")
		}
	})

	t.Run("multiple tags match second key", func(t *testing.T) {
		res := resourcetest.Instance("i1").Prop("Tags", []string{"env=prod", "team=backend"}).Build()
		if !TagKey("team").Match(res) {
			t.Fatal("expected match on second tag key to succeed")
		}
	})

	t.Run("tags property wrong type", func(t *testing.T) {
		res := resourcetest.Instance("i1").Prop("Tags", 42).Build()
		if TagKey("env").Match(res) {
			t.Fatal("expected non-slice Tags to fail")
		}
	})
}

func TestTagValueMatcher(t *testing.T) {
	t.Run("value exists in tag", func(t *testing.T) {
		res := resourcetest.Instance("i1").Prop("Tags", []string{"env=prod"}).Build()
		if !TagValue("prod").Match(res) {
			t.Fatal("expected existing value to match")
		}
	})

	t.Run("value does not exist", func(t *testing.T) {
		res := resourcetest.Instance("i1").Prop("Tags", []string{"env=prod"}).Build()
		if TagValue("staging").Match(res) {
			t.Fatal("expected missing value to fail")
		}
	})

	t.Run("no tags property", func(t *testing.T) {
		res := resourcetest.Instance("i1").Build()
		if TagValue("prod").Match(res) {
			t.Fatal("expected missing Tags property to fail")
		}
	})

	t.Run("empty tags", func(t *testing.T) {
		res := resourcetest.Instance("i1").Prop("Tags", []string{}).Build()
		if TagValue("prod").Match(res) {
			t.Fatal("expected empty tags to fail")
		}
	})

	t.Run("multiple tags match value in second tag", func(t *testing.T) {
		res := resourcetest.Instance("i1").Prop("Tags", []string{"env=prod", "team=backend"}).Build()
		if !TagValue("backend").Match(res) {
			t.Fatal("expected match on second tag value to succeed")
		}
	})

	t.Run("tags property wrong type", func(t *testing.T) {
		res := resourcetest.Instance("i1").Prop("Tags", "not-a-slice").Build()
		if TagValue("prod").Match(res) {
			t.Fatal("expected non-slice Tags to fail")
		}
	})

	t.Run("tag without equals sign does not match", func(t *testing.T) {
		res := resourcetest.Instance("i1").Prop("Tags", []string{"standalone"}).Build()
		if TagValue("standalone").Match(res) {
			t.Fatal("expected tag without = to not match as a value")
		}
	})
}

func TestNestedCombinators(t *testing.T) {
	res := resourcetest.Instance("i1").
		Prop("Name", "web-server").
		Prop("State", "running").
		Prop("Tags", []string{"env=prod", "team=backend"}).
		Build()

	t.Run("And nested in Or", func(t *testing.T) {
		m := Or(
			And(Property("Name", "db-server"), Property("State", "running")),
			And(Property("Name", "web-server"), Property("State", "running")),
		)
		if !m.Match(res) {
			t.Fatal("expected nested And-in-Or to match")
		}
	})

	t.Run("Or nested in And", func(t *testing.T) {
		m := And(
			Or(Property("Name", "web-server"), Property("Name", "db-server")),
			TagKey("env"),
		)
		if !m.Match(res) {
			t.Fatal("expected nested Or-in-And to match")
		}
	})

	t.Run("nested combinators no match", func(t *testing.T) {
		m := And(
			Or(Property("Name", "api-server"), Property("Name", "db-server")),
			TagKey("env"),
		)
		if m.Match(res) {
			t.Fatal("expected nested combinators to fail when inner Or fails")
		}
	})
}
