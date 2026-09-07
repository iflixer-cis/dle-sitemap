package sitemap

import "testing"

func TestFileNamesExcludesEmptySitemap(t *testing.T) {
	sm := &SmSitemap{}
	if err := sm.Init("example.com", t.TempDir(), "sitemap_collections.xml"); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := sm.Close(); err != nil {
			t.Errorf("close sitemap: %v", err)
		}
	})

	if names := sm.FileNames(); len(names) != 0 {
		t.Fatalf("empty sitemap returned filenames: %v", names)
	}

	if err := sm.Add(SmSitemapRow{Loc: "https://example.com/collection"}); err != nil {
		t.Fatal(err)
	}

	names := sm.FileNames()
	if len(names) != 1 || names[0] != "sitemap_collections.xml" {
		t.Fatalf("non-empty sitemap returned unexpected filenames: %v", names)
	}
}
