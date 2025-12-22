import { serve } from "bun";

interface ScriptRequest {
    script: string;
    type?: "typescript" | "javascript";
}

const server = serve({
    port: 3000,
    async fetch(req) {
        const url = new URL(req.url);

        // CORS headers for browser access
        const corsHeaders = {
            "Access-Control-Allow-Origin": "*",
            "Access-Control-Allow-Methods": "POST, GET, OPTIONS",
            "Access-Control-Allow-Headers": "Content-Type",
        };

        // Handle OPTIONS for CORS preflight
        if (req.method === "OPTIONS") {
            return new Response(null, { headers: corsHeaders });
        }

        // Root endpoint - info page
        if (url.pathname === "/" && req.method === "GET") {
            return new Response(
                `
        <html>
          <head>
            <title>Script Runner</title>
            <style>
              body { font-family: system-ui; max-width: 800px; margin: 50px auto; padding: 20px; }
              pre { background: #f5f5f5; padding: 15px; border-radius: 5px; overflow-x: auto; }
              code { background: #e0e0e0; padding: 2px 6px; border-radius: 3px; }
              h1 { color: #333; }
              .example { margin: 20px 0; }
            </style>
          </head>
          <body>
            <h1>🚀 Script Runner Server</h1>
            <p>Server is running on <strong>http://localhost:3000</strong></p>
            
            <h2>Usage</h2>
            <div class="example">
              <h3>Execute a script:</h3>
              <pre><code>curl -X POST http://localhost:3000/run \\
  -H "Content-Type: application/json" \\
  -d '{"script": "console.log(\\"Hello World!\\")", "type": "javascript"}'</code></pre>
            </div>

            <div class="example">
              <h3>Or with TypeScript:</h3>
              <pre><code>curl -X POST http://localhost:3000/run \\
  -H "Content-Type: application/json" \\
  -d '{"script": "const x: number = 42; console.log(x);", "type": "typescript"}'</code></pre>
            </div>

            <div class="example">
              <h3>Using fetch from JavaScript:</h3>
              <pre><code>fetch('http://localhost:3000/run', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    script: 'console.log("Hello from fetch!")',
    type: 'javascript'
  })
})
.then(r => r.json())
.then(console.log);</code></pre>
            </div>

            <h2>Endpoints</h2>
            <ul>
              <li><code>GET /</code> - This info page</li>
              <li><code>POST /run</code> - Execute a script</li>
            </ul>
          </body>
        </html>
        `,
                {
                    headers: {
                        "Content-Type": "text/html",
                        ...corsHeaders,
                    },
                }
            );
        }

        // Execute script endpoint
        if (url.pathname === "/run" && req.method === "POST") {
            try {
                const body = (await req.json()) as ScriptRequest;
                const { script, type = "javascript" } = body;

                if (!script) {
                    return Response.json(
                        { error: "No script provided" },
                        { status: 400, headers: corsHeaders }
                    );
                }

                console.log(`\n📝 Executing ${type} script:`);
                console.log("─".repeat(50));
                console.log(script);
                console.log("─".repeat(50));

                // Capture console output
                const logs: string[] = [];
                const originalLog = console.log;
                const originalError = console.error;
                const originalWarn = console.warn;

                console.log = (...args: any[]) => {
                    logs.push(args.map(a => String(a)).join(" "));
                    originalLog(...args);
                };
                console.error = (...args: any[]) => {
                    logs.push("[ERROR] " + args.map(a => String(a)).join(" "));
                    originalError(...args);
                };
                console.warn = (...args: any[]) => {
                    logs.push("[WARN] " + args.map(a => String(a)).join(" "));
                    originalWarn(...args);
                };

                let result;
                let error = null;

                try {
                    // Execute the script
                    result = eval(script);

                    // If result is a promise, wait for it
                    if (result instanceof Promise) {
                        result = await result;
                    }
                } catch (e: any) {
                    error = {
                        message: e.message,
                        stack: e.stack,
                        name: e.name,
                    };
                    console.error("Script execution error:", e);
                } finally {
                    // Restore console
                    console.log = originalLog;
                    console.error = originalError;
                    console.warn = originalWarn;
                }

                console.log("─".repeat(50));
                console.log(`✅ Script execution completed\n`);

                return Response.json(
                    {
                        success: !error,
                        result: result,
                        logs: logs,
                        error: error,
                        timestamp: new Date().toISOString(),
                    },
                    { headers: corsHeaders }
                );
            } catch (e: any) {
                console.error("Request handling error:", e);
                return Response.json(
                    {
                        success: false,
                        error: {
                            message: e.message,
                            stack: e.stack,
                        },
                    },
                    { status: 500, headers: corsHeaders }
                );
            }
        }

        // 404 for unknown routes
        return new Response("Not Found", { status: 404, headers: corsHeaders });
    },
});

console.log(`
╔════════════════════════════════════════════════════════════╗
║                                                            ║
║  🚀 Script Runner Server                                   ║
║                                                            ║
║  Server running at: http://localhost:${server.port}              ║
║                                                            ║
║  Send POST requests to /run with JSON body:                ║
║  { "script": "your code here", "type": "javascript" }      ║
║                                                            ║
║  Or visit http://localhost:${server.port} for examples           ║
║                                                            ║
╚════════════════════════════════════════════════════════════╝
`);