import { useState, useCallback } from "react";
import { LinkInfo } from "./LinkInfo";
import "./App.css";

function App() {
  const [url, setUrl] = useState("");
  const [isLoading, setLoading] = useState(false);
  const [data, setData] = useState(null);
  const [errorObj, setError] = useState(null);
  const [elapsed, setElapsed] = useState(-1);

  const handleUrlChange = useCallback((e) => {
    setUrl(e.target.value);
  }, []);

  const handleAnalyze = useCallback(async () => {
    if (isLoading) {
      return;
    }

    setError(null);
    setData(null);
    setLoading(true);
    setElapsed(-1);
    const t1 = Date.now();
    try {
      const res = await fetch("/api/analyze", {
        method: "POST",
        body: JSON.stringify({
          url,
        }),
      });

      const json = await res.json();
      if (!res.ok) {
        setError(json);
      } else {
        setData(json);
      }
    } catch (err) {
      setError(err);
      console.log(err);
    } finally {
      setLoading(false);
      setElapsed(Date.now() - t1);
    }
  }, [url, isLoading]);

  return (
    <div
      style={{ display: "flex", justifyContent: "center", paddingTop: "120px" }}
    >
      <div style={{ width: "50%", flexDirection: "column" }}>
        <h1 style={{ textAlign: "center", fontSize: "48px" }}>Link-Lens</h1>
        <div
          style={{ display: "flex", flexDirection: "column", width: "100%" }}
        >
          <label>Enter URL:</label>
          <input
            className="url-input"
            type="url"
            onChange={handleUrlChange}
            value={url}
            disabled={isLoading}
          ></input>
          <button
            className="btn-analyze"
            onClick={handleAnalyze}
            disabled={isLoading}
          >
            {isLoading ? "Analyzing..." : "Analyze"}
          </button>
        </div>
        <div style={{ margin: "12px 0" }}>
          {errorObj && (
            <div class="err-label">Error! {errorObj["message"]}</div>
          )}
          {data && <LinkInfo data={data} elapsed={elapsed} />}
        </div>
      </div>
    </div>
  );
}

export default App;
