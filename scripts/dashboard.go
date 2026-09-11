// Dashboard serves a web-based monitoring UI for all parser engine pipelines.
//
// Usage:
//
//	go run scripts/dashboard.go [-port 8088]
//
// Open http://localhost:8088 in your browser.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type Batch struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	File        string   `json:"file"`
	Status      string   `json:"status"`
	Description string   `json:"description"`
	Rules       []string `json:"rules"`
	DependsOn   []int    `json:"depends_on"`
	Tests       []string `json:"tests"`
	Error       string   `json:"error,omitempty"`
}

type Progress struct {
	Version    int     `json:"version"`
	Target     string  `json:"target"`
	Status     string  `json:"status"`
	AuditRound int     `json:"audit_round"`
	Batches    []Batch `json:"batches"`
}

type EngineSummary struct {
	Name        string  `json:"name"`
	Key         string  `json:"key"`
	Total       int     `json:"total"`
	Done        int     `json:"done"`
	InProgress  int     `json:"in_progress"`
	Pending     int     `json:"pending"`
	Failed      int     `json:"failed"`
	Percent     int     `json:"percent"`
	AuditRound  int     `json:"audit_round"`
	Batches     []Batch `json:"batches"`
	LastUpdated string  `json:"last_updated"`
	Phase       string  `json:"phase"` // "building", "auditing", "complete", "idle"
}

type DashboardData struct {
	Engines   []EngineSummary `json:"engines"`
	Timestamp string          `json:"timestamp"`
}

type AuditStatus struct {
	Running bool   `json:"running"`
	Engine  string `json:"engine"`
	StartedAt string `json:"started_at,omitempty"`
	Output  string `json:"output,omitempty"`
}

var (
	auditMu     sync.Mutex
	auditStatus = map[string]*AuditStatus{}
)

var engines = []struct {
	name string
	key  string
	path string
}{
	{"PostgreSQL", "pg", "pg/parser/PROGRESS.json"},
	{"MySQL", "mysql", "mysql/parser/PROGRESS.json"},
	{"MSSQL", "mssql", "mssql/parser/PROGRESS.json"},
	{"Oracle", "oracle", "oracle/parser/PROGRESS.json"},
}

func omniRoot() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Dir(filepath.Dir(file))
}

func loadEngine(root, name, key, relPath string) EngineSummary {
	path := filepath.Join(root, relPath)
	s := EngineSummary{Name: name, Key: key}

	data, err := os.ReadFile(path)
	if err != nil {
		s.Batches = []Batch{}
		s.Phase = "idle"
		return s
	}

	var p Progress
	if err := json.Unmarshal(data, &p); err != nil {
		s.Batches = []Batch{}
		s.Phase = "idle"
		return s
	}

	s.Total = len(p.Batches)
	s.Batches = p.Batches
	s.AuditRound = p.AuditRound
	for _, b := range p.Batches {
		switch b.Status {
		case "done":
			s.Done++
		case "in_progress":
			s.InProgress++
		case "failed":
			s.Failed++
		default:
			s.Pending++
		}
	}
	if s.Total > 0 {
		s.Percent = s.Done * 100 / s.Total
	}

	// Determine phase
	auditMu.Lock()
	as := auditStatus[key]
	auditMu.Unlock()
	if as != nil && as.Running {
		s.Phase = "auditing"
	} else if s.InProgress > 0 {
		s.Phase = "building"
	} else if s.Pending == 0 && s.Failed == 0 && s.Total > 0 {
		s.Phase = "complete"
	} else {
		s.Phase = "idle"
	}

	info, err := os.Stat(path)
	if err == nil {
		s.LastUpdated = info.ModTime().Format("15:04:05")
	}
	return s
}

func apiHandler(root string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var d DashboardData
		d.Timestamp = time.Now().Format("2006-01-02 15:04:05")
		for _, e := range engines {
			d.Engines = append(d.Engines, loadEngine(root, e.name, e.key, e.path))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(d)
	}
}

