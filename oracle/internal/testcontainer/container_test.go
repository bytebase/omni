package testcontainer

import (
	"testing"

	"github.com/bytebase/omni/enginetest"
)

func TestSameServerPinAllowsDifferentRoleAndJudge(t *testing.T) {
	base := enginetest.Pin{
		Role:   "default",
		Image:  "gvenzl/oracle-free:23.26.3-slim-faststart",
		Digest: "sha256:f5ff19033860d662c821cb04eb10483fa94f14f78eae252d054291ea07028093",
		Tag:    "23.26.3.0.0",
	}
	parserPin := base
	parserPin.Role = "parser"
	parserPin.Judge = "4"

	if !sameServerPin(base, parserPin) {
		t.Fatalf("expected sameServerPin(%v, %v) = true", base, parserPin)
	}

	diffImage := base
	diffImage.Image = "gvenzl/oracle-free:23.99.0-slim-faststart"
	if sameServerPin(base, diffImage) {
		t.Fatalf("expected different image to return false")
	}

	diffDigest := base
	diffDigest.Digest = "sha256:0000"
	if sameServerPin(base, diffDigest) {
		t.Fatalf("expected different digest to return false")
	}

	diffTag := base
	diffTag.Tag = "23.99.0.0.0"
	if sameServerPin(base, diffTag) {
		t.Fatalf("expected different tag to return false")
	}
}
