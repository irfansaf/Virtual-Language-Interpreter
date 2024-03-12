from flask import Flask, render_template
from flask_cors import CORS
from flask_socketio import SocketIO, emit

from vali.whisper.audio import transcribe_audio

app = Flask(__name__, template_folder='vali/templates')
socketio = SocketIO(app)
CORS(app)

@app.route('/')
def index():
    return render_template('index.html')

@socketio.on('stream')
def handle_stream(data):
    # Pass the audio data to the audio processing function
    transcribed_text = transcribe_audio(data)
    emit('transcription', {'transcribed_text': transcribed_text})

if __name__ == '__main__':
    socketio.run(app, debug=True)
