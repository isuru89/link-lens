import { useState, useCallback } from "react";
import { LinkInfo } from "./LinkInfo";

function App() {
  const [url, setUrl] = useState("");
  const [isLoading, setLoading] = useState(false);
  const [data, setData] = useState(null);
  const [errorObj, setError] = useState(null);

  const handleUrlChange = useCallback((e) => {
    setUrl(e.target.value);
  }, []);

  const handleAnalyze = useCallback(() => {
    if (isLoading) {
      return
    }

    setLoading(true);
    fetch("/api/analyze", {
      method: "POST",
      body: JSON.stringify({
        url
      })
    })
    .then(res => res.json())
    .then(data => {
      setData(data)
    }).catch(err => {
      console.error(err);
      setError(err)
    }).finally(() => {
      setLoading(false)
    })
  }, [url, isLoading]);

  return (
    <div>
      <div style={{ display: "flex", flexDirection: "column", width: "60%" }}>
        <label>Enter URL:</label>
        <input type="url" onChange={handleUrlChange} value={url} disabled={isLoading}></input>
        <button onClick={handleAnalyze} disabled={isLoading}>Analyze</button>
      </div>
      {
        errorObj && <div>Error</div>
      }
      {
        data && <LinkInfo data={data}/>
      }
    </div>
  );
}

export default App;
