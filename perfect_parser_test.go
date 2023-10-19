package parser_test

import (
	"testing"

	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
)

func TestPerfectParse(t *testing.T) {
	parser := parser.New()

	stmt, _, err := parser.PerfectParse("OPTIMIZE TABLE foo;", "", "")
	if err != nil {
		t.Error(err)
		return
	}
	if _, ok := stmt[0].(*ast.UnparsedStmt); !ok {
		t.Errorf("expect stmt type is unparsedStmt, actual is %T", stmt)
		return
	}

	type testCase struct {
		sql    string
		expect []string
	}

	tc := []testCase{
		{
			sql: `SELECT * FROM db1.t1`,
			expect: []string{
				`SELECT * FROM db1.t1`,
			},
		},
		{
			sql: `SELECT * FROM db1.t1;SELECT * FROM db2.t2`,
			expect: []string{
				"SELECT * FROM db1.t1;",
				"SELECT * FROM db2.t2",
			},
		},
		{
			sql: "SELECT * FROM db1.t1;OPTIMIZE TABLE foo;SELECT * FROM db2.t2",
			expect: []string{
				"SELECT * FROM db1.t1;",
				"OPTIMIZE TABLE foo;",
				"SELECT * FROM db2.t2",
			},
		},
		{
			sql: "OPTIMIZE TABLE foo;SELECT * FROM db1.t1;SELECT * FROM db2.t2",
			expect: []string{
				"OPTIMIZE TABLE foo;",
				"SELECT * FROM db1.t1;",
				"SELECT * FROM db2.t2",
			},
		},
		{
			sql: "SELECT * FROM db1.t1;SELECT * FROM db2.t2;OPTIMIZE TABLE foo",
			expect: []string{
				"SELECT * FROM db1.t1;",
				"SELECT * FROM db2.t2;",
				"OPTIMIZE TABLE foo",
			},
		},
		{
			sql: "SELECT FROM db2.t2 where a=\"asd;\"; SELECT * FROM db1.t1;",
			expect: []string{
				"SELECT FROM db2.t2 where a=\"asd;\";",
				" SELECT * FROM db1.t1;",
			},
		},
		{
			sql: "SELECT * FROM db1.t1;OPTIMIZE TABLE foo;OPTIMIZE TABLE foo;SELECT * FROM db2.t2",
			expect: []string{
				"SELECT * FROM db1.t1;",
				"OPTIMIZE TABLE foo;",
				"OPTIMIZE TABLE foo;",
				"SELECT * FROM db2.t2",
			},
		},
		{
			sql: "OPTIMIZE TABLE foo;SELECT * FROM db1.t1;OPTIMIZE TABLE foo;SELECT * FROM db2.t2",
			expect: []string{
				"OPTIMIZE TABLE foo;",
				"SELECT * FROM db1.t1;",
				"OPTIMIZE TABLE foo;",
				"SELECT * FROM db2.t2",
			},
		},
		{
			sql: "SELECT * FROM db1.t1;OPTIMIZE TABLE foo;SELECT * FROM db2.t2;OPTIMIZE TABLE foo",
			expect: []string{
				"SELECT * FROM db1.t1;",
				"OPTIMIZE TABLE foo;",
				"SELECT * FROM db2.t2;",
				"OPTIMIZE TABLE foo",
			},
		},
		{
			sql: `
CREATE PROCEDURE proc1(OUT s int)
BEGIN
END;
`,
			expect: []string{
				`
CREATE PROCEDURE proc1(OUT s int)
BEGIN
END;`,
			},
		},
		{
			sql: `
CREATE PROCEDURE proc1(OUT s int)
BEGIN
SELECT COUNT(*)  FROM user;
END;
`,
			expect: []string{
				`
CREATE PROCEDURE proc1(OUT s int)
BEGIN
SELECT COUNT(*)  FROM user;
END;`,
			},
		},
		{
			sql: `
CREATE PROCEDURE proc1(OUT s int)
BEGIN
SELECT COUNT(*)  FROM user;
SELECT COUNT(*)  FROM user;
END;
`,
			expect: []string{
				`
CREATE PROCEDURE proc1(OUT s int)
BEGIN
SELECT COUNT(*)  FROM user;
SELECT COUNT(*)  FROM user;
END;`,
			},
		},
		{
			sql: `
SELECT * FROM db1.t1;
CREATE PROCEDURE proc1(OUT s int)
BEGIN
END;
`,
			expect: []string{
				`SELECT * FROM db1.t1;`,
				`
CREATE PROCEDURE proc1(OUT s int)
BEGIN
END;`,
			},
		},
		{
			sql: `
SELECT * FROM db1.t1;
CREATE PROCEDURE proc1(OUT s int)
BEGIN
 SELECT COUNT(*)  FROM user;
 SELECT COUNT(*)  FROM user;
END;
`,
			expect: []string{
				`SELECT * FROM db1.t1;`,
				`
CREATE PROCEDURE proc1(OUT s int)
BEGIN
 SELECT COUNT(*)  FROM user;
 SELECT COUNT(*)  FROM user;
END;`,
			},
		},
		{
			sql: `
CREATE PROCEDURE proc1(OUT s int)
BEGIN
 SELECT COUNT(*)  FROM user;
 SELECT COUNT(*)  FROM user;
END;
SELECT * FROM db1.t1;
`,
			expect: []string{
				`
CREATE PROCEDURE proc1(OUT s int)
BEGIN
 SELECT COUNT(*)  FROM user;
 SELECT COUNT(*)  FROM user;
END;`,
				`SELECT * FROM db1.t1;`,
			},
		},
		{
			sql: `
SELECT * FROM db1.t1;
CREATE PROCEDURE proc1(OUT s int)
BEGIN
 SELECT COUNT(*)  FROM user;
 SELECT COUNT(*)  FROM user;
END;
CREATE PROCEDURE proc1(OUT s int)
BEGIN
 SELECT COUNT(*)  FROM user;
 SELECT COUNT(*)  FROM user;
 SELECT COUNT(*)  FROM user;
 SELECT COUNT(*)  FROM user;
END;
SELECT * FROM db1.t1;
`,
			expect: []string{
				`SELECT * FROM db1.t1;`,
				`
CREATE PROCEDURE proc1(OUT s int)
BEGIN
 SELECT COUNT(*)  FROM user;
 SELECT COUNT(*)  FROM user;
END;`,
				`
CREATE PROCEDURE proc1(OUT s int)
BEGIN
 SELECT COUNT(*)  FROM user;
 SELECT COUNT(*)  FROM user;
 SELECT COUNT(*)  FROM user;
 SELECT COUNT(*)  FROM user;
END;`,
				`SELECT * FROM db1.t1;`,
			},
		},
		{
			sql: `
SELECT * FROM db1.t1;
CREATE PROCEDURE proc1(OUT s int)
BEGIN
 SELECT COUNT(*)  FROM user;
 SELECT COUNT(*)  FROM user;
END;
SELECT * FROM db1.t1;
CREATE PROCEDURE proc1(OUT s int)
BEGIN
 SELECT COUNT(*)  FROM user;
 SELECT COUNT(*)  FROM user;
 SELECT COUNT(*)  FROM user;
 SELECT COUNT(*)  FROM user;
END;
SELECT * FROM db1.t1;
`,
			expect: []string{
				`SELECT * FROM db1.t1;`,
				`
CREATE PROCEDURE proc1(OUT s int)
BEGIN
 SELECT COUNT(*)  FROM user;
 SELECT COUNT(*)  FROM user;
END;`,
				`SELECT * FROM db1.t1;`,
				`
CREATE PROCEDURE proc1(OUT s int)
BEGIN
 SELECT COUNT(*)  FROM user;
 SELECT COUNT(*)  FROM user;
 SELECT COUNT(*)  FROM user;
 SELECT COUNT(*)  FROM user;
END;`,
				`SELECT * FROM db1.t1;`,
			},
		},
		{ // 匹配特殊字符结束
			sql: "select * from  �E",
			expect: []string{
				`select * from  �E`,
			},
		},
		{ // 匹配特殊字符后是;
			sql: "select * from  �E;select * from t1",
			expect: []string{
				`select * from  �E;`,
				"select * from t1",
			},
		},
		{ // 匹配特殊字符在中间
			sql: "select * from  �E where id = 1;select * from  �E ",
			expect: []string{
				`select * from  �E where id = 1;`,
				`select * from  �E `,
			},
		},
		{ // 匹配特殊字符在开头
			sql: " where id = 1;select * from  �E ",
			expect: []string{
				` where id = 1;`,
				`select * from  �E `,
			},
		},
		{ // 匹配特殊字符在SQL开头
			sql: "select * from  �E ; where id = 1",
			expect: []string{
				`select * from  �E ;`,
				` where id = 1`,
			},
		},
		{ // 匹配其他invalid场景
			sql: "@`",
			expect: []string{
				"@`",
			},
		},
		{ // 匹配其他invalid场景
			sql: "@` ;select * from t1",
			expect: []string{
				"@` ;select * from t1",
			},
		},
	}
	for _, c := range tc {
		stmt, _, err := parser.PerfectParse(c.sql, "", "")
		if err != nil {
			t.Error(err)
			return
		}
		if len(c.expect) != len(stmt) {
			t.Errorf("expect sql length is %d, actual is %d, sql is [%s]", len(c.expect), len(stmt), c.sql)
		} else {
			for i, s := range stmt {
				if s.Text() != c.expect[i] {
					t.Errorf("expect sql is [%s], actual is [%s]", c.expect[i], s.Text())
				}
			}
		}
	}
}
