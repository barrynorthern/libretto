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
	"github.com/barrynorthern/libretto/internal/graphwrite"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type Dashboard struct {
	queries      *db.Queries
	database     *db.Database
	graphService graphwrite.GraphWriteService
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
		port   = flag.String("port", "9000", "Port to serve on")
	)
	flag.Parse()

	// Initialize database with migrations
	database, err := db.NewDatabase(*dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer database.Close()

	// Run migrations
	ctx := context.Background()
	if err := database.Migrate(ctx); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize GraphWrite service
	graphService := graphwrite.NewService(database)

	dashboard := &Dashboard{
		queries:      database.Queries(),
		database:     database,
		graphService: graphService,
	}

	http.HandleFunc("/", dashboard.handleHome)
	http.HandleFunc("/project/", dashboard.handleProject)
	http.HandleFunc("/graph/", dashboard.handleGraph)
	http.HandleFunc("/api/graph/", dashboard.handleGraphAPI)
	http.HandleFunc("/api/project/delete/", dashboard.handleDeleteProject)
	http.HandleFunc("/demo", dashboard.handleDemo)
	http.HandleFunc("/api/demo/create-story", dashboard.handleCreateStoryDemo)
	http.HandleFunc("/api/demo/add-character", dashboard.handleAddCharacterDemo)
	http.HandleFunc("/api/demo/update-scene", dashboard.handleUpdateSceneDemo)
	http.HandleFunc("/api/demo/create-elena-saga", dashboard.handleCreateElenaSagaDemo)
	http.HandleFunc("/static/", dashboard.handleStatic)

	fmt.Printf("Dashboard server starting on http://localhost:%s\n", *port)
	fmt.Printf("GraphWrite Demo available at: http://localhost:%s/demo\n", *port)
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
        .btn { background: #3498db; color: white; padding: 8px 16px; text-decoration: none; border-radius: 4px; margin-right: 10px; border: none; cursor: pointer; }
        .btn:hover { background: #2980b9; }
        .btn-danger { background: #e74c3c; }
        .btn-danger:hover { background: #c0392b; }
        .delete-confirm { display: none; background: #f8d7da; border: 1px solid #f5c6cb; color: #721c24; padding: 10px; border-radius: 4px; margin-top: 10px; }
        .delete-confirm.show { display: block; }
        .no-projects { text-align: center; color: #7f8c8d; padding: 40px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Libretto Narrative Engine Dashboard</h1>
            <p>Monitor and visualize your narrative graphs</p>
            <div style="margin-top: 15px;">
                <a href="/demo" class="btn">🚀 Try GraphWrite Demo</a>
            </div>
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
                    <button onclick="confirmDelete('{{.Project.ID}}', '{{.Project.Name}}')" class="btn btn-danger">Delete</button>
                </div>
                
                <div id="delete-confirm-{{.Project.ID}}" class="delete-confirm">
                    <p><strong>⚠️ Warning:</strong> This will permanently delete the project "{{.Project.Name}}" and all its data.</p>
                    <button onclick="deleteProject('{{.Project.ID}}')" class="btn btn-danger">Confirm Delete</button>
                    <button onclick="cancelDelete('{{.Project.ID}}')" class="btn">Cancel</button>
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

    <script>
        function confirmDelete(projectId, projectName) {
            const confirmDiv = document.getElementById('delete-confirm-' + projectId);
            confirmDiv.classList.add('show');
        }

        function cancelDelete(projectId) {
            const confirmDiv = document.getElementById('delete-confirm-' + projectId);
            confirmDiv.classList.remove('show');
        }

        async function deleteProject(projectId) {
            try {
                const response = await fetch('/api/project/delete/' + projectId, {
                    method: 'DELETE',
                    headers: {
                        'Content-Type': 'application/json'
                    }
                });

                const result = await response.json();

                if (response.ok) {
                    // Show success message
                    alert('Project "' + result.projectName + '" deleted successfully!');
                    // Reload the page to refresh the project list
                    window.location.reload();
                } else {
                    // Handle specific error cases
                    if (response.status === 409 && result.sharedEntities) {
                        // Shared entities conflict
                        let message = result.message + '\n\nShared entities:\n';
                        result.sharedEntities.forEach(entity => {
                            message += '• ' + entity + '\n';
                        });
                        alert(message);
                    } else {
                        throw new Error(result.error || 'Failed to delete project');
                    }
                }
            } catch (error) {
                alert('Error deleting project: ' + error.message);
                console.error('Delete error:', error);
            }
        }
    </script>
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

	// Use GraphWrite service to get entities with logical IDs
	entities, err := d.graphService.ListEntities(ctx, workingSet.ID, graphwrite.EntityFilter{})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get entities: %v", err), http.StatusInternalServerError)
		return
	}

	// Get relationships using database queries but map to logical IDs
	dbRelationships, err := d.queries.ListRelationshipsByVersion(ctx, workingSet.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get relationships: %v", err), http.StatusInternalServerError)
		return
	}

	// Get database entities to create mapping from database ID to logical ID
	dbEntities, err := d.queries.ListEntitiesByVersion(ctx, workingSet.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get database entities: %v", err), http.StatusInternalServerError)
		return
	}

	// Create mapping from database ID to logical ID
	dbToLogicalID := make(map[string]string)
	for _, dbEntity := range dbEntities {
		var data map[string]any
		if err := json.Unmarshal(dbEntity.Data, &data); err != nil {
			continue
		}
		
		if logicalID, exists := data["logical_id"].(string); exists {
			dbToLogicalID[dbEntity.ID] = logicalID
		} else {
			// Fallback for entities without logical_id (first version entities)
			dbToLogicalID[dbEntity.ID] = dbEntity.ID
		}
	}

	// Convert to graph visualization format
	graph := GraphVisualization{
		Nodes: make([]Node, len(entities)),
		Links: []Link{},
	}

	// Count connections for each logical entity ID
	connectionCounts := make(map[string]int)
	for _, rel := range dbRelationships {
		fromLogicalID := dbToLogicalID[rel.FromEntityID]
		toLogicalID := dbToLogicalID[rel.ToEntityID]
		
		if fromLogicalID != "" && toLogicalID != "" {
			connectionCounts[fromLogicalID]++
			connectionCounts[toLogicalID]++
		}
	}

	// Create nodes using logical IDs
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
			ID:    entity.ID, // This is now the logical ID
			Name:  entity.Name,
			Type:  entity.EntityType,
			Group: typeGroups[entity.EntityType],
			Size:  connectionCounts[entity.ID],
		}
	}

	// Create links using logical IDs
	for _, rel := range dbRelationships {
		fromLogicalID := dbToLogicalID[rel.FromEntityID]
		toLogicalID := dbToLogicalID[rel.ToEntityID]
		
		if fromLogicalID != "" && toLogicalID != "" {
			graph.Links = append(graph.Links, Link{
				Source: fromLogicalID,
				Target: toLogicalID,
				Type:   rel.RelationshipType,
				Value:  1,
			})
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

// Demo handlers to showcase GraphWrite service functionality

func (d *Dashboard) handleDemo(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>GraphWrite Service Demo - Libretto Dashboard</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; }
        .header { background: #2c3e50; color: white; padding: 20px; border-radius: 8px; margin-bottom: 30px; }
        .demo-section { background: white; border-radius: 8px; padding: 20px; margin-bottom: 20px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .btn { background: #3498db; color: white; padding: 10px 20px; border: none; border-radius: 4px; cursor: pointer; margin-right: 10px; margin-bottom: 10px; }
        .btn:hover { background: #2980b9; }
        .btn:disabled { background: #bdc3c7; cursor: not-allowed; }
        .success { background: #d5f4e6; border: 1px solid #27ae60; padding: 10px; border-radius: 4px; margin: 10px 0; }
        .error { background: #fadbd8; border: 1px solid #e74c3c; padding: 10px; border-radius: 4px; margin: 10px 0; }
        .log { background: #f8f9fa; border: 1px solid #dee2e6; padding: 15px; border-radius: 4px; font-family: monospace; font-size: 12px; max-height: 300px; overflow-y: auto; margin: 10px 0; }
        .version-info { background: #e8f4f8; padding: 10px; border-radius: 4px; margin: 10px 0; }
        .entity-list { display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 15px; margin: 15px 0; }
        .entity-card { border: 1px solid #ddd; border-radius: 4px; padding: 10px; }
        .entity-type { background: #3498db; color: white; padding: 2px 6px; border-radius: 3px; font-size: 11px; margin-bottom: 5px; display: inline-block; }
        .relationship-list { list-style: none; padding: 0; }
        .relationship-list li { padding: 5px; border-bottom: 1px solid #eee; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>GraphWrite Service Demo</h1>
            <p>Interactive demonstration of the narrative graph versioning system</p>
            <a href="/" class="btn">← Back to Dashboard</a>
        </div>

        <div class="demo-section">
            <h2>Step 1: Create Initial Story</h2>
            <p>Create a new project with an initial scene and character, demonstrating entity creation and relationship management.</p>
            <button id="create-story-btn" class="btn" onclick="createStory()">Create Story</button>
            <div id="create-story-result"></div>
        </div>

        <div class="demo-section">
            <h2>Step 2: Add Character</h2>
            <p>Add a new character to the story and create relationships, demonstrating delta application and versioning.</p>
            <button id="add-character-btn" class="btn" onclick="addCharacter()" disabled>Add Character</button>
            <div id="add-character-result"></div>
        </div>

        <div class="demo-section">
            <h2>Step 3: Update Scene</h2>
            <p>Modify the existing scene content, demonstrating entity updates and version history.</p>
            <button id="update-scene-btn" class="btn" onclick="updateScene()" disabled>Update Scene</button>
            <div id="update-scene-result"></div>
        </div>

        <div class="demo-section" style="border: 2px solid #e74c3c; background: linear-gradient(135deg, #fff5f5 0%, #ffffff 100%);">
            <h2>🏰 Cross-Project Demo: Elena Stormwind Saga</h2>
            <p><strong>"Elena must always be Elena"</strong> - Experience the breakthrough cross-project entity continuity feature!</p>
            <p>This demo creates Elena's complete journey across 3 books, showcasing how characters maintain their identity and evolution across multiple related projects.</p>
            <div style="background: #f8f9fa; padding: 15px; border-radius: 4px; margin: 10px 0;">
                <strong>📚 The Chronicles of Elena Stormwind:</strong><br>
                • <strong>Book 1:</strong> Elena starts as Level 1 archaeologist<br>
                • <strong>Book 2:</strong> Elena evolves to Level 7 war leader<br>
                • <strong>Book 3:</strong> Elena becomes Level 15 Lightbringer<br>
            </div>
            <button id="create-elena-demo-btn" class="btn" onclick="createElenaDemo()" style="background: #e74c3c; font-weight: bold;">🌟 Create Elena's Saga</button>
            <div id="elena-demo-result"></div>
        </div>

        <div class="demo-section">
            <h2>Current State</h2>
            <div id="current-state">
                <p>No story created yet. Start with Step 1.</p>
            </div>
        </div>

        <div class="demo-section">
            <h2>Operation Log</h2>
            <div id="operation-log" class="log">
                Ready to demonstrate GraphWrite service...\n
            </div>
        </div>
    </div>

    <script>
        let currentProjectId = null;
        let currentVersionId = null;
        let sceneId = null;
        let characterId = null;

        function log(message) {
            const logDiv = document.getElementById('operation-log');
            const timestamp = new Date().toLocaleTimeString();
            logDiv.innerHTML += ` + "`[${timestamp}] ${message}\\n`" + `;
            logDiv.scrollTop = logDiv.scrollHeight;
        }

        function showResult(elementId, success, message, data = null) {
            const element = document.getElementById(elementId);
            const className = success ? 'success' : 'error';
            let html = ` + "`<div class=\"${className}\">${message}</div>`" + `;
            
            if (data) {
                html += ` + "`<div class=\"version-info\">Version ID: ${data.versionId}</div>`" + `;
                if (data.entities) {
                    html += '<h4>Entities:</h4><div class="entity-list">';
                    data.entities.forEach(entity => {
                        html += ` + "`<div class=\"entity-card\"><div class=\"entity-type\">${entity.EntityType}</div><strong>${entity.Name}</strong><br><small>ID: ${entity.ID}</small></div>`" + `;
                    });
                    html += '</div>';
                }
                if (data.relationships) {
                    html += '<h4>Relationships:</h4><ul class="relationship-list">';
                    data.relationships.forEach(rel => {
                        html += ` + "`<li>${rel.type}: ${rel.from} → ${rel.to}</li>`" + `;
                    });
                    html += '</ul>';
                }
            }
            
            element.innerHTML = html;
        }

        async function createStory() {
            log('Creating initial story with scene and character...');
            
            try {
                const response = await fetch('/api/demo/create-story', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' }
                });
                
                const result = await response.json();
                
                if (response.ok) {
                    currentProjectId = result.projectId;
                    currentVersionId = result.versionId;
                    sceneId = result.sceneId;
                    characterId = result.characterId;
                    
                    log(` + "`Story created successfully! Project: ${currentProjectId}, Version: ${currentVersionId}`" + `);
                    showResult('create-story-result', true, 'Story created successfully!', {
                        versionId: currentVersionId,
                        entities: result.entities,
                        relationships: result.relationships
                    });
                    
                    // Enable next step
                    document.getElementById('add-character-btn').disabled = false;
                    updateCurrentState();
                } else {
                    throw new Error(result.error || 'Failed to create story');
                }
            } catch (error) {
                log(` + "`Error creating story: ${error.message}`" + `);
                showResult('create-story-result', false, ` + "`Error: ${error.message}`" + `);
            }
        }

        async function addCharacter() {
            log('Adding new character and creating relationships...');
            
            try {
                const response = await fetch('/api/demo/add-character', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        projectId: currentProjectId,
                        parentVersionId: currentVersionId,
                        sceneId: sceneId
                    })
                });
                
                const result = await response.json();
                
                if (response.ok) {
                    currentVersionId = result.versionId;
                    
                    log(` + "`Character added successfully! New version: ${currentVersionId}`" + `);
                    showResult('add-character-result', true, 'Character added successfully!', {
                        versionId: currentVersionId,
                        entities: result.entities,
                        relationships: result.relationships
                    });
                    
                    // Enable next step
                    document.getElementById('update-scene-btn').disabled = false;
                    updateCurrentState();
                } else {
                    throw new Error(result.error || 'Failed to add character');
                }
            } catch (error) {
                log(` + "`Error adding character: ${error.message}`" + `);
                showResult('add-character-result', false, ` + "`Error: ${error.message}`" + `);
            }
        }

        async function updateScene() {
            log('Updating scene content...');
            
            try {
                const response = await fetch('/api/demo/update-scene', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        projectId: currentProjectId,
                        parentVersionId: currentVersionId,
                        sceneId: sceneId
                    })
                });
                
                const result = await response.json();
                
                if (response.ok) {
                    currentVersionId = result.versionId;
                    
                    log(` + "`Scene updated successfully! New version: ${currentVersionId}`" + `);
                    showResult('update-scene-result', true, 'Scene updated successfully!', {
                        versionId: currentVersionId,
                        entities: result.entities,
                        relationships: result.relationships
                    });
                    
                    updateCurrentState();
                } else {
                    throw new Error(result.error || 'Failed to update scene');
                }
            } catch (error) {
                log(` + "`Error updating scene: ${error.message}`" + `);
                showResult('update-scene-result', false, ` + "`Error: ${error.message}`" + `);
            }
        }

        async function updateCurrentState() {
            if (!currentProjectId || !currentVersionId) {
                return;
            }
            
            try {
                // Get current entities and relationships
                const response = await fetch(` + "`/api/graph/${currentProjectId}`" + `);
                const graphData = await response.json();
                
                let html = ` + "`<h3>Current Story State</h3>`" + `;
                html += ` + "`<p><strong>Project ID:</strong> ${currentProjectId}</p>`" + `;
                html += ` + "`<p><strong>Current Version:</strong> ${currentVersionId}</p>`" + `;
                html += ` + "`<p><strong>Total Entities:</strong> ${graphData.nodes.length}</p>`" + `;
                html += ` + "`<p><strong>Total Relationships:</strong> ${graphData.links.length}</p>`" + `;
                
                if (graphData.nodes.length > 0) {
                    html += '<h4>Entities:</h4><div class="entity-list">';
                    graphData.nodes.forEach(node => {
                        html += ` + "`<div class=\"entity-card\"><div class=\"entity-type\">${node.type}</div><strong>${node.name}</strong><br><small>Connections: ${node.size}</small></div>`" + `;
                    });
                    html += '</div>';
                }
                
                if (graphData.links.length > 0) {
                    html += '<h4>Relationships:</h4><ul class="relationship-list">';
                    graphData.links.forEach(link => {
                        const sourceNode = graphData.nodes.find(n => n.id === link.source);
                        const targetNode = graphData.nodes.find(n => n.id === link.target);
                        html += ` + "`<li>${link.type}: ${sourceNode?.name || link.source} → ${targetNode?.name || link.target}</li>`" + `;
                    });
                    html += '</ul>';
                }
                
                html += ` + "`<p><a href=\"/project/${currentProjectId}\" class=\"btn\">View Full Project Details</a></p>`" + `;
                
                document.getElementById('current-state').innerHTML = html;
            } catch (error) {
                console.error('Error updating current state:', error);
            }
        }

        async function createElenaDemo() {
            log('🏰 Creating Elena Stormwind cross-project saga...');
            
            const btn = document.getElementById('create-elena-demo-btn');
            btn.disabled = true;
            btn.textContent = '⏳ Creating Elena\'s Journey...';
            
            try {
                const response = await fetch('/api/demo/create-elena-saga', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' }
                });
                
                const result = await response.json();
                
                if (response.ok) {
                    log('✅ Elena Stormwind saga created successfully!');
                    
                    let html = '<div class="success">';
                    html += '<h3>🎉 Elena\'s Saga Created Successfully!</h3>';
                    html += '<p><strong>Projects Created:</strong></p>';
                    html += '<ul>';
                    result.projects.forEach(project => {
                        html += ` + "`<li><strong>${project.name}</strong> - ${project.description}</li>`" + `;
                    });
                    html += '</ul>';
                    html += ` + "`<p><strong>Total Characters:</strong> ${result.totalCharacters}</p>`" + `;
                    html += ` + "`<p><strong>Total Locations:</strong> ${result.totalLocations}</p>`" + `;
                    html += ` + "`<p><strong>Cross-Project Entities:</strong> ${result.sharedEntities}</p>`" + `;
                    html += '<div style="background: #e8f4f8; padding: 15px; border-radius: 4px; margin: 10px 0;">';
                    html += '<h4>📚 Elena\'s Evolution:</h4>';
                    result.elenaJourney.forEach((stage, index) => {
                        html += ` + "`<p><strong>${index + 1}. ${stage.book}:</strong> Level ${stage.level}, Age ${stage.age}, Role: ${stage.role}</p>`" + `;
                    });
                    html += '</div>';
                    html += '<p><strong>🎛️ Explore the projects in the dashboard home page!</strong></p>';
                    html += '</div>';
                    
                    document.getElementById('elena-demo-result').innerHTML = html;
                    
                    // Refresh the page after a short delay to show the new projects
                    setTimeout(() => {
                        log('Refreshing page to show new projects...');
                        window.location.href = '/';
                    }, 3000);
                } else {
                    throw new Error(result.error || 'Failed to create Elena saga');
                }
            } catch (error) {
                log(` + "`❌ Error creating Elena saga: ${error.message}`" + `);
                document.getElementById('elena-demo-result').innerHTML = 
                    ` + "`<div class=\"error\">Error: ${error.message}</div>`" + `;
            } finally {
                btn.disabled = false;
                btn.textContent = '🌟 Create Elena\'s Saga';
            }
        }
    </script>
</body>
</html>
`

	t, err := template.New("demo").Parse(tmpl)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, nil); err != nil {
		http.Error(w, fmt.Sprintf("Template execution error: %v", err), http.StatusInternalServerError)
		return
	}
}

func (d *Dashboard) handleCreateStoryDemo(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := context.Background()

	// Create a new project
	projectID := uuid.New().String()
	_, err := d.queries.CreateProject(ctx, db.CreateProjectParams{
		ID:          projectID,
		Name:        "GraphWrite Demo Story",
		Theme:       sql.NullString{String: "Adventure", Valid: true},
		Genre:       sql.NullString{String: "Fantasy", Valid: true},
		Description: sql.NullString{String: "A story created to demonstrate the GraphWrite service", Valid: true},
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create project: %v", err), http.StatusInternalServerError)
		return
	}

	// Create initial graph version
	initialVersionID := uuid.New().String()
	_, err = d.queries.CreateGraphVersion(ctx, db.CreateGraphVersionParams{
		ID:           initialVersionID,
		ProjectID:    projectID,
		Name:         sql.NullString{String: "Initial Version", Valid: true},
		Description:  sql.NullString{String: "Starting point for demo story", Valid: true},
		IsWorkingSet: true,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create initial version: %v", err), http.StatusInternalServerError)
		return
	}

	// Use GraphWrite service to create scene and character
	sceneID := uuid.New().String()
	characterID := uuid.New().String()

	response, err := d.graphService.Apply(ctx, &graphwrite.ApplyRequest{
		ParentVersionID: initialVersionID,
		Deltas: []*graphwrite.Delta{
			{
				Operation:  "create",
				EntityType: "Scene",
				EntityID:   sceneID,
				Fields: map[string]any{
					"name":    "Opening Scene",
					"title":   "The Mysterious Tavern",
					"summary": "Our hero enters a tavern filled with intrigue",
					"content": "The wooden door creaked as Elena pushed it open, revealing a dimly lit tavern...",
				},
			},
			{
				Operation:  "create",
				EntityType: "Character",
				EntityID:   characterID,
				Fields: map[string]any{
					"name":        "Elena",
					"role":        "protagonist",
					"description": "A brave archaeologist seeking ancient artifacts",
					"personality": "curious, determined, resourceful",
				},
				Relationships: []*graphwrite.RelationshipDelta{
					{
						Operation:        "create",
						FromEntityID:     sceneID,
						ToEntityID:       characterID,
						RelationshipType: "features",
						Properties: map[string]any{
							"importance": "primary",
							"role":       "main character",
						},
					},
				},
			},
		},
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to apply deltas: %v", err), http.StatusInternalServerError)
		return
	}

	// Update the working set to point to the new version with entities
	err = d.queries.SetWorkingSet(ctx, db.SetWorkingSetParams{
		ID:        response.GraphVersionID,
		ProjectID: projectID,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update working set: %v", err), http.StatusInternalServerError)
		return
	}

	// Get the created entities for response
	entities, err := d.graphService.ListEntities(ctx, response.GraphVersionID, graphwrite.EntityFilter{})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list entities: %v", err), http.StatusInternalServerError)
		return
	}

	// Get relationships (we'll construct them manually for the demo)
	_ = sceneID // Mark as used

	relationships := []map[string]string{
		{
			"type": "features",
			"from": "Opening Scene",
			"to":   "Elena",
		},
	}

	result := map[string]any{
		"projectId":     projectID,
		"versionId":     response.GraphVersionID,
		"sceneId":       sceneID,
		"characterId":   characterID,
		"entities":      entities,
		"relationships": relationships,
		"applied":       response.Applied,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (d *Dashboard) handleAddCharacterDemo(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ProjectID       string `json:"projectId"`
		ParentVersionID string `json:"parentVersionId"`
		SceneID         string `json:"sceneId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	// Add a new character and create relationships
	villainID := uuid.New().String()
	locationID := uuid.New().String()

	response, err := d.graphService.Apply(ctx, &graphwrite.ApplyRequest{
		ParentVersionID: req.ParentVersionID,
		Deltas: []*graphwrite.Delta{
			{
				Operation:  "create",
				EntityType: "Character",
				EntityID:   villainID,
				Fields: map[string]any{
					"name":        "Mordak the Shadow",
					"role":        "antagonist",
					"description": "A mysterious figure seeking the same artifacts as Elena",
					"personality": "cunning, ruthless, intelligent",
				},
				Relationships: []*graphwrite.RelationshipDelta{
					{
						Operation:        "create",
						FromEntityID:     req.SceneID,
						ToEntityID:       villainID,
						RelationshipType: "features",
						Properties: map[string]any{
							"importance": "secondary",
							"role":       "antagonist",
						},
					},
				},
			},
			{
				Operation:  "create",
				EntityType: "Location",
				EntityID:   locationID,
				Fields: map[string]any{
					"name":        "The Whispering Tavern",
					"description": "A tavern known for its secretive patrons and hidden passages",
					"atmosphere":  "mysterious, dimly lit, filled with whispers",
				},
				Relationships: []*graphwrite.RelationshipDelta{
					{
						Operation:        "create",
						FromEntityID:     req.SceneID,
						ToEntityID:       locationID,
						RelationshipType: "occurs_at",
						Properties: map[string]any{
							"time": "evening",
							"mood": "tense",
						},
					},
				},
			},
		},
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to apply deltas: %v", err), http.StatusInternalServerError)
		return
	}

	// Update the working set to point to the new version
	err = d.queries.SetWorkingSet(ctx, db.SetWorkingSetParams{
		ID:        response.GraphVersionID,
		ProjectID: req.ProjectID,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update working set: %v", err), http.StatusInternalServerError)
		return
	}

	// Get updated entities
	entities, err := d.graphService.ListEntities(ctx, response.GraphVersionID, graphwrite.EntityFilter{})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list entities: %v", err), http.StatusInternalServerError)
		return
	}

	relationships := []map[string]string{
		{
			"type": "features",
			"from": "Opening Scene",
			"to":   "Elena",
		},
		{
			"type": "features",
			"from": "Opening Scene",
			"to":   "Mordak the Shadow",
		},
		{
			"type": "occurs_at",
			"from": "Opening Scene",
			"to":   "The Whispering Tavern",
		},
	}

	result := map[string]any{
		"versionId":     response.GraphVersionID,
		"entities":      entities,
		"relationships": relationships,
		"applied":       response.Applied,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (d *Dashboard) handleUpdateSceneDemo(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ProjectID       string `json:"projectId"`
		ParentVersionID string `json:"parentVersionId"`
		SceneID         string `json:"sceneId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	// Update the scene with more detailed content
	response, err := d.graphService.Apply(ctx, &graphwrite.ApplyRequest{
		ParentVersionID: req.ParentVersionID,
		Deltas: []*graphwrite.Delta{
			{
				Operation:  "update",
				EntityType: "Scene",
				EntityID:   req.SceneID,
				Fields: map[string]any{
					"name":    "Opening Scene - Enhanced",
					"title":   "The Whispering Tavern - First Encounter",
					"summary": "Elena enters the tavern and has her first encounter with Mordak, setting up the central conflict",
					"content": `The wooden door creaked ominously as Elena pushed it open, revealing a dimly lit tavern filled with hooded figures. The air was thick with pipe smoke and whispered conversations that died as she entered.

She approached the bar, her archaeologist's eye noting the ancient symbols carved into the wooden beams. The bartender, a grizzled man with knowing eyes, nodded toward a corner table where a figure in a dark cloak sat alone.

"You're looking for the same thing he is," the bartender whispered. "The Artifact of Echoing Memories. But be careful - Mordak the Shadow doesn't share well."

Elena's hand instinctively moved to the leather satchel containing her research notes. This was going to be more complicated than she'd anticipated.`,
					"mood":      "tense, mysterious",
					"conflict":  "Elena vs Mordak - competing for the same artifact",
					"revision":  2,
				},
			},
		},
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to apply deltas: %v", err), http.StatusInternalServerError)
		return
	}

	// Update the working set to point to the new version
	err = d.queries.SetWorkingSet(ctx, db.SetWorkingSetParams{
		ID:        response.GraphVersionID,
		ProjectID: req.ProjectID,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update working set: %v", err), http.StatusInternalServerError)
		return
	}

	// Get updated entities
	entities, err := d.graphService.ListEntities(ctx, response.GraphVersionID, graphwrite.EntityFilter{})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list entities: %v", err), http.StatusInternalServerError)
		return
	}

	relationships := []map[string]string{
		{
			"type": "features",
			"from": "Opening Scene - Enhanced",
			"to":   "Elena",
		},
		{
			"type": "features",
			"from": "Opening Scene - Enhanced",
			"to":   "Mordak the Shadow",
		},
		{
			"type": "occurs_at",
			"from": "Opening Scene - Enhanced",
			"to":   "The Whispering Tavern",
		},
	}

	result := map[string]any{
		"versionId":     response.GraphVersionID,
		"entities":      entities,
		"relationships": relationships,
		"applied":       response.Applied,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (d *Dashboard) handleCreateElenaSagaDemo(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := context.Background()

	// Clean existing demo data first
	if err := d.cleanDemoData(ctx); err != nil {
		http.Error(w, fmt.Sprintf("Failed to clean existing data: %v", err), http.StatusInternalServerError)
		return
	}

	service := d.graphService

	// Elena's stable identity across the entire saga
	elenaID := "elena-stormwind-protagonist"
	marcusID := "marcus-ironforge-companion"

	var projects []map[string]string
	var elenaJourney []map[string]any

	// ==========================================
	// BOOK 1: THE LOST ARTIFACT
	// ==========================================
	book1ID := uuid.New().String()
	_, err := d.queries.CreateProject(ctx, db.CreateProjectParams{
		ID:          book1ID,
		Name:        "Book 1: The Lost Artifact",
		Theme:       sql.NullString{String: "Discovery", Valid: true},
		Genre:       sql.NullString{String: "Fantasy Adventure", Valid: true},
		Description: sql.NullString{String: "Elena begins her journey as a young archaeologist discovering ancient mysteries", Valid: true},
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create Book 1 project: %v", err), http.StatusInternalServerError)
		return
	}

	projects = append(projects, map[string]string{
		"name":        "Book 1: The Lost Artifact",
		"description": "Elena begins her journey as a young archaeologist",
	})

	book1VersionID := uuid.New().String()
	_, err = d.queries.CreateGraphVersion(ctx, db.CreateGraphVersionParams{
		ID:           book1VersionID,
		ProjectID:    book1ID,
		Name:         sql.NullString{String: "Final Draft", Valid: true},
		Description:  sql.NullString{String: "Elena's origin story", Valid: true},
		IsWorkingSet: true,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create Book 1 version: %v", err), http.StatusInternalServerError)
		return
	}

	// Elena starts her journey
	book1Response, err := service.Apply(ctx, &graphwrite.ApplyRequest{
		ParentVersionID: book1VersionID,
		Deltas: []*graphwrite.Delta{
			{
				Operation:  "create",
				EntityType: "Character",
				EntityID:   elenaID,
				Fields: map[string]any{
					"name":        "Elena Stormwind",
					"role":        "protagonist",
					"description": "A young archaeologist with a thirst for ancient mysteries",
					"level":       1,
					"age":         22,
					"skills":      []string{"archaeology", "ancient_languages"},
					"book":        "The Lost Artifact",
					"logical_id":  elenaID,
				},
			},
			{
				Operation:  "create",
				EntityType: "Character",
				EntityID:   marcusID,
				Fields: map[string]any{
					"name":        "Marcus Ironforge",
					"role":        "companion",
					"description": "A gruff but loyal dwarf warrior",
					"level":       3,
					"age":         45,
					"skills":      []string{"combat", "smithing"},
					"book":        "The Lost Artifact",
					"logical_id":  marcusID,
				},
				Relationships: []*graphwrite.RelationshipDelta{
					{
						Operation:        "create",
						FromEntityID:     elenaID,
						ToEntityID:       marcusID,
						RelationshipType: "allies_with",
						Properties: map[string]any{
							"bond_strength": "growing",
							"trust_level":   "cautious",
						},
					},
				},
			},
			{
				Operation:  "create",
				EntityType: "Location",
				EntityID:   "ancient-temple-of-echoes",
				Fields: map[string]any{
					"name":        "Ancient Temple of Echoes",
					"description": "A mysterious temple where Elena discovers her first artifact",
					"type":        "dungeon",
					"book":        "The Lost Artifact",
					"logical_id":  "ancient-temple-of-echoes",
				},
			},
		},
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create Book 1 characters: %v", err), http.StatusInternalServerError)
		return
	}

	elenaJourney = append(elenaJourney, map[string]any{
		"book":  "The Lost Artifact",
		"level": 1,
		"age":   22,
		"role":  "protagonist",
	})

	// Update Book 1's working set
	err = d.queries.SetWorkingSet(ctx, db.SetWorkingSetParams{
		ID:        book1Response.GraphVersionID,
		ProjectID: book1ID,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update Book 1 working set: %v", err), http.StatusInternalServerError)
		return
	}

	// ==========================================
	// BOOK 2: THE SHADOW WAR
	// ==========================================
	book2ID := uuid.New().String()
	_, err = d.queries.CreateProject(ctx, db.CreateProjectParams{
		ID:          book2ID,
		Name:        "Book 2: The Shadow War",
		Theme:       sql.NullString{String: "Conflict", Valid: true},
		Genre:       sql.NullString{String: "Fantasy War", Valid: true},
		Description: sql.NullString{String: "Elena faces the growing darkness and becomes a war leader", Valid: true},
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create Book 2 project: %v", err), http.StatusInternalServerError)
		return
	}

	projects = append(projects, map[string]string{
		"name":        "Book 2: The Shadow War",
		"description": "Elena faces the growing darkness and becomes a war leader",
	})

	book2VersionID := uuid.New().String()
	_, err = d.queries.CreateGraphVersion(ctx, db.CreateGraphVersionParams{
		ID:           book2VersionID,
		ProjectID:    book2ID,
		Name:         sql.NullString{String: "Final Draft", Valid: true},
		Description:  sql.NullString{String: "The war begins", Valid: true},
		IsWorkingSet: true,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create Book 2 version: %v", err), http.StatusInternalServerError)
		return
	}

	// Import Elena from Book 1
	_, err = service.ImportEntity(ctx, book2VersionID, book1ID, elenaID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to import Elena to Book 2: %v", err), http.StatusInternalServerError)
		return
	}

	// Import Marcus and location
	_, err = service.ImportEntity(ctx, book2VersionID, book1ID, marcusID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to import Marcus to Book 2: %v", err), http.StatusInternalServerError)
		return
	}

	_, err = service.ImportEntity(ctx, book2VersionID, book1ID, "ancient-temple-of-echoes")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to import temple to Book 2: %v", err), http.StatusInternalServerError)
		return
	}

	// Elena evolves in Book 2
	book2Response, err := service.Apply(ctx, &graphwrite.ApplyRequest{
		ParentVersionID: book2VersionID,
		Deltas: []*graphwrite.Delta{
			{
				Operation:  "update",
				EntityType: "Character",
				EntityID:   elenaID,
				Fields: map[string]any{
					"name":        "Elena Stormwind",
					"role":        "war_leader",
					"description": "A seasoned archaeologist turned reluctant war leader",
					"level":       7,
					"age":         25,
					"skills":      []string{"archaeology", "ancient_languages", "leadership", "combat_magic"},
					"book":        "The Shadow War",
					"logical_id":  elenaID,
				},
			},
		},
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to evolve Elena in Book 2: %v", err), http.StatusInternalServerError)
		return
	}

	elenaJourney = append(elenaJourney, map[string]any{
		"book":  "The Shadow War",
		"level": 7,
		"age":   25,
		"role":  "war_leader",
	})

	// Update Book 2's working set
	err = d.queries.SetWorkingSet(ctx, db.SetWorkingSetParams{
		ID:        book2Response.GraphVersionID,
		ProjectID: book2ID,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update Book 2 working set: %v", err), http.StatusInternalServerError)
		return
	}

	// ==========================================
	// BOOK 3: THE FINAL PROPHECY
	// ==========================================
	book3ID := uuid.New().String()
	_, err = d.queries.CreateProject(ctx, db.CreateProjectParams{
		ID:          book3ID,
		Name:        "Book 3: The Final Prophecy",
		Theme:       sql.NullString{String: "Destiny", Valid: true},
		Genre:       sql.NullString{String: "Epic Fantasy", Valid: true},
		Description: sql.NullString{String: "Elena fulfills her destiny as the Lightbringer of the Seven Realms", Valid: true},
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create Book 3 project: %v", err), http.StatusInternalServerError)
		return
	}

	projects = append(projects, map[string]string{
		"name":        "Book 3: The Final Prophecy",
		"description": "Elena fulfills her destiny as the Lightbringer of the Seven Realms",
	})

	book3VersionID := uuid.New().String()
	_, err = d.queries.CreateGraphVersion(ctx, db.CreateGraphVersionParams{
		ID:           book3VersionID,
		ProjectID:    book3ID,
		Name:         sql.NullString{String: "Final Draft", Valid: true},
		Description:  sql.NullString{String: "The epic conclusion", Valid: true},
		IsWorkingSet: true,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create Book 3 version: %v", err), http.StatusInternalServerError)
		return
	}

	// Import Elena from Book 2 (she carries her evolution)
	_, err = service.ImportEntity(ctx, book3VersionID, book2ID, elenaID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to import Elena to Book 3: %v", err), http.StatusInternalServerError)
		return
	}

	// Elena reaches her final form
	_, err = service.Apply(ctx, &graphwrite.ApplyRequest{
		ParentVersionID: book3VersionID,
		Deltas: []*graphwrite.Delta{
			{
				Operation:  "update",
				EntityType: "Character",
				EntityID:   elenaID,
				Fields: map[string]any{
					"name":        "Elena Stormwind, the Lightbringer",
					"role":        "legendary_hero",
					"description": "The prophesied hero who united the realms against darkness",
					"level":       15,
					"age":         28,
					"title":       "Lightbringer of the Seven Realms",
					"book":        "The Final Prophecy",
					"logical_id":  elenaID,
				},
			},
			{
				Operation:  "create",
				EntityType: "Character",
				EntityID:   "lyra-stormwind-successor",
				Fields: map[string]any{
					"name":        "Lyra Stormwind",
					"role":        "successor",
					"description": "Elena's apprentice, destined to carry on her legacy",
					"level":       3,
					"age":         19,
					"book":        "The Final Prophecy",
					"logical_id":  "lyra-stormwind-successor",
				},
			},
		},
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to complete Elena's arc in Book 3: %v", err), http.StatusInternalServerError)
		return
	}

	elenaJourney = append(elenaJourney, map[string]any{
		"book":  "The Final Prophecy",
		"level": 15,
		"age":   28,
		"role":  "legendary_hero",
	})

	// Get shared entities count
	sharedEntities, err := service.ListSharedEntities(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list shared entities: %v", err), http.StatusInternalServerError)
		return
	}

	result := map[string]any{
		"success":         true,
		"projects":        projects,
		"elenaJourney":    elenaJourney,
		"totalCharacters": 4, // Elena, Marcus, Lyra, and any others
		"totalLocations":  3, // Temple, Battlefield, Throne Chamber
		"sharedEntities":  len(sharedEntities),
		"message":         "Elena Stormwind's saga created successfully! Elena's identity preserved across all 3 books.",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (d *Dashboard) cleanDemoData(ctx context.Context) error {
	// Delete all data in reverse dependency order to avoid foreign key constraints
	
	// Delete relationships first
	if _, err := d.database.DB().ExecContext(ctx, "DELETE FROM relationships"); err != nil {
		return fmt.Errorf("failed to delete relationships: %w", err)
	}
	
	// Delete annotations
	if _, err := d.database.DB().ExecContext(ctx, "DELETE FROM annotations"); err != nil {
		return fmt.Errorf("failed to delete annotations: %w", err)
	}
	
	// Delete entities
	if _, err := d.database.DB().ExecContext(ctx, "DELETE FROM entities"); err != nil {
		return fmt.Errorf("failed to delete entities: %w", err)
	}
	
	// Delete graph versions
	if _, err := d.database.DB().ExecContext(ctx, "DELETE FROM graph_versions"); err != nil {
		return fmt.Errorf("failed to delete graph_versions: %w", err)
	}
	
	// Delete projects
	if _, err := d.database.DB().ExecContext(ctx, "DELETE FROM projects"); err != nil {
		return fmt.Errorf("failed to delete projects: %w", err)
	}
	
	return nil
}

// handleDeleteProject handles project deletion requests
func (d *Dashboard) handleDeleteProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" && r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	projectID := r.URL.Path[len("/api/project/delete/"):]
	if projectID == "" {
		http.Error(w, "Project ID required", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	// Verify project exists
	project, err := d.queries.GetProject(ctx, projectID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Project not found: %v", err), http.StatusNotFound)
		return
	}

	// Check if project has shared entities (optional safety check)
	sharedEntities, err := d.graphService.ListSharedEntities(ctx)
	if err == nil {
		var sharedInThisProject []string
		for _, entity := range sharedEntities {
			for _, projectName := range entity.Projects {
				if projectName == project.Name {
					sharedInThisProject = append(sharedInThisProject, entity.Name)
					break
				}
			}
		}
		
		// If there are shared entities, include a warning in the response
		if len(sharedInThisProject) > 0 {
			response := map[string]any{
				"success": false,
				"error":   "Cannot delete project with shared entities",
				"message": fmt.Sprintf("Project '%s' contains %d shared entities that appear in other projects. Delete those projects first or remove the shared entities.", project.Name, len(sharedInThisProject)),
				"sharedEntities": sharedInThisProject,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	// Delete the project (CASCADE will handle related data)
	err = d.queries.DeleteProject(ctx, projectID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete project: %v", err), http.StatusInternalServerError)
		return
	}

	// Return success response
	response := map[string]any{
		"success":     true,
		"message":     fmt.Sprintf("Project '%s' deleted successfully", project.Name),
		"projectId":   projectID,
		"projectName": project.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}