func auditHandler(root string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			auditMu.Lock()
			defer auditMu.Unlock()
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(auditStatus)
			return
		}

		// POST: trigger audit
		engine := r.URL.Query().Get("engine")
		if engine == "" {
			http.Error(w, "missing engine parameter", 400)
			return
		}

		// Validate engine
		valid := false
		for _, e := range engines {
			if e.key == engine {
				valid = true
				break
			}
		}
		if !valid {
			http.Error(w, "invalid engine: "+engine, 400)
			return
		}

		auditMu.Lock()
		if as, ok := auditStatus[engine]; ok && as.Running {
			auditMu.Unlock()
			http.Error(w, "audit already running for "+engine, 409)
			return
		}
		as := &AuditStatus{Running: true, Engine: engine, StartedAt: time.Now().Format("15:04:05")}
		auditStatus[engine] = as
		auditMu.Unlock()

		// Run audit in background
		go func() {
			defer func() {
				auditMu.Lock()
				as.Running = false
				auditMu.Unlock()
			}()

			auditSkill := filepath.Join(root, "scripts", "prompts", "audit-skill.md")
			skillData, err := os.ReadFile(auditSkill)
			if err != nil {
				auditMu.Lock()
				as.Output = "Error reading audit skill: " + err.Error()
				auditMu.Unlock()
				return
			}

			// Substitute engine name
			prompt := string(skillData)
			for {
				newPrompt := ""
				idx := 0
				for i := 0; i < len(prompt); i++ {
					if i+len("{{ENGINE}}") <= len(prompt) && prompt[i:i+len("{{ENGINE}}")] == "{{ENGINE}}" {
						newPrompt += prompt[idx:i] + engine
						idx = i + len("{{ENGINE}}")
						i = idx - 1
					}
				}
				newPrompt += prompt[idx:]
				if newPrompt == prompt {
					break
				}
				prompt = newPrompt
			}

			cmd := exec.Command("claude", "-p", prompt, "--dangerously-skip-permissions")
			cmd.Dir = root
			output, err := cmd.CombinedOutput()

			auditMu.Lock()
			if err != nil {
				as.Output = string(output) + "\nError: " + err.Error()
			} else {
				as.Output = string(output)
			}
			auditMu.Unlock()

			log.Printf("Audit complete for %s", engine)
		}()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "started", "engine": engine})
	}
}

