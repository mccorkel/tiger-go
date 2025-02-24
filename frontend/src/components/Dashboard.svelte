<script lang="ts">
  import { onMount } from 'svelte';
  import { GetStageInfo, CreateParticipantToken, StartWebcamCapture, StartScreenCapture, StartStreaming, StopStreaming } from '../../wailsjs/go/main/App';
  
  let logs: string[] = [];
  let currentStream: MediaStream | null = null;
  let stageArn: string = '';
  let token: string = '';
  let whipConnection: any = null;
  
  function log(message: string) {
    logs = [...logs, `${new Date().toISOString()}: ${message}`];
    console.log(message);
  }

  async function checkPermissions() {
    try {
      log('Requesting camera and microphone permissions...');
      const result = await StartWebcamCapture();
      log('✅ Camera and microphone permissions granted');
      
      // Display preview
      const videoEl = document.getElementById('preview') as HTMLVideoElement;
      if (videoEl) {
        videoEl.srcObject = currentStream;
      }
    } catch (error) {
      log(`❌ Error accessing media devices: ${error}`);
    }
  }
  
  async function getStageInfo() {
    try {
      log('Fetching stage information...');
      const info = await GetStageInfo();
      stageArn = info.StageArn;
      log(`✅ Retrieved stage: ${stageArn}`);
    } catch (error) {
      log(`❌ Error getting stage: ${error}`);
    }
  }
  
  async function getToken() {
    try {
      log('Requesting participant token...');
      const response = await CreateParticipantToken('user-' + Date.now());
      token = response.Token;
      log('✅ Token received');
    } catch (error) {
      log(`❌ Error getting token: ${error}`);
    }
  }
  
  async function startStreaming() {
    if (!token) {
      log('❌ Stream or token not available');
      return;
    }
    
    try {
      log('Initializing Pion WHIP client...');
      log('Connecting to stage...');
      const result = await StartStreaming(token);
      log('✅ Streaming started');
    } catch (error) {
      log(`❌ Streaming error: ${error}`);
    }
  }
  
  async function switchToScreenShare() {
    try {
      log('Requesting screen sharing permissions...');
      const result = await StartScreenCapture();
      log('✅ Screen sharing started');
      
      // Update preview
      const videoEl = document.getElementById('preview') as HTMLVideoElement;
      if (videoEl) {
        videoEl.srcObject = currentStream;
      }
    } catch (error) {
      log(`❌ Screen sharing error: ${error}`);
    }
  }
  
  async function stopStreaming() {
    try {
      log('Stopping stream...');
      const result = await StopStreaming();
      log('✅ Streaming stopped');
    } catch (error) {
      log(`❌ Error stopping stream: ${error}`);
    }
  }
  
  function clearLogs() {
    logs = [];
  }
</script>

<div class="dashboard">
  <div class="controls">
    <h2>IVS Real-Time Streaming</h2>
    
    <div class="button-group">
      <button on:click={checkPermissions}>1. Check Camera Permissions</button>
      <button on:click={getStageInfo}>2. Get Stage Info</button>
      <button on:click={getToken}>3. Get Participant Token</button>
      <button on:click={startStreaming}>4. Start Streaming</button>
      <button on:click={switchToScreenShare}>Switch to Screen Share</button>
      <button on:click={stopStreaming}>Stop Streaming</button>
    </div>
    
    <div class="preview-container">
      <h3>Preview</h3>
      <video id="preview" autoplay muted playsinline></video>
    </div>
  </div>
  
  <div class="logs">
    <div class="log-header">
      <h3>Logs</h3>
      <button class="clear-button" on:click={clearLogs}>Clear</button>
    </div>
    <div class="log-content">
      {#each logs as log}
        <div class="log-entry">{log}</div>
      {/each}
    </div>
  </div>
</div>

<style>
  .dashboard {
    display: flex;
    height: 100vh;
  }
  
  .controls {
    flex: 1;
    padding: 20px;
    border-right: 1px solid #ccc;
  }
  
  .logs {
    flex: 1;
    display: flex;
    flex-direction: column;
    padding: 20px;
  }
  
  .button-group {
    display: flex;
    flex-direction: column;
    gap: 10px;
    margin-bottom: 20px;
  }
  
  button {
    padding: 10px;
    background: #0066cc;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }
  
  button:hover {
    background: #0052a3;
  }
  
  .preview-container {
    margin-top: 20px;
  }
  
  video {
    width: 100%;
    background: #000;
    border-radius: 4px;
  }
  
  .log-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 10px;
  }
  
  .clear-button {
    background: #666;
    padding: 5px 10px;
    font-size: 12px;
  }
  
  .log-content {
    flex: 1;
    overflow-y: auto;
    background: #f5f5f5;
    padding: 10px;
    border-radius: 4px;
    font-family: monospace;
  }
  
  .log-entry {
    margin-bottom: 5px;
    border-bottom: 1px solid #eee;
    padding-bottom: 5px;
  }
</style> 