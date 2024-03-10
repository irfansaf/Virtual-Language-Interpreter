import numpy as np
import torch
import whisper
import pyaudio
import wave

# Load Whisper model
model = whisper.load_model('medium')


# Transcribe audio from microphone
def transcribe_audio():
    CHUNK = 1024
    FORMAT = pyaudio.paInt16
    CHANNELS = 1
    RATE = 16000
    RECORD_SECONDS = 5

    p = pyaudio.PyAudio()

    stream = p.open(format=FORMAT,
                    channels=CHANNELS,
                    rate=RATE,
                    input=True,
                    frames_per_buffer=CHUNK)

    audio_frames = []
    for i in range(0, int(RATE / CHUNK * RECORD_SECONDS)):
        data = stream.read(CHUNK)
        audio_frames.append(data)

    stream.stop_stream()
    stream.close()
    p.terminate()

    audio_data = b''.join(audio_frames)

    # Save audio data to file
    with wave.open('audio.wav', 'wb') as wf:
        wf.setnchannels(CHANNELS)
        wf.setsampwidth(p.get_sample_size(FORMAT))
        wf.setframerate(RATE)
        wf.writeframes(audio_data)

    # Convert audio data to numpy array
    audio = np.frombuffer(audio_data, dtype=np.int16).astype(np.float32) / 32768.0

    # Transcribe audio
    result = model.transcribe(audio, fp16=torch.cuda.is_available(), language='Indonesian')

    return result
