
import './style.css';
import './app.css';

import logo from './assets/images/test.png';
import { StartService, StopService, FetchLogs } from '../wailsjs/go/main/App'; // Assume FetchLogs is a Go method that fetches logs

// Add SVG icon for the button
const playIcon = `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" class="feather feather-play" viewBox="0 0 24 24"><path d="M5 3L19 12 5 21 5 3z"/></svg>`;
const stopIcon = `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" class="feather feather-stop" viewBox="0 0 24 24"><rect width="18" height="18" x="3" y="3" rx="2" ry="2"/></svg>`;

document.querySelector('#app').innerHTML = `
    <img id="logo" class="logo">
    <div class="result" id="result">Service is currently stopped</div>
    <div class="input-box">
        <button id="toggleButton" class="btn">
            ${playIcon} Start
        </button>
    </div>
    <div class="log-container">
        <h3>Service Logs</h3>
        <pre id="logs">No logs available</pre>
    </div>
`;

// Set logo image
document.getElementById('logo').src = logo;

let resultElement = document.getElementById("result");
let buttonElement = document.getElementById("toggleButton");
let logElement = document.getElementById("logs"); // Element to display logs
let isServiceRunning = false; // Tracks if the service is running

const toggleService = () => {
    if (isServiceRunning) {
        StopService()
            .then(() => {
                isServiceRunning = false;
                resultElement.innerText = "Service is currently stopped";
                buttonElement.innerHTML = `${playIcon} Start Service`;
            })
            .catch((err) => {
                console.error("Error stopping service:", err);
            });
    } else {
        StartService()
            .then(() => {
                isServiceRunning = true;
                resultElement.innerText = "Service is currently running";
                buttonElement.innerHTML = `${stopIcon} Stop Service`;
                fetchLogs(); // Fetch logs after starting the service
            })
            .catch((err) => {
                console.error("Error starting service:", err);
            });
    }
};

const fetchLogs = () => {
    FetchLogs()
        .then((logs) => {
            logElement.innerText = JSON.stringify(logs, null, 2); // Pretty print JSON logs
        })
        .catch((err) => {
            logElement.innerText = "Error fetching logs: " + err;
        });
};

// Event listener for service toggle button
buttonElement.addEventListener("click", toggleService);
