package build

import "testing"

func TestPythonRequirement(t *testing.T) {
	fakeFiles := []string{`
psycopg2==2.9.3
pysqlite==2.8.3
`, `
psycopg2
pysqlite
`, `
psycopg2 < 2
pysqlite >= 123
`, `
psycopg2-binary==2.9.3
pysqlite3-binary==0.4.6
`, `
psycopg2==2.9.3
psycopg2-binary==2.9.3 
pysqlite==2.8.3
pysqlite3-binary==0.4.6
`,
	}
	expected := `
psycopg2-binary==2.9.3
pysqlite3-binary==0.4.6
`
	for _, file := range fakeFiles {
		result := ReplaceLambdaRequirements(file)
		if result != expected {
			t.Errorf("Expected %s, got %s", expected, result)
		}
	}
}
