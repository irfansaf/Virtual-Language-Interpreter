import numpy as np
import torch
import whisper
import pyaudio
import wave
import io

# Load Whisper model
model = whisper.load_model('medium')

# Function to transcribe audio data
def transcribe_audio(audio_data):
    CHUNK = 1024
    FORMAT = pyaudio.paInt16
    CHANNELS = 1
    RATE = 16000

    p = pyaudio.PyAudio()

    stream = p.open(format=FORMAT,
                    channels=CHANNELS,
                    rate=RATE,
                    input=True,
                    frames_per_buffer=CHUNK)

    audio_frames = io.BytesIO()

    try:
        while True:
            data = stream.read(CHUNK)
            audio_frames.write(data)
    except KeyboardInterrupt:
        pass

    stream.stop_stream()
    stream.close()
    p.terminate()

    audio_frames.seek(0)
    audio_data = audio_frames.read()

    # Convert audio data to numpy array
    audio = np.frombuffer(audio_data, dtype=np.int16).astype(np.float32) / 32768.0

    # Transcribe audio
    result = model.transcribe(audio, fp16=torch.cuda.is_available(), language='Indonesian')

    return result
