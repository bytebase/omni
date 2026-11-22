package parser

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/compat"
)

// TestCompatibility holds omni to testdata/compat.tsv: a statement for each
// alternative of ../testdata/SqlBaseParser.g4 (Apache Spark's official ANTLR 4
// SQL grammar, which Databricks SQL is built on), its near-misses, and hand-
// curated suite rows, with omni's verdict and Apache Spark 4.0.5's verdict
// from spark.sessionState.sqlParser.parsePlan. The test replays the committed
// file with no container or network. OMNI_RECORD=1 records unjudged rows and
// OMNI_RECORD=verify re-checks every committed verdict against the pinned
// apache/spark:4.0.5 container.
func TestCompatibility(t *testing.T) {
	grammar := filepath.Join("..", "testdata", "SqlBaseParser.g4")
	lexer := filepath.Join("..", "testdata", "SqlBaseLexer.g4")
	compat.Check(t, compat.Suite{
		Path:    filepath.Join("testdata", "compat.tsv"),
		Pins:    filepath.Join("..", "testdata", "pins.tsv"),
		Grammar: grammar,
		Lexer:   lexer,
		Tokens:  filepath.Join("..", "testdata", "tokens.tsv"),
		Omni:    omniVerdict,
		Source:  compat.SourceEngine,
		Judge:   sparkJudge,
	})
}

// omniVerdict accepts a row that parses cleanly as one Databricks SQL statement.
func omniVerdict(sql string) string {
	res, err := Parse(sql)
	if err != nil || res == nil || len(res.Statements) != 1 {
		return compat.Reject
	}
	return compat.Accept
}

const sparkJudgeJavaSource = `
import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.io.PrintStream;
import java.nio.charset.StandardCharsets;
import java.util.Base64;
import org.apache.spark.SparkThrowable;
import org.apache.spark.sql.SparkSession;
import org.apache.spark.sql.catalyst.parser.ParserInterface;

public class SparkJudge {
    public static void main(String[] args) throws Exception {
        PrintStream realOut = System.out;
        System.setOut(System.err);
        SparkSession spark = SparkSession.builder()
            .master("local[1]")
            .config("spark.driver.bindAddress", "127.0.0.1")
            .config("spark.ui.enabled", "false")
            .config("spark.sql.scripting.enabled", "true")
            .getOrCreate();
        ParserInterface parser = spark.sessionState().sqlParser();
        realOut.println("READY");
        realOut.flush();
        BufferedReader in = new BufferedReader(new InputStreamReader(System.in, StandardCharsets.UTF_8));
        String line;
        while ((line = in.readLine()) != null) {
            if (line.isEmpty()) {
                continue;
            }
            String sql = new String(Base64.getDecoder().decode(line), StandardCharsets.UTF_8);
            try {
                parser.parsePlan(sql);
                realOut.println("accept");
            } catch (Throwable e) {
                String cond = e.getClass().getSimpleName();
                if (e instanceof SparkThrowable) {
                    String c = ((SparkThrowable) e).getCondition();
                    if (c != null && !c.isEmpty()) {
                        cond = c;
                    }
                }
                realOut.println("reject\t" + cond);
            }
            realOut.flush();
        }
        spark.stop();
    }
}
`

// sparkJudge runs a single persistent apache/spark container pinned by digest
// and asks spark.sessionState.sqlParser.parsePlan (with SQL scripting enabled)
// about each statement over a base64 line protocol on stdin/stdout.
func sparkJudge(t testing.TB, pin enginetest.Pin) func(string) (string, error) {
	t.Helper()
	if !enginetest.Recording(t) {
		t.Fatalf("sparkJudge called outside OMNI_RECORD=1 or OMNI_RECORD=verify")
	}
	script := "cat > /tmp/SparkJudge.java << 'SPARK_JUDGE_EOF'\n" +
		sparkJudgeJavaSource +
		"\nSPARK_JUDGE_EOF\n" +
		"exec java -cp '/opt/spark/jars/*' /tmp/SparkJudge.java"
	cmd := exec.Command("docker", "run", "--rm", "-i", pin.ImageRef(), "bash", "-c", script)
	in, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("sparkJudge stdin pipe: %v", err)
	}
	out, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("sparkJudge stdout pipe: %v", err)
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		t.Fatalf("start %s: %v", pin.ImageRef(), err)
	}
	reader := bufio.NewReader(out)
	t.Cleanup(func() {
		_ = in.Close()
		_ = cmd.Wait()
	})
	ready, err := reader.ReadString('\n')
	if err != nil || strings.TrimSpace(ready) != "READY" {
		_ = in.Close()
		_ = cmd.Wait()
		t.Fatalf("sparkJudge %s did not become ready (got %q, err=%v)\nstderr:\n%s",
			pin.ImageRef(), strings.TrimSpace(ready), err, stderr.String())
	}

	var mu sync.Mutex
	return func(sql string) (string, error) {
		mu.Lock()
		defer mu.Unlock()
		encoded := base64.StdEncoding.EncodeToString([]byte(sql))
		if _, err := io.WriteString(in, encoded+"\n"); err != nil {
			return "", fmt.Errorf("write to sparkJudge: %w", err)
		}
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("read from sparkJudge: %w (stderr tail: %s)", err, stderrTail(stderr.String()))
		}
		line = strings.TrimRight(line, "\r\n")
		switch {
		case line == compat.Accept:
			return compat.Accept, nil
		case line == compat.Reject || strings.HasPrefix(line, compat.Reject+"\t"):
			return compat.Reject, nil
		default:
			return "", fmt.Errorf("unexpected sparkJudge response %q", line)
		}
	}
}

func stderrTail(s string) string {
	const maxLen = 1000
	if len(s) <= maxLen {
		return s
	}
	return s[len(s)-maxLen:]
}
