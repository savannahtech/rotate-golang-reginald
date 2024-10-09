
import './style.css';
import './app.css';

import logo from './assets/images/test.png';
import { StartService, StopService, FetchLogs } from '../wailsjs/go/app/App';

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

document.getElementById('logo').src = logo;

let resultElement = document.getElementById("result");
let buttonElement = document.getElementById("toggleButton");
let logElement = document.getElementById("logs"); 
let isServiceRunning = false;
let pollingInterval = null;   

const fetchLogs = () => {
    FetchLogs()
        .then((logs) => {
            logElement.innerText = logs
        })
        .catch((err) => {
            logElement.innerText = "Error fetching logs: " + err;
        });
};

const startPollingLogs = () => {
    if (pollingInterval === null) {
        pollingInterval = setInterval(fetchLogs, 3000); 
    }
};

const stopPollingLogs = () => {
    if (pollingInterval !== null) {
        clearInterval(pollingInterval);
        pollingInterval = null;
    }
};

const toggleService = () => {
    if (isServiceRunning) {
        StopService()
            .then(() => {
                isServiceRunning = false;
                resultElement.innerText = "Service is currently stopped";
                buttonElement.innerHTML = `${playIcon} Start Service`;
                stopPollingLogs(); 
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
                startPollingLogs(); 
            })
            .catch((err) => {
                console.error("Error starting service:", err);
            });
    }
};

buttonElement.addEventListener("click", toggleService);