const htmlPage = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<title>Omni Parser Dashboard</title>
<style>
  :root {
    --bg: #0d1117; --surface: #161b22; --border: #30363d;
    --text: #e6edf3; --muted: #8b949e; --accent: #58a6ff;
    --green: #3fb950; --yellow: #d29922; --red: #f85149; --blue: #58a6ff;
    --purple: #bc8cff;
  }
  * { box-sizing: border-box; margin: 0; padding: 0; }
  body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Helvetica, Arial, sans-serif;
         background: var(--bg); color: var(--text); padding: 24px; }
  h1 { font-size: 24px; font-weight: 600; margin-bottom: 4px; }
  .header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; }
  .timestamp { color: var(--muted); font-size: 13px; }
  .grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(480px, 1fr)); gap: 20px; }
  .card { background: var(--surface); border: 1px solid var(--border); border-radius: 8px; padding: 20px; }
  .card-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
  .card-title { font-size: 18px; font-weight: 600; }
  .card-stats { display: flex; gap: 8px; font-size: 13px; align-items: center; }
  .stat { padding: 2px 8px; border-radius: 12px; font-weight: 500; }
  .stat-done { background: rgba(63,185,80,0.15); color: var(--green); }
  .stat-progress { background: rgba(210,153,34,0.15); color: var(--yellow); }
  .stat-failed { background: rgba(248,81,73,0.15); color: var(--red); }
  .stat-pending { background: rgba(139,148,158,0.15); color: var(--muted); }
  .stat-audit { background: rgba(188,140,255,0.15); color: var(--purple); }
  .progress-bar { height: 8px; background: var(--border); border-radius: 4px; margin-bottom: 16px; overflow: hidden; display: flex; }
  .progress-fill-done { background: var(--green); transition: width 0.5s; }
  .progress-fill-progress { background: var(--yellow); transition: width 0.5s; }
  .progress-fill-failed { background: var(--red); transition: width 0.5s; }
  .percent { font-size: 13px; color: var(--muted); margin-bottom: 8px; text-align: right; }
  .batches { display: flex; flex-wrap: wrap; gap: 4px; }
  .batch { width: 28px; height: 28px; border-radius: 4px; display: flex; align-items: center;
           justify-content: center; font-size: 11px; font-weight: 600; cursor: default;
           transition: transform 0.15s, box-shadow 0.15s; position: relative; }
  .batch:hover { transform: scale(1.3); z-index: 10; box-shadow: 0 4px 12px rgba(0,0,0,0.4); }
  .batch-done { background: var(--green); color: #fff; }
  .batch-in_progress { background: var(--yellow); color: #000; animation: pulse 2s infinite; }
  .batch-pending { background: var(--border); color: var(--muted); }
  .batch-failed { background: var(--red); color: #fff; }
  .tooltip { display: none; position: absolute; bottom: 110%; left: 50%; transform: translateX(-50%);
             background: #1c2128; border: 1px solid var(--border); border-radius: 6px; padding: 8px 12px;
             white-space: nowrap; font-size: 12px; font-weight: 400; z-index: 100; pointer-events: none;
             box-shadow: 0 8px 24px rgba(0,0,0,0.4); max-width: 400px; white-space: normal; }
  .tooltip-name { font-weight: 600; margin-bottom: 2px; }
  .tooltip-desc { color: var(--muted); }
  .tooltip-error { color: var(--red); margin-top: 4px; }
  .batch:hover .tooltip { display: block; }
  @keyframes pulse { 0%,100% { opacity: 1; } 50% { opacity: 0.6; } }
  .overall { display: flex; gap: 24px; margin-bottom: 24px; }
  .overall-card { background: var(--surface); border: 1px solid var(--border); border-radius: 8px;
                  padding: 16px 24px; flex: 1; text-align: center; }
  .overall-num { font-size: 32px; font-weight: 700; }
  .overall-label { font-size: 13px; color: var(--muted); margin-top: 4px; }
  .updated { font-size: 12px; color: var(--muted); margin-top: 12px; display: flex; justify-content: space-between; align-items: center; }
  .phase-badge { font-size: 11px; padding: 2px 10px; border-radius: 12px; font-weight: 600; text-transform: uppercase; letter-spacing: 0.5px; }
  .phase-building { background: rgba(210,153,34,0.2); color: var(--yellow); }
  .phase-auditing { background: rgba(188,140,255,0.2); color: var(--purple); animation: pulse 2s infinite; }
  .phase-complete { background: rgba(63,185,80,0.2); color: var(--green); }
  .phase-idle { background: rgba(139,148,158,0.15); color: var(--muted); }
  .audit-btn { background: rgba(188,140,255,0.15); color: var(--purple); border: 1px solid rgba(188,140,255,0.3);
               border-radius: 6px; padding: 4px 12px; font-size: 12px; cursor: pointer; font-weight: 500; }
  .audit-btn:hover { background: rgba(188,140,255,0.25); }
  .audit-btn:disabled { opacity: 0.4; cursor: not-allowed; }
  .audit-row { display: flex; align-items: center; gap: 8px; }
</style>
</head>
<body>
<div class="header">
  <div>
    <h1>Omni Parser Dashboard</h1>
    <div class="timestamp" id="timestamp"></div>
  </div>
  <div class="timestamp">Auto-refresh: 3s</div>
</div>
<div class="overall" id="overall"></div>
<div class="grid" id="grid"></div>
<script>
let auditStates = {};

function render(data) {
  document.getElementById('timestamp').textContent = data.timestamp;

  let totalAll = 0, doneAll = 0, progAll = 0, failAll = 0;
  data.engines.forEach(e => {
    totalAll += e.total; doneAll += e.done; progAll += e.in_progress; failAll += e.failed;
  });
  const pendAll = totalAll - doneAll - progAll - failAll;
  const pctAll = totalAll ? Math.round(doneAll * 100 / totalAll) : 0;

  document.getElementById('overall').innerHTML =
    '<div class="overall-card"><div class="overall-num" style="color:var(--accent)">' + totalAll + '</div><div class="overall-label">Total Batches</div></div>' +
    '<div class="overall-card"><div class="overall-num" style="color:var(--green)">' + doneAll + '</div><div class="overall-label">Done</div></div>' +
    '<div class="overall-card"><div class="overall-num" style="color:var(--yellow)">' + progAll + '</div><div class="overall-label">In Progress</div></div>' +
    '<div class="overall-card"><div class="overall-num" style="color:var(--muted)">' + pendAll + '</div><div class="overall-label">Pending</div></div>' +
    '<div class="overall-card"><div class="overall-num" style="color:var(--red)">' + failAll + '</div><div class="overall-label">Failed</div></div>' +
    '<div class="overall-card"><div class="overall-num" style="color:var(--accent)">' + pctAll + '%</div><div class="overall-label">Overall</div></div>';

  const grid = document.getElementById('grid');
  grid.innerHTML = '';
  data.engines.forEach(e => {
    const progW = e.total ? (e.in_progress * 100 / e.total) : 0;
    const failW = e.total ? (e.failed * 100 / e.total) : 0;
    const doneW = e.total ? (e.done * 100 / e.total) : 0;

    let batchesHtml = '';
    (e.batches || []).forEach(b => {
      const err = b.error ? '<div class="tooltip-error">' + esc(b.error) + '</div>' : '';
      batchesHtml += '<div class="batch batch-' + b.status + '">' + b.id +
        '<div class="tooltip"><div class="tooltip-name">' + esc(b.name) + '</div>' +
        '<div class="tooltip-desc">' + esc(b.description || '') + '</div>' +
        (b.file ? '<div class="tooltip-desc">File: ' + esc(b.file) + '</div>' : '') +
        err + '</div></div>';
    });

    const phaseClass = 'phase-' + e.phase;
    const phaseLabel = e.phase === 'building' ? 'Building' :
                       e.phase === 'auditing' ? 'Auditing' :
                       e.phase === 'complete' ? 'Complete' : 'Idle';

    const isAuditing = e.phase === 'auditing';
    const auditBtnDisabled = isAuditing ? ' disabled' : '';
    const auditRoundText = e.audit_round > 0 ? 'Round ' + e.audit_round : 'Not run';

    grid.innerHTML += '<div class="card">' +
      '<div class="card-header"><span class="card-title">' + esc(e.name) + ' <span class="phase-badge ' + phaseClass + '">' + phaseLabel + '</span></span>' +
      '<div class="card-stats">' +
        '<span class="stat stat-done">' + e.done + ' done</span>' +
        (e.in_progress ? '<span class="stat stat-progress">' + e.in_progress + ' running</span>' : '') +
        (e.failed ? '<span class="stat stat-failed">' + e.failed + ' failed</span>' : '') +
        '<span class="stat stat-pending">' + e.pending + ' pending</span>' +
      '</div></div>' +
      '<div class="percent">' + e.percent + '% complete (' + e.done + '/' + e.total + ')</div>' +
      '<div class="progress-bar">' +
        '<div class="progress-fill-done" style="width:' + doneW + '%"></div>' +
        '<div class="progress-fill-progress" style="width:' + progW + '%"></div>' +
        '<div class="progress-fill-failed" style="width:' + failW + '%"></div>' +
      '</div>' +
      '<div class="batches">' + batchesHtml + '</div>' +
      '<div class="updated">' +
        '<span>Last updated: ' + (e.last_updated || 'N/A') + '</span>' +
        '<div class="audit-row">' +
          '<span class="stat stat-audit">Audit: ' + auditRoundText + '</span>' +
          '<button class="audit-btn" onclick="triggerAudit(\'' + e.key + '\')"' + auditBtnDisabled + '>Run Audit</button>' +
        '</div>' +
      '</div>' +
    '</div>';
  });
}

async function triggerAudit(engine) {
  if (!confirm('Run audit for ' + engine + '? This will call Claude to analyze gaps and may take a few minutes.')) return;
  try {
    const r = await fetch('/api/audit?engine=' + engine, {method: 'POST'});
    const d = await r.json();
    if (r.ok) {
      alert('Audit started for ' + engine);
    } else {
      alert('Error: ' + (d.error || r.statusText));
    }
  } catch(e) {
    alert('Failed to trigger audit: ' + e.message);
  }
}

function esc(s) {
  const d = document.createElement('div');
  d.textContent = s;
  return d.innerHTML;
}

async function refresh() {
  try {
    const r = await fetch('/api/status');
    const d = await r.json();
    render(d);
  } catch(e) {
    console.error('fetch error:', e);
  }
}

refresh();
setInterval(refresh, 3000);
</script>
</body>
</html>`

func main() {
	port := flag.Int("port", 8088, "HTTP port")
	flag.Parse()

	root := omniRoot()
	log.Printf("Omni root: %s", root)
	log.Printf("Dashboard: http://localhost:%d", *port)

	http.HandleFunc("/api/status", apiHandler(root))
	http.HandleFunc("/api/audit", auditHandler(root))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, htmlPage)
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
