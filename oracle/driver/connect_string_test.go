package driver

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseConnectString(t *testing.T) {
	ep := func(host, port string) Endpoint { return Endpoint{Host: host, Port: port} }
	svc := func(service string, endpoints ...Endpoint) ConnectTarget {
		return ConnectTarget{Endpoints: endpoints, Service: service}
	}
	sid := func(sid string, endpoints ...Endpoint) ConnectTarget {
		return ConnectTarget{Endpoints: endpoints, SID: sid}
	}

	tests := []struct {
		connectString string
		want          ConnectTarget
		ok            bool
	}{
		// Easy Connect forms.
		{"localhost:1521/FREEPDB1", svc("FREEPDB1", ep("localhost", "1521")), true},
		{"//db.example.com:1522/ORCL", svc("ORCL", ep("db.example.com", "1522")), true},
		{"DB.EXAMPLE.COM/orcl", svc("orcl", ep("db.example.com", "1521")), true},
		{"host:1521/svc:dedicated/inst1", svc("svc", ep("host", "1521")), true},
		{"[::1]:1521/svc", svc("svc", ep("::1", "1521")), true},
		{"[::1]/svc", svc("svc", ep("::1", "1521")), true},
		{"  host:1521/svc  ", svc("svc", ep("host", "1521")), true},

		// Easy Connect Plus: scheme, parameters, several addresses.
		{"tcp://db.example.com:1521/orcl", svc("orcl", ep("db.example.com", "1521")), true},
		{"TCPS://adb.example.com:1522/abc_high.adb.oraclecloud.com", svc("abc_high.adb.oraclecloud.com", ep("adb.example.com", "1522")), true},
		{"tcps://adb.example.com:1522/abc_high.adb.oraclecloud.com", svc("abc_high.adb.oraclecloud.com", ep("adb.example.com", "1522")), true},
		{"//host:1521/svc?connect_timeout=5&retry_count=2", svc("svc", ep("host", "1521")), true},
		{"rac1:1521,rac2:1522/svc", svc("svc", ep("rac1", "1521"), ep("rac2", "1522")), true},
		// A port applies to the hosts listed before it.
		{"rac1,rac2:1522/svc", svc("svc", ep("rac1", "1522"), ep("rac2", "1522")), true},
		{"rac1,rac2:1522,rac3/svc", svc("svc", ep("rac1", "1522"), ep("rac2", "1522"), ep("rac3", "1521")), true},
		{"[::1],[::2]:1522/svc", svc("svc", ep("::1", "1522"), ep("::2", "1522")), true},

		// TNS descriptors.
		{"(DESCRIPTION=(ADDRESS=(PROTOCOL=TCP)(HOST=db.example.com)(PORT=1521))(CONNECT_DATA=(SERVICE_NAME=orcl)))", svc("orcl", ep("db.example.com", "1521")), true},
		{"(DESCRIPTION=(ADDRESS=(PROTOCOL=TCP)(HOST=db)(PORT=1522))(CONNECT_DATA=(SID=orcl)))", sid("orcl", ep("db", "1522")), true},
		{"(DESCRIPTION=(ADDRESS=(PROTOCOL=TCP)(HOST=db))(CONNECT_DATA=(SERVICE_NAME=orcl)))", svc("orcl", ep("db", "1521")), true},
		{"( DESCRIPTION = ( ADDRESS = ( HOST = db ) ( PORT = 1521 ) ) ( CONNECT_DATA = ( SERVICE_NAME = orcl ) ) )", svc("orcl", ep("db", "1521")), true},
		{"(DESCRIPTION=(ADDRESS=(PROTOCOL=TCP)(HOST=\"DB.Example.com\")(PORT=\"1521\"))(CONNECT_DATA=(SERVICE_NAME=\"orcl\")))", svc("orcl", ep("db.example.com", "1521")), true},
		{"(DESCRIPTION=(ADDRESS=(PROTOCOL=TCP)(HOST='db.example.com')(PORT='1522'))(CONNECT_DATA=(SID='orcl')))", sid("orcl", ep("db.example.com", "1522")), true},
		{"(DESCRIPTION_LIST=(DESCRIPTION=(ADDRESS=(HOST=\"h1\")(PORT=1521))(CONNECT_DATA=(SERVICE_NAME=\"proddb\")))(DESCRIPTION=(ADDRESS=(HOST=h2)(PORT=1521))(CONNECT_DATA=(SERVICE_NAME=PRODDB))))", svc("proddb", ep("h1", "1521"), ep("h2", "1521")), true},
		{"(DESCRIPTION=(SOURCE_ROUTE=no)(ADDRESS=(HOST=db)(PORT=1521))(CONNECT_DATA=(SERVICE_NAME=orcl)))", svc("orcl", ep("db", "1521")), true},
		// RAC and Data Guard: several addresses, one service.
		{"(DESCRIPTION=(ADDRESS_LIST=(LOAD_BALANCE=on)(ADDRESS=(PROTOCOL=TCP)(HOST=rac1-vip)(PORT=1521))(ADDRESS=(PROTOCOL=TCP)(HOST=rac2-vip)(PORT=1521)))(CONNECT_DATA=(SERVICE_NAME=proddb)))", svc("proddb", ep("rac1-vip", "1521"), ep("rac2-vip", "1521")), true},
		// A DESCRIPTION_LIST repeating one service or SID is one database (Data Guard failover).
		{"(DESCRIPTION_LIST=(DESCRIPTION=(ADDRESS=(HOST=h1)(PORT=1521))(CONNECT_DATA=(SERVICE_NAME=proddb)))(DESCRIPTION=(ADDRESS=(HOST=h2)(PORT=1521))(CONNECT_DATA=(SERVICE_NAME=PRODDB))))", svc("proddb", ep("h1", "1521"), ep("h2", "1521")), true},
		{"(DESCRIPTION_LIST=(DESCRIPTION=(ADDRESS=(HOST=h1)(PORT=1521))(CONNECT_DATA=(SID=orcl)))(DESCRIPTION=(ADDRESS=(HOST=h2)(PORT=1521))(CONNECT_DATA=(SID=ORCL))))", sid("orcl", ep("h1", "1521"), ep("h2", "1521")), true},

		// A TNS alias names nothing on its own.
		{"PRODDB", ConnectTarget{}, false},
		{"", ConnectTarget{}, false},
		{"   ", ConnectTarget{}, false},

		// Unsupported schemes.
		{"http://host:1521/svc", ConnectTarget{}, false},
		{"ldap://host:389/svc", ConnectTarget{}, false},
		{"tcps:////host:1522/svc", ConnectTarget{}, false},

		// No service or SID.
		{"host:1521", ConnectTarget{}, false},
		{"host:1521/", ConnectTarget{}, false},
		{"host:1521/?connect_timeout=5", ConnectTarget{}, false},
		{"/svc", ConnectTarget{}, false},
		{"(DESCRIPTION=(ADDRESS=(HOST=a)(PORT=1521)))", ConnectTarget{}, false},

		// SOURCE_ROUTE hop chain passes through a proxy; its first address is not the database.
		{"(DESCRIPTION=(SOURCE_ROUTE=yes)(ADDRESS=(HOST=cman)(PORT=1521))(ADDRESS=(HOST=real)(PORT=1521))(CONNECT_DATA=(SERVICE_NAME=svc)))", ConnectTarget{}, false},
		{"(DESCRIPTION=(SOURCE_ROUTE=ON)(ADDRESS=(HOST=cman)(PORT=1521))(ADDRESS=(HOST=real)(PORT=1521))(CONNECT_DATA=(SERVICE_NAME=svc)))", ConnectTarget{}, false},
		{"(DESCRIPTION=(SOURCE_ROUTE=true)(ADDRESS=(HOST=cman)(PORT=1521))(ADDRESS=(HOST=real)(PORT=1521))(CONNECT_DATA=(SERVICE_NAME=svc)))", ConnectTarget{}, false},
		{"(DESCRIPTION=(SOURCE_ROUTE=\"yes\")(ADDRESS=(HOST=cman)(PORT=1521))(ADDRESS=(HOST=real)(PORT=1521))(CONNECT_DATA=(SERVICE_NAME=svc)))", ConnectTarget{}, false},

		// Conflicting or mixed identities.
		{"(DESCRIPTION_LIST=(DESCRIPTION=(ADDRESS=(HOST=h1)(PORT=1521))(CONNECT_DATA=(SERVICE_NAME=proddb)))(DESCRIPTION=(ADDRESS=(HOST=h2)(PORT=1521))(CONNECT_DATA=(SERVICE_NAME=reportdb))))", ConnectTarget{}, false},
		{"(DESCRIPTION=(ADDRESS=(HOST=a)(PORT=1521))(CONNECT_DATA=(SERVICE_NAME=orcl)(SID=orcl)))", ConnectTarget{}, false},
		{"(DESCRIPTION=(ADDRESS=(HOST=a)(PORT=1521))(CONNECT_DATA=(SERVICE_NAME=orcl)(SID=other)))", ConnectTarget{}, false},

		// Malformed connect strings.
		{"host1 host2:1521/svc", ConnectTarget{}, false},
		{"rac1:1521,/svc", ConnectTarget{}, false},
		{"rac1:/svc", ConnectTarget{}, false},
		{"rac1:1521:1522/svc", ConnectTarget{}, false},
		{"[::1:1521/svc", ConnectTarget{}, false},
		{"[::1]1521/svc", ConnectTarget{}, false},
		{"(DESCRIPTION=(ADDRESS=(PORT=1521))(CONNECT_DATA=(SERVICE_NAME=orcl)))", ConnectTarget{}, false},
		{"(DESCRIPTION=(ADDRESS=(HOST=a)(HOST=b)(PORT=1521))(CONNECT_DATA=(SERVICE_NAME=orcl)))", ConnectTarget{}, false},
		{"(DESCRIPTION=(ADDRESS=(HOST=a)(PORT=1521)(PORT=1522))(CONNECT_DATA=(SERVICE_NAME=orcl)))", ConnectTarget{}, false},
		{"(DESCRIPTION=(ADDRESS=(HOST= )(PORT=1521))(CONNECT_DATA=(SERVICE_NAME=orcl)))", ConnectTarget{}, false},
		{"(DESCRIPTION=(ADDRESS=(HOST=\"\")(PORT=1521))(CONNECT_DATA=(SERVICE_NAME=orcl)))", ConnectTarget{}, false},
		{"(DESCRIPTION=(ADDRESS=(HOST=\"db)(PORT=1521))(CONNECT_DATA=(SERVICE_NAME=orcl)))", ConnectTarget{}, false},
		{"(DESCRIPTION=(ADDRESS=(HOST=db\")(PORT=1521))(CONNECT_DATA=(SERVICE_NAME=orcl)))", ConnectTarget{}, false},
		{"(DESCRIPTION=(ADDRESS=(HOST=\"a\" \"b\")(PORT=1521))(CONNECT_DATA=(SERVICE_NAME=orcl)))", ConnectTarget{}, false},
		{"(DESCRIPTION=(ADDRESS=(HOST=a)(PORT=1521))(CONNECT_DATA=(SERVICE_NAME= )))", ConnectTarget{}, false},
		{"(DESCRIPTION=(ADDRESS=(HOST=a)(PORT=1521))(CONNECT_DATA=(SERVICE_NAME='')))", ConnectTarget{}, false},
		{"(DESCRIPTION=(ADDRESS=(HOST=a)(PORT=1521)(CONNECT_DATA=(SERVICE_NAME=orcl)))", ConnectTarget{}, false},
		{"(DESCRIPTION=(CONNECT_DATA=(SERVICE_NAME=orcl)))", ConnectTarget{}, false},
	}

	for _, tc := range tests {
		got, ok := ParseConnectString(tc.connectString)
		if ok != tc.ok {
			t.Errorf("ParseConnectString(%q) ok = %v, want %v (got %+v)", tc.connectString, ok, tc.ok, got)
			continue
		}
		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Errorf("ParseConnectString(%q) mismatch (-want +got):\n%s", tc.connectString, diff)
		}
	}
}
