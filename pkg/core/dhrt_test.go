package dhrt

import (
	"testing"
)

func TestSet(t *testing.T) {
	db, err := Open(Config{})
	if err != nil {
		t.Fail()
	}
	key1 := "key1"
	dt, err := db.Set(key1, "one two three")
	if err != nil {
		t.Fail()
	}
	if dt.Key != key1 {
		t.Fatalf("Expected %s to be %s", dt.Key, key1)
	}
}

func TestGet(t *testing.T) {
	db, err := Open(Config{})
	if err != nil {
		t.Fail()
	}
	key1 := "key1"
	value1 := "this is sparta"
	_, _ = db.Set(key1, value1)
	out, err := db.Get(key1)
	if err != nil {
		t.Fail()
	}
	if out.Value != value1 {
		t.Fail()
	}
	if out.Created == 0 {
		t.Fail()
	}
	if out.Created != out.Updated {
		t.Fail()
	}
}

func TestSetUUID(t *testing.T) {
	db, err := Open(Config{})
	if err != nil {
		t.Fail()
	}
	key1 := ""
	value1 := "this is sparta"
	dt, _ := db.Set(key1, value1)
	out, err := db.Get(dt.Key)
	if err != nil {
		t.Fail()
	}
	if out.Value != value1 {
		t.Fail()
	}
}

func TestDel(t *testing.T) {
	db, _ := Open(Config{})
	dt, _ := db.Set("", "this is sparta")
	out, _ := db.Get(dt.Key)
	if out.Value == "" {
		t.Fail()
	}
	err := db.Del(out.Key)
	if err != nil {
		t.Fail()
	}
	out1, err := db.Get(dt.Key)
	if out1.Value != "" {
		t.Fail()
	}
}

func TestDelAll(t *testing.T) {
	db, _ := Open(Config{})
	dt1, _ := db.Set("key1", "this is sparta")
	dt2, _ := db.Set("key2", "this is not sparta")

	err := db.DelAll([]string{dt1.Key, dt2.Key})
	if err != nil {
		t.Fail()
	}
	out1, err := db.Get(dt1.Key)
	if out1.Value != "" {
		t.Fail()
	}
	out2, err := db.Get(dt2.Key)
	if out2.Value != "" {
		t.Fail()
	}
}

func TestUpdate(t *testing.T) {
	db, err := Open(Config{})
	if err != nil {
		t.Fail()
	}
	key1 := "key1"
	value1 := "this is sparta"
	_, _ = db.Set(key1, value1)
	out, err := db.Get(key1)
	if out.Value != value1 {
		t.Fail()
	}
	value2 := "echo..."
	_, _ = db.Set(key1, value2)
	out, err = db.Get(key1)
	if out.Value != value2 {
		t.Fail()
	}
	if out.Created == out.Updated {
		t.Fail()
	}
}

func TestStats(t *testing.T) {
	db, _ := Open(Config{})
	stats := db.Stats()
	if stats.ItemCount != 0 {
		t.Fail()
	}
	_, _ = db.Set("wat", "lol")
	stats1 := db.Stats()
	if stats1.ItemCount != 1 {
		t.Fail()
	}
}

func TestClose(t *testing.T) {
	db, _ := Open(Config{})
	value := "value"
	dt, _ := db.Set("wat", value)
	out, _ := db.Get(dt.Key)
	if !db.Open {
		t.Fail()
	}
	if out.Value != value {
		t.Fail()
	}
	db.Close()
	out1, _ := db.Get(dt.Key)
	if out1.Value != "" {
		t.Fail()
	}
	if db.Stats().Size != 0 {
		t.Fail()
	}
	if db.Open {
		t.Fail()
	}
}
