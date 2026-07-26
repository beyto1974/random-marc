const statusEl = document.getElementById("status");
const previewEl = document.getElementById("preview");
const downloadEl = document.getElementById("download");
const copyEl = document.getElementById("copy");
const formEl = document.getElementById("form");

function setStatus(msg, isError) {
  statusEl.textContent = msg;
  statusEl.classList.toggle("error", !!isError);
}

const go = new Go();
WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject)
  .then((result) => {
    go.run(result.instance);
    setStatus("Ready.");
  })
  .catch((err) => setStatus("Failed to load WebAssembly: " + err, true));

formEl.addEventListener("submit", (event) => {
  event.preventDefault();

  const count = parseInt(document.getElementById("count").value, 10);
  const format = document.getElementById("format").value;
  const seedInput = document.getElementById("seed").value;
  const seed = seedInput === "" ? 0 : parseInt(seedInput, 10);

  const result = window.generateMARC(count, format, seed);

  if (!result.ok) {
    setStatus("Error: " + result.error, true);
    previewEl.textContent = "";
    downloadEl.style.display = "none";
    copyEl.style.display = "none";
    return;
  }

  const bytes = new Uint8Array(result.data);
  const blob = new Blob([bytes], { type: result.mime });
  downloadEl.href = URL.createObjectURL(blob);
  downloadEl.download = "records" + result.ext;
  downloadEl.style.display = "inline-block";

  if (format === "mrc") {
    previewEl.textContent = "";
    copyEl.style.display = "none";
    setStatus(`Generated ${count} record(s), ${bytes.length} bytes (binary - use Download).`);
  } else {
    previewEl.textContent = new TextDecoder("utf-8").decode(bytes);
    copyEl.style.display = "inline-block";
    setStatus(`Generated ${count} record(s).`);
  }
});

copyEl.addEventListener("click", () => {
  navigator.clipboard.writeText(previewEl.textContent).then(
    () => {
      const original = copyEl.textContent;
      copyEl.textContent = "Copied!";
      setTimeout(() => { copyEl.textContent = original; }, 1500);
    },
    (err) => setStatus("Copy failed: " + err, true)
  );
});
