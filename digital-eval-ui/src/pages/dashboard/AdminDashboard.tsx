import React, { useState, useEffect, useRef } from "react";
import Card from "../../components/Card";
import Button from "../../components/Button";
import { performServiceAction, getServiceLogs } from "../../api/admin";

const AdminDashboard: React.FC = () => {
  const [logs, setLogs] = useState<string[]>([
    "[SYSTEM] Dashboard initialized...",
    "[INFO] Connecting to backend..."
  ]);
  const [selectedService, setSelectedService] = useState<string>("IPFS");
  const logEndRef = useRef<HTMLDivElement>(null);

  const addLog = (msg: string) => {
    const timestamp = new Date().toLocaleTimeString();
    setLogs(prev => [...prev, `[${timestamp}] ${msg}`]);
  };

  const handleServiceAction = async (service: string, action: "start" | "stop" | "restart") => {
    addLog(`Requesting ${action} on ${service}...`);
    try {
      await performServiceAction(service, action);
      addLog(`[SUCCESS] ${service} ${action} command sent.`);
    } catch (error: any) {
      addLog(`[ERROR] Failed to ${action} ${service}: ${error.response?.data || error.message}`);
    }
  };

  // Poll logs for selected service
  useEffect(() => {
    const interval = setInterval(async () => {
      try {
        const fetchedLogs = await getServiceLogs(selectedService);
        if (fetchedLogs && fetchedLogs.length > 0) {
           // Simple diffing or just replace for now (optimization: append only new)
           // For this demo, we'll just show the last 50 lines from backend
           setLogs(fetchedLogs);
        }
      } catch (error) {
        // Silent fail on log fetch to avoid spamming
      }
    }, 2000);

    return () => clearInterval(interval);
  }, [selectedService]);

  // Auto-scroll to bottom
  useEffect(() => {
    logEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [logs]);

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold text-slate-800">System Administration</h1>
        <div className="flex gap-2">
          {/* Global actions could be implemented here */}
        </div>
      </div>

      {/* Service Orchestration */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <Card className={`border-l-4 border-blue-500 ${selectedService === "IPFS" ? "ring-2 ring-blue-300" : ""}`} onClick={() => setSelectedService("IPFS")}>
          <div className="flex justify-between items-start mb-4">
            <div>
              <h3 className="font-bold text-lg">IPFS Daemon</h3>
              <span className="text-xs font-mono bg-slate-100 text-slate-700 px-2 py-1 rounded">SERVICE</span>
            </div>
            <div className="flex gap-1">
              <button onClick={(e) => { e.stopPropagation(); handleServiceAction("IPFS", "restart"); }} className="p-1 hover:bg-slate-100 rounded" title="Restart">🔄</button>
              <button onClick={(e) => { e.stopPropagation(); handleServiceAction("IPFS", "stop"); }} className="p-1 hover:bg-red-100 text-red-600 rounded" title="Stop">⏹</button>
              <button onClick={(e) => { e.stopPropagation(); handleServiceAction("IPFS", "start"); }} className="p-1 hover:bg-green-100 text-green-600 rounded" title="Start">▶️</button>
            </div>
          </div>
          <div className="text-sm text-slate-600 space-y-1">
            <p>Click card to view logs</p>
          </div>
        </Card>

        <Card className={`border-l-4 border-cyan-500 ${selectedService === "Python Extractor" ? "ring-2 ring-cyan-300" : ""}`} onClick={() => setSelectedService("Python Extractor")}>
          <div className="flex justify-between items-start mb-4">
            <div>
              <h3 className="font-bold text-lg">Python Extractor</h3>
              <span className="text-xs font-mono bg-slate-100 text-slate-700 px-2 py-1 rounded">SERVICE</span>
            </div>
            <div className="flex gap-1">
              <button onClick={(e) => { e.stopPropagation(); handleServiceAction("Python Extractor", "restart"); }} className="p-1 hover:bg-slate-100 rounded">🔄</button>
              <button onClick={(e) => { e.stopPropagation(); handleServiceAction("Python Extractor", "stop"); }} className="p-1 hover:bg-red-100 text-red-600 rounded">⏹</button>
              <button onClick={(e) => { e.stopPropagation(); handleServiceAction("Python Extractor", "start"); }} className="p-1 hover:bg-green-100 text-green-600 rounded">▶️</button>
            </div>
          </div>
          <div className="text-sm text-slate-600 space-y-1">
            <p>Click card to view logs</p>
          </div>
        </Card>

        <Card className={`border-l-4 border-yellow-500 ${selectedService === "Python Validator" ? "ring-2 ring-yellow-300" : ""}`} onClick={() => setSelectedService("Python Validator")}>
          <div className="flex justify-between items-start mb-4">
            <div>
              <h3 className="font-bold text-lg">Python Validator</h3>
              <span className="text-xs font-mono bg-slate-100 text-slate-700 px-2 py-1 rounded">SERVICE</span>
            </div>
            <div className="flex gap-1">
              <button onClick={(e) => { e.stopPropagation(); handleServiceAction("Python Validator", "restart"); }} className="p-1 hover:bg-slate-100 rounded">🔄</button>
              <button onClick={(e) => { e.stopPropagation(); handleServiceAction("Python Validator", "stop"); }} className="p-1 hover:bg-red-100 text-red-600 rounded">⏹</button>
              <button onClick={(e) => { e.stopPropagation(); handleServiceAction("Python Validator", "start"); }} className="p-1 hover:bg-green-100 text-green-600 rounded">▶️</button>
            </div>
          </div>
          <div className="text-sm text-slate-600 space-y-1">
            <p>Click card to view logs</p>
          </div>
        </Card>
      </div>

      {/* Terminal / Logs */}
      <div className="bg-slate-900 rounded-xl overflow-hidden shadow-xl border border-slate-700">
        <div className="bg-slate-800 px-4 py-2 flex items-center justify-between border-b border-slate-700">
          <span className="text-xs font-mono text-slate-400">logs: {selectedService}</span>
          <div className="flex gap-2">
            <div className="w-3 h-3 rounded-full bg-red-500"></div>
            <div className="w-3 h-3 rounded-full bg-yellow-500"></div>
            <div className="w-3 h-3 rounded-full bg-green-500"></div>
          </div>
        </div>
        <div className="p-4 h-96 overflow-y-auto font-mono text-sm text-green-400 bg-black/50 backdrop-blur">
          {logs.map((log, i) => (
            <div key={i} className="mb-1 border-b border-slate-800/50 pb-1 last:border-0 break-all whitespace-pre-wrap">{log}</div>
          ))}
          <div ref={logEndRef} />
          <div className="animate-pulse">_</div>
        </div>
      </div>
    </div>
  );
};

export default AdminDashboard;
