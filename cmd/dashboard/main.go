package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/barrynorthern/libretto/internal/db"
	_ "github.com/mattn/go-sqlite3"
)

type Dashboard struct {
	queries *db.Queries
}

type ProjectSummary struct {
	Project   db.Project
	Versions  []db.GraphVersion
	Stats     ProjectStats
}

type ProjectStats struct {
	TotalEntities     int64
	TotalRelationships int
	TotalAnnotations  int
	EntityCounts      map[string]int64
	RelationshipCounts map[string]int
}

type GraphVisualization struct {
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
}

type Node struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Group    int    `json:"group"`
	Size     int    `json:"size"`
}

type Link struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Type   string `json:"type"`
	Value  int    `json:"value"`
}

func main() {
	var (
		dbPath = flag.String("db", "libretto.db", "Path to SQLite database")
		port   = flag.String("port", "8080", "Port to serve on")
	)
	flag.Parse()

	database, err := sql.Open("sqlite3", *dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer database.Close()

	dashboard := &Dashboard{
		queries: db.New(database),
	}

	http.HandleFunc("/", dashboard.handleHome)
	http.HandleFunc("/project/", dashboard.handleProject)
	http.HandleFunc("/graph/", dashboard.handleGraph)
	http.HandleFunc("/api/graph/", dashboard.handleGraphAPI)
	http.HandleFunc("/static/", dashboard.handleStatic)

	fmt.Printf("Dashboard server starting on http://localhost:%s\n", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func (d *Dashboard) handleHome(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	
	projects, err := d.queries.ListProjects(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list projects: %v", err), http.StatusInternalServerError)
		return
	}

	var projectSummaries []ProjectSummary
	for _, project := range projects {
		versions, err := d.queries.ListGraphVersionsByProject(ctx, project.ID)
		if err != nil {
			log.Printf("Failed to get versions for project %s: %v", project.ID, err)
			continue
		}

		var stats ProjectStats
		if len(versions) > 0 {
			// Get stats for working set version
			for _, version := range versions {
				if version.IsWorkingSet {
					stats = d.getProjectStats(ctx, version.ID)
					break
				}
			}
		}

		projectSummaries = append(projectSummaries, ProjectSummary{
			Project:  project,
			Versions: versions,
			Stats:    stats,
		})
	}

	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Libretto Dashboard</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; }
        .header { background: #2c3e50; color: white; padding: 20px; border-radius: 8px; margin-bottom: 30px; }
        .project-card { background: white; border-radius: 8px; padding: 20px; margin-bottom: 20px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .project-title { color: #2c3e50; margin-bottom: 10px; }
        .project-meta { color: #7f8c8d; margin-bottom: 15px; }
        .stats { display: flex; gap: 20px; margin-top: 15px; }
        .stat { background: #ecf0f1; padding: 10px; border-radius: 4px; text-align: center; min-width: 80px; }
        .stat-value { font-size: 24px; font-weight: bold; color: #2c3e50; }
        .stat-label { font-size: 12px; color: #7f8c8d; }
        .actions { margin-top: 15px; }
        .btn { background: #3498db; color: white; padding: 8px 16px; text-decoration: none; border-radius: 4px; margin-right: 10px; }
        .btn:hover { background: #2980b9; }
        .no-projects { text-align: center; color: #7f8c8d; padding: 40px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Libretto Narrative Engine Dashboard</h1>
            <p>Monitor and visualize your narrative graphs</p>
        </div>

        {{if .}}
            {{range .}}
            <div class="project-card">
                <h2 class="project-title">{{.Project.Name}}</h2>
                <div class="project-meta">
                    <strong>Theme:</strong> {{if .Project.Theme.Valid}}{{.Project.Theme.String}}{{else}}Not set{{end}} | 
                    <strong>Genre:</strong> {{if .Project.Genre.Valid}}{{.Project.Genre.String}}{{else}}Not set{{end}} | 
                    <strong>Versions:</strong> {{len .Versions}}
                </div>
                {{if .Project.Description.Valid}}
                <p>{{.Project.Description.String}}</p>
                {{end}}
                
                <div class="stats">
                    <div class="stat">
                        <div class="stat-value">{{.Stats.TotalEntities}}</div>
                        <div class="stat-label">Entities</div>
                    </div>
                    <div class="stat">
                        <div class="stat-value">{{.Stats.TotalRelationships}}</div>
                        <div class="stat-label">Relationships</div>
                    </div>
                    <div class="stat">
                        <div class="stat-value">{{.Stats.TotalAnnotations}}</div>
                        <div class="stat-label">Annotations</div>
                    </div>
                </div>

                <div class="actions">
                    <a href="/project/{{.Project.ID}}" class="btn">View Details</a>
                    <a href="/graph/{{.Project.ID}}" class="btn">Visualize Graph</a>
                </div>
            </div>
            {{end}}
        {{else}}
            <div class="no-projects">
                <h3>No projects found</h3>
                <p>Use the dbseed tool to create sample data: <code>go run cmd/dbseed/main.go</code></p>
            </div>
        {{end}}
    </div>
</body>
</html>
`

	t, err := template.New("home").Parse(tmpl)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, projectSummaries); err != nil {
		http.Error(w, fmt.Sprintf("Template execution error: %v", err), http.StatusInternalServerError)
		return
	}
}

func (d *Dashboard) handleProject(w http.ResponseWriter, r *http.Request) {
	projectID := r.URL.Path[len("/project/"):]
	if projectID == "" {
		http.Error(w, "Project ID required", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	
	project, err := d.queries.GetProject(ctx, projectID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get project: %v", err), http.StatusInternalServerError)
		return
	}

	versions, err := d.queries.ListGraphVersionsByProject(ctx, projectID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get versions: %v", err), http.StatusInternalServerError)
		return
	}

	// Get working set version details
	var workingSetVersion *db.GraphVersion
	for _, version := range versions {
		if version.IsWorkingSet {
			workingSetVersion = &version
			break
		}
	}

	var entities []db.Entity
	var relationships []db.Relationship
	var entityCounts map[string]int64
	
	if workingSetVersion != nil {
		entities, err = d.queries.ListEntitiesByVersion(ctx, workingSetVersion.ID)
		if err != nil {
			log.Printf("Failed to get entities: %v", err)
		}

		relationships, err = d.queries.ListRelationshipsByVersion(ctx, workingSetVersion.ID)
		if err != nil {
			log.Printf("Failed to get relationships: %v", err)
		}

		entityCounts = make(map[string]int64)
		entityTypes := []string{"Scene", "Character", "Location", "Theme", "PlotPoint", "Arc"}
		for _, entityType := range entityTypes {
			params := db.CountEntitiesByTypeParams{
				VersionID:  workingSetVersion.ID,
				EntityType: entityType,
			}
			count, err := d.queries.CountEntitiesByType(ctx, params)
			if err == nil {
				entityCounts[entityType] = count
			}
		}
	}

	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>{{.Project.Name}} - Libretto Dashboard</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; }
        .header { background: #2c3e50; color: white; padding: 20px; border-radius: 8px; margin-bottom: 30px; }
        .section { background: white; border-radius: 8px; padding: 20px; margin-bottom: 20px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .entity-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 20px; }
        .entity-card { border: 1px solid #ddd; border-radius: 4px; padding: 15px; }
        .entity-type { background: #3498db; color: white; padding: 4px 8px; border-radius: 4px; font-size: 12px; margin-bottom: 10px; display: inline-block; }
        .relationship-list { list-style: none; padding: 0; }
        .relationship-list li { padding: 8px; border-bottom: 1px solid #eee; }
        .btn { background: #3498db; color: white; padding: 8px 16px; text-decoration: none; border-radius: 4px; margin-right: 10px; }
        .btn:hover { background: #2980b9; }
        .stats { display: flex; gap: 20px; margin-bottom: 20px; }
        .stat { background: #ecf0f1; padding: 15px; border-radius: 4px; text-align: center; flex: 1; }
        .stat-value { font-size: 24px; font-weight: bold; color: #2c3e50; }
        .stat-label { font-size: 12px; color: #7f8c8d; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>{{.Project.Name}}</h1>
            <p>{{if .Project.Description.Valid}}{{.Project.Description.String}}{{end}}</p>
            <a href="/" class="btn">← Back to Dashboard</a>
            <a href="/graph/{{.Project.ID}}" class="btn">Visualize Graph</a>
        </div>

        {{if .WorkingSetVersion}}
        <div class="section">
            <h2>Statistics</h2>
            <div class="stats">
                {{range $type, $count := .EntityCounts}}
                <div class="stat">
                    <div class="stat-value">{{$count}}</div>
                    <div class="stat-label">{{$type}}</div>
                </div>
                {{end}}
                <div class="stat">
                    <div class="stat-value">{{len .Relationships}}</div>
                    <div class="stat-label">Relationships</div>
                </div>
            </div>
        </div>

        <div class="section">
            <h2>Entities ({{len .Entities}})</h2>
            <div class="entity-grid">
                {{range .Entities}}
                <div class="entity-card">
                    <div class="entity-type">{{.EntityType}}</div>
                    <h3>{{.Name}}</h3>
                    <p><strong>ID:</strong> {{.ID}}</p>
                    <p><strong>Created:</strong> {{.CreatedAt.Format "2006-01-02 15:04"}}</p>
                </div>
                {{end}}
            </div>
        </div>

        <div class="section">
            <h2>Relationships ({{len .Relationships}})</h2>
            <ul class="relationship-list">
                {{range .Relationships}}
                <li>
                    <strong>{{.RelationshipType}}</strong>: 
                    {{.FromEntityID}} → {{.ToEntityID}}
                    <small>({{.CreatedAt.Format "2006-01-02 15:04"}})</small>
                </li>
                {{end}}
            </ul>
        </div>
        {{else}}
        <div class="section">
            <h2>No Working Set Version</h2>
            <p>This project doesn't have a working set version yet.</p>
        </div>
        {{end}}

        <div class="section">
            <h2>Versions ({{len .Versions}})</h2>
            {{range .Versions}}
            <div style="padding: 10px; border: 1px solid #ddd; margin-bottom: 10px; border-radius: 4px;">
                <h4>{{if .Name.Valid}}{{.Name.String}}{{else}}Unnamed Version{{end}} 
                    {{if .IsWorkingSet}}<span style="background: #27ae60; color: white; padding: 2px 6px; border-radius: 3px; font-size: 10px;">WORKING SET</span>{{end}}
                </h4>
                <p>{{if .Description.Valid}}{{.Description.String}}{{end}}</p>
                <small>Created: {{.CreatedAt.Format "2006-01-02 15:04"}}</small>
            </div>
            {{end}}
        </div>
    </div>
</body>
</html>
`

	data := struct {
		Project           db.Project
		Versions          []db.GraphVersion
		WorkingSetVersion *db.GraphVersion
		Entities          []db.Entity
		Relationships     []db.Relationship
		EntityCounts      map[string]int64
	}{
		Project:           project,
		Versions:          versions,
		WorkingSetVersion: workingSetVersion,
		Entities:          entities,
		Relationships:     relationships,
		EntityCounts:      entityCounts,
	}

	t, err := template.New("project").Parse(tmpl)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, data); err != nil {
		http.Error(w, fmt.Sprintf("Template execution error: %v", err), http.StatusInternalServerError)
		return
	}
}

func (d *Dashboard) handleGraph(w http.ResponseWriter, r *http.Request) {
	projectID := r.URL.Path[len("/graph/"):]
	if projectID == "" {
		http.Error(w, "Project ID required", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	
	project, err := d.queries.GetProject(ctx, projectID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get project: %v", err), http.StatusInternalServerError)
		return
	}

	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>{{.Name}} - Graph Visualization</title>
    <script src="https://d3js.org/d3.v7.min.js"></script>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; background: #f5f5f5; }
        .header { background: #2c3e50; color: white; padding: 20px; }
        .container { display: flex; height: calc(100vh - 80px); }
        .sidebar { width: 300px; background: white; padding: 20px; overflow-y: auto; }
        .graph-container { flex: 1; background: white; margin: 20px; border-radius: 8px; }
        #graph { width: 100%; height: 100%; }
        .node { cursor: pointer; }
        .link { stroke: #999; stroke-opacity: 0.6; }
        .node text { font: 12px sans-serif; pointer-events: none; }
        .legend { margin-bottom: 20px; }
        .legend-item { display: flex; align-items: center; margin-bottom: 5px; }
        .legend-color { width: 20px; height: 20px; margin-right: 10px; border-radius: 3px; }
        .btn { background: #3498db; color: white; padding: 8px 16px; text-decoration: none; border-radius: 4px; margin-right: 10px; }
        .btn:hover { background: #2980b9; }
    </style>
</head>
<body>
    <div class="header">
        <h1>{{.Name}} - Graph Visualization</h1>
        <a href="/project/{{.ID}}" class="btn">← Back to Project</a>
    </div>
    
    <div class="container">
        <div class="sidebar">
            <div class="legend">
                <h3>Entity Types</h3>
                <div class="legend-item">
                    <div class="legend-color" style="background: #e74c3c;"></div>
                    <span>Scene</span>
                </div>
                <div class="legend-item">
                    <div class="legend-color" style="background: #3498db;"></div>
                    <span>Character</span>
                </div>
                <div class="legend-item">
                    <div class="legend-color" style="background: #2ecc71;"></div>
                    <span>Location</span>
                </div>
                <div class="legend-item">
                    <div class="legend-color" style="background: #f39c12;"></div>
                    <span>Theme</span>
                </div>
                <div class="legend-item">
                    <div class="legend-color" style="background: #9b59b6;"></div>
                    <span>PlotPoint</span>
                </div>
                <div class="legend-item">
                    <div class="legend-color" style="background: #1abc9c;"></div>
                    <span>Arc</span>
                </div>
            </div>
            
            <div id="node-info">
                <h3>Node Information</h3>
                <p>Click on a node to see details</p>
            </div>
        </div>
        
        <div class="graph-container">
            <svg id="graph"></svg>
        </div>
    </div>

    <script>
        const projectId = "{{.ID}}";
        
        // Color mapping for entity types
        const colors = {
            'Scene': '#e74c3c',
            'Character': '#3498db',
            'Location': '#2ecc71',
            'Theme': '#f39c12',
            'PlotPoint': '#9b59b6',
            'Arc': '#1abc9c'
        };

        // Set up SVG
        const svg = d3.select("#graph");
        const container = d3.select(".graph-container");
        const width = container.node().getBoundingClientRect().width;
        const height = container.node().getBoundingClientRect().height;
        
        svg.attr("width", width).attr("height", height);

        // Load and visualize graph data
        fetch('/api/graph/' + projectId)
            .then(response => response.json())
            .then(data => {
                createGraph(data);
            })
            .catch(error => {
                console.error('Error loading graph data:', error);
            });

        function createGraph(data) {
            // Create force simulation
            const simulation = d3.forceSimulation(data.nodes)
                .force("link", d3.forceLink(data.links).id(d => d.id).distance(100))
                .force("charge", d3.forceManyBody().strength(-300))
                .force("center", d3.forceCenter(width / 2, height / 2));

            // Create links
            const link = svg.append("g")
                .selectAll("line")
                .data(data.links)
                .enter().append("line")
                .attr("class", "link")
                .attr("stroke-width", d => Math.sqrt(d.value));

            // Create nodes
            const node = svg.append("g")
                .selectAll("circle")
                .data(data.nodes)
                .enter().append("circle")
                .attr("class", "node")
                .attr("r", d => 5 + d.size)
                .attr("fill", d => colors[d.type] || '#95a5a6')
                .call(d3.drag()
                    .on("start", dragstarted)
                    .on("drag", dragged)
                    .on("end", dragended))
                .on("click", function(event, d) {
                    showNodeInfo(d);
                });

            // Add labels
            const label = svg.append("g")
                .selectAll("text")
                .data(data.nodes)
                .enter().append("text")
                .text(d => d.name)
                .attr("font-size", "10px")
                .attr("dx", 12)
                .attr("dy", 4);

            // Update positions on simulation tick
            simulation.on("tick", () => {
                link
                    .attr("x1", d => d.source.x)
                    .attr("y1", d => d.source.y)
                    .attr("x2", d => d.target.x)
                    .attr("y2", d => d.target.y);

                node
                    .attr("cx", d => d.x)
                    .attr("cy", d => d.y);

                label
                    .attr("x", d => d.x)
                    .attr("y", d => d.y);
            });

            // Drag functions
            function dragstarted(event, d) {
                if (!event.active) simulation.alphaTarget(0.3).restart();
                d.fx = d.x;
                d.fy = d.y;
            }

            function dragged(event, d) {
                d.fx = event.x;
                d.fy = event.y;
            }

            function dragended(event, d) {
                if (!event.active) simulation.alphaTarget(0);
                d.fx = null;
                d.fy = null;
            }
        }

        function showNodeInfo(node) {
            const infoDiv = document.getElementById('node-info');
            infoDiv.innerHTML = ` + "`" + `
                <h3>${node.name}</h3>
                <p><strong>Type:</strong> ${node.type}</p>
                <p><strong>ID:</strong> ${node.id}</p>
                <p><strong>Connections:</strong> ${node.size}</p>
            ` + "`" + `;
        }
    </script>
</body>
</html>
`

	t, err := template.New("graph").Parse(tmpl)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, project); err != nil {
		http.Error(w, fmt.Sprintf("Template execution error: %v", err), http.StatusInternalServerError)
		return
	}
}

func (d *Dashboard) handleGraphAPI(w http.ResponseWriter, r *http.Request) {
	projectID := r.URL.Path[len("/api/graph/"):]
	if projectID == "" {
		http.Error(w, "Project ID required", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	
	// Get working set version
	workingSet, err := d.queries.GetWorkingSetVersion(ctx, projectID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get working set: %v", err), http.StatusInternalServerError)
		return
	}

	// Get entities and relationships
	entities, err := d.queries.ListEntitiesByVersion(ctx, workingSet.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get entities: %v", err), http.StatusInternalServerError)
		return
	}

	relationships, err := d.queries.ListRelationshipsByVersion(ctx, workingSet.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get relationships: %v", err), http.StatusInternalServerError)
		return
	}

	// Convert to graph visualization format
	graph := GraphVisualization{
		Nodes: make([]Node, len(entities)),
		Links: make([]Link, len(relationships)),
	}

	// Count connections for each entity
	connectionCounts := make(map[string]int)
	for _, rel := range relationships {
		connectionCounts[rel.FromEntityID]++
		connectionCounts[rel.ToEntityID]++
	}

	// Create nodes
	typeGroups := map[string]int{
		"Scene":     1,
		"Character": 2,
		"Location":  3,
		"Theme":     4,
		"PlotPoint": 5,
		"Arc":       6,
	}

	for i, entity := range entities {
		graph.Nodes[i] = Node{
			ID:    entity.ID,
			Name:  entity.Name,
			Type:  entity.EntityType,
			Group: typeGroups[entity.EntityType],
			Size:  connectionCounts[entity.ID],
		}
	}

	// Create links
	for i, rel := range relationships {
		graph.Links[i] = Link{
			Source: rel.FromEntityID,
			Target: rel.ToEntityID,
			Type:   rel.RelationshipType,
			Value:  1,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(graph)
}

func (d *Dashboard) handleStatic(w http.ResponseWriter, r *http.Request) {
	// Serve static files if needed
	http.NotFound(w, r)
}

func (d *Dashboard) getProjectStats(ctx context.Context, versionID string) ProjectStats {
	stats := ProjectStats{
		EntityCounts:       make(map[string]int64),
		RelationshipCounts: make(map[string]int),
	}

	// Get entity counts
	entityTypes := []string{"Scene", "Character", "Location", "Theme", "PlotPoint", "Arc"}
	for _, entityType := range entityTypes {
		params := db.CountEntitiesByTypeParams{
			VersionID:  versionID,
			EntityType: entityType,
		}
		count, err := d.queries.CountEntitiesByType(ctx, params)
		if err == nil {
			stats.EntityCounts[entityType] = count
			stats.TotalEntities += count
		}
	}

	// Get relationship counts
	relationships, err := d.queries.ListRelationshipsByVersion(ctx, versionID)
	if err == nil {
		stats.TotalRelationships = len(relationships)
		for _, rel := range relationships {
			stats.RelationshipCounts[rel.RelationshipType]++
		}
	}

	// Get annotation count (approximate - would need to query all entities)
	entities, err := d.queries.ListEntitiesByVersion(ctx, versionID)
	if err == nil {
		for _, entity := range entities {
			annotations, err := d.queries.ListAnnotationsByEntity(ctx, entity.ID)
			if err == nil {
				stats.TotalAnnotations += len(annotations)
			}
		}
	}

	return stats
}