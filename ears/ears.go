package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html>
		<html>
		<head>
			<title>Record and Upload Audio</title>
		</head>
		<body>
			<button id="start-recording-btn">Start Recording</button>
			<button id="stop-recording-btn" disabled>Stop Recording</button>
			<button id="upload-recording-btn" disabled>Upload Recording</button>
		
			<script>
				// Global variables
				let audioCtx;
				let recorder;
				let recordedBlob;
		
				// Get the buttons
				let startRecordingBtn = document.getElementById("start-recording-btn");
				let stopRecordingBtn = document.getElementById("stop-recording-btn");
				let uploadRecordingBtn = document.getElementById("upload-recording-btn");
		
				// Add event listeners to the buttons
				startRecordingBtn.addEventListener("click", startRecording);
				stopRecordingBtn.addEventListener("click", stopRecording);
				uploadRecordingBtn.addEventListener("click", uploadRecording);
		
				function startRecording() {
					// Request access to the user's microphone
					navigator.mediaDevices.getUserMedia({ audio: true })
						.then(stream => {
							// Initialize the audio context
							audioCtx = new (window.AudioContext || window.webkitAudioContext)();
				
							// Create a MediaStreamSource node
							let source = audioCtx.createMediaStreamSource(stream);
				
							// Create a MediaRecorder
							recorder = new MediaRecorder(stream);
				
							// Start recording
							recorder.start();
				
							// Disable the start recording button
							startRecordingBtn.disabled = true;
				
							// Enable the stop recording button
							stopRecordingBtn.disabled = false;
						})
						.catch(error => {
							console.log(error);
						});
				}
				
				// Stop recording
				function stopRecording() {
					// Stop the recorder
					recorder.stop();
			
					// Disable the stop recording button
					stopRecordingBtn.disabled = true;
			
					// Enable the upload recording button
					uploadRecordingBtn.disabled = false;
			
					// Get the recorded blob
					recorder.ondataavailable = (e) => {
						recordedBlob = e.data;
					}
				}
			
				// Upload recording
				function uploadRecording() {
					// Create a FormData object
					let formData = new FormData();
			
					// Add the recorded blob to the form data
					formData.append("audio", recordedBlob);
			
					// POST the form data to the server
					fetch("/upload", {
						method: "POST",
						body: formData
					})
					.then(response => response.json())
					.then(data => {
						console.log(data);
					})
					.catch(error => {
						console.log(error);
					});
				}
			</script>
			
			</body>
			</html>
`))
	})
	http.HandleFunc("/upload", uploadHandler)
	http.ListenAndServe(":16000", nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		file, _, err := r.FormFile("audio")
		if err != nil {
			http.Error(w, "Error uploading file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		out, err := os.Create("audio.mp3")
		if err != nil {
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, "File uploaded successfully")
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
