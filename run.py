from flask import Flask, request, render_template, jsonify

from vali.whisper.audio import transcribe_audio

app = Flask(__name__, template_folder='vali/templates')


@app.route('/')
def index():
    return render_template('index.html')


@app.route('/speech-to-text', methods=['POST'])
def speech_to_text():
    transcribed_text = transcribe_audio()
    return jsonify({'transcribed_text': transcribed_text})


if __name__ == '__main__':
    app.run(debug=True)
