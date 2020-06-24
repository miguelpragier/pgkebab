package pgkebab

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"
)

var db *DBLink

func TestConnect(t *testing.T) {
	var cs *ConnectionString

	if os.Getenv("{APPDBCS}") != "" {
		cs = ConnStringEnvVar("{APPDBCS}")
	} else {
		cs = ConnStringDirect("postgres://pgkebab:001001001001@localhost/pgkebab?sslmode=disable")
	}

	opts := Options(cs, 10, 10, 5, 5, 10, true)

	if _db, err := NewConnected(opts); err != nil {
		t.Fatal(err)
	} else {
		db = _db
	}
}

func TestExec(t *testing.T) {
	if _, err := db.Exec("DROP TABLE IF EXISTS kbd"); err != nil {
		t.Fatal(err)
	} else {
		t.Log("table droped")
	}

	if _, err := db.Exec("CREATE TABLE kbd (id SERIAL NOT NULL PRIMARY KEY,txt VARCHAR(30), ts TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, dec DECIMAL, x INT, j JSONB)"); err != nil {
		t.Fatal(err)
	} else {
		t.Log("table kbd created")
	}

	if i, err := db.Exec("INSERT INTO kbd (txt, dec, x) VALUES ('test one',1.5,8),('test two',2.6,90),('test three',3.1,120),('test four',4.6,311),('test five',5.10,501)"); err != nil {
		t.Fatal(err)
	} else {
		if i != 5 {
			t.Logf("%d records inserted\n", i)
		}
	}
}

func TestGetString(t *testing.T) {
	if s, err := db.GetString("SELECT CURRENT_TIMESTAMP now"); err != nil {
		t.Fatal(err)
	} else {
		t.Log(s)
	}
}

func TestGetInt64(t *testing.T) {
	if i, err := db.GetInt64("SELECT COUNT(*) FROM kbd WHERE x>=120"); err != nil {
		t.Fatal(err)
	} else {
		if i != 3 {
			t.Fatalf("expected 3, got %d", i)
		}
	}
}

func TestGetFloat64(t *testing.T) {
	if f, err := db.GetFloat64("SELECT SUM(dec) FROM kbd"); err != nil {
		t.Fatal(err)
	} else {
		if f != 16.9 {
			t.Fatalf("expected 16.9, got %f", f)
		}
	}
}

func TestGetBool(t *testing.T) {
	if b, err := db.GetBool("SELECT EXISTS ( SELECT*FROM kbd)"); err != nil {
		t.Fatal(err)
	} else {
		if !b {
			t.Fatalf("'true' expected, got false")
		}
	}
}

func TestDBLink_Insert(t *testing.T) {
	pairs := map[string]interface{}{
		"txt": "testing Insert()",
		"dec": 88.8888888,
		"x":   time.Now().Unix(),
	}

	if err := db.Insert("kbd", pairs); err != nil {
		t.Fatal(err)
	}
}

func TestDBLink_InsertID(t *testing.T) {
	pairs := map[string]interface{}{
		"txt": "testing Insert()",
		"dec": 99.9999999,
		"x":   time.Now().Unix(),
	}

	if id, err := db.InsertID("kbd", pairs); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("last inserted id: %d\n", id)
	}
}

func TestDBLink_GetOne(t *testing.T) {
	if r, err := db.GetOne("SELECT*FROM kbd WHERE true LIMIT 1"); err != nil {
		t.Fatal(err)
	} else {
		if r.Ready() {
			t.Log(r.JSON())
		} else {
			fmt.Println("empty row")
		}
	}
}

func TestDBLink_GetMany(t *testing.T) {
	if rs, err := db.GetMany("SELECT*FROM kbd"); err != nil {
		t.Fatal(err)
	} else {
		for rs.Next() {
			row, _ := rs.Row()

			fmt.Printf("id:%d\ttxt:%s\tdec:%f\tx:%d\n", row.Int64("id"), row.String("txt"), row.Float64("dec"), row.Int64("x"))
		}

		t.Logf(rs.JSON(false))
	}
}

func TestDBLink_MustGetMany(t *testing.T) {
	if _, err := db.MustGetMany("SELECT*FROM kbd WHERE id<0"); err == nil {
		t.Fatal(err)
	} else {
		t.Log("Super ok!")
	}
}

func TestDBLink_InsertJSONMap(t *testing.T) {
	d, err0 := json.Marshal(map[string]interface{}{"xyz": 987, "b": true, "i": 987654321})

	if err0 != nil {
		t.Fatal(err0)
	}

	pairs := map[string]interface{}{
		"txt": "TestDBLink_InsertJSONMap",
		"dec": 0,
		"x":   0,
		"J":   string(d),
	}

	if err := db.Insert("kbd", pairs); err != nil {
		t.Fatal(err)
	}
}

func TestDBLink_GetJSONMap(t *testing.T) {
	if m, err := db.GetJSONMap("SELECT j FROM kbd WHERE txt='TestDBLink_InsertJSONMap' LIMIT 1"); err != nil {
		t.Error(err)
	} else {
		t.Log(m)
	}
}

type jstruct struct {
	ABC int  `json:"abc"`
	B   bool `json:"b"`
	I   int  `json:"i"`
}

func TestDBLink_InsertJSONStruct(t *testing.T) {
	x := jstruct{ABC: 123, B: true, I: 9}

	d, err0 := json.Marshal(x)

	if err0 != nil {
		t.Fatal(err0)
	}

	pairs := map[string]interface{}{
		"txt": "TestDBLink_InsertJSONStruct",
		"dec": 0,
		"x":   0,
		"j":   string(d),
	}

	if err := db.Insert("kbd", pairs); err != nil {
		t.Fatal(err)
	}
}

func TestDBLink_GetJSONStruct(t *testing.T) {
	var j jstruct

	if err := db.GetJSONStruct(&j, "SELECT j FROM kbd WHERE txt='TestDBLink_InsertJSONStruct' LIMIT 1"); err != nil {
		t.Error(err)
	} else {
		t.Log(j)
	}
}

func Test_RowJSONMap(t *testing.T) {
	if rs, err := db.GetOne("SELECT j FROM kbd WHERE txt='TestDBLink_InsertJSONMap' LIMIT 1"); err != nil {
		t.Error(err)
	} else {
		t.Log(rs.JSONMap("j"))
	}
}

func Test_RowJSONStruct(t *testing.T) {
	if rs, err := db.GetOne("SELECT j FROM kbd WHERE txt='TestDBLink_InsertJSONStruct' LIMIT 1"); err != nil {
		t.Error(err)
	} else {
		var j jstruct
		if err0 := rs.JSONStruct("j", &j); err0 != nil {
			t.Error(err0)

			return
		}

		t.Log(j)
	}
}
