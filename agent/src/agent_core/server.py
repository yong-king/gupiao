import json
from http.server import BaseHTTPRequestHandler, HTTPServer

from .analyzer import analyze
from .chat_assistant import chat
from .config import load_config
from .health import health_payload
from .langgraph_workflow import run_research_workflow


class Handler(BaseHTTPRequestHandler):
    def do_GET(self):
        if self.path != "/healthz":
            self.send_response(404)
            self.end_headers()
            return

        payload = json.dumps(health_payload()).encode("utf-8")
        self.send_response(200)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(payload)))
        self.end_headers()
        self.wfile.write(payload)

    def do_POST(self):
        if self.path not in {"/analyze", "/workflow/research", "/assistant/chat"}:
            self.send_response(404)
            self.end_headers()
            return
        length = int(self.headers.get("Content-Length", "0"))
        raw = self.rfile.read(length)
        try:
            payload = json.loads(raw.decode("utf-8"))
            cfg = load_config()
            if self.path == "/analyze":
                result = analyze(payload)
            elif self.path == "/workflow/research":
                result = run_research_workflow(payload, cfg)
            else:
                result = chat(payload, cfg)
            body = json.dumps(result).encode("utf-8")
            self.send_response(200)
            self.send_header("Content-Type", "application/json")
            self.send_header("Content-Length", str(len(body)))
            self.end_headers()
            self.wfile.write(body)
        except Exception as exc:
            body = json.dumps({"error": str(exc)}).encode("utf-8")
            self.send_response(400)
            self.send_header("Content-Type", "application/json")
            self.send_header("Content-Length", str(len(body)))
            self.end_headers()
            self.wfile.write(body)

    def log_message(self, format, *args):
        return


def main():
    cfg = load_config()
    server = HTTPServer((cfg.host, cfg.port), Handler)
    print(f"agent listening on {cfg.host}:{cfg.port}")
    server.serve_forever()


if __name__ == "__main__":
    main()